import { test, expect } from '@playwright/test';
import { createClient } from '@supabase/supabase-js';
import { randomUUID } from 'node:crypto';
const live = process.env.E2E_LIVE_AUTH === 'true';
test.describe('真实 Supabase 认证（需要专用测试项目）', () => {
  test.skip(!live, '未配置 E2E_LIVE_AUTH；不执行真实认证联调');
  test('注册、登录、删除、重放及跨用户隔离', async ({
    page,
    context,
    request,
  }) => {
    const url = process.env.SUPABASE_URL!,
      key = process.env.SUPABASE_PUBLISHABLE_KEY!,
      secret = process.env.TEST_SUPABASE_SECRET_KEY!;
    expect(url).toBeTruthy();
    expect(key).toBeTruthy();
    expect(secret).toBeTruthy();
    const admin = createClient(url, secret, {
      auth: { persistSession: false, autoRefreshToken: false },
    });
    const users: string[] = [];
    async function account() {
      const email = `e2e-${randomUUID()}@example.com`;
      const created = await admin.auth.admin.createUser({
        email,
        email_confirm: true,
      });
      expect(created.error).toBeNull();
      users.push(created.data.user!.id);
      const link = await admin.auth.admin.generateLink({
        type: 'magiclink',
        email,
      });
      expect(link.error).toBeNull();
      return { email, code: link.data.properties!.email_otp };
    }
    try {
      const a = await account(),
        b = await account();
      const publicClient = createClient(url, key, {
        auth: { persistSession: false, autoRefreshToken: false },
      });
      const bSession = await publicClient.auth.verifyOtp({
        email: b.email,
        token: b.code,
        type: 'email',
      });
      expect(bSession.error).toBeNull();
      const token = bSession.data.session!.access_token;
      const appHeaders = {
        Authorization: `Bearer ${token}`,
        'X-Client-Type': 'app',
      };
      const cdp = await context.newCDPSession(page);
      await cdp.send('WebAuthn.enable');
      await cdp.send('WebAuthn.addVirtualAuthenticator', {
        options: {
          protocol: 'ctap2',
          transport: 'internal',
          hasResidentKey: true,
          hasUserVerification: true,
          isUserVerified: true,
          automaticPresenceSimulation: true,
        },
      });
      await page.goto('/login');
      await page.getByRole('button', { name: '使用邮件验证码' }).click();
      // Use an admin-generated OTP to test the verification flow without requiring SMTP.
      await page.route('**/api/v1/auth/otp/send', (r) =>
        r.fulfill({ json: { data: { sent: true, resendAfterSeconds: 60 } } }),
      );
      await page.getByLabel('邮箱').fill(a.email);
      await page.getByRole('button', { name: '发送验证码' }).click();
      await page.getByLabel('邮件验证码').fill(a.code);
      await page.getByRole('button', { name: '验证并继续' }).click();
      await expect(page).toHaveURL('/');
      await page.goto('/settings');
      await expect(
        page.getByRole('heading', { name: '账号设置' }),
      ).toBeVisible();
      let registration: unknown;
      let crossStatus = 0;
      await page.route(
        '**/api/v1/auth/passkeys/register/verify',
        async (route) => {
          registration = route.request().postDataJSON();
          const cross = await request.post(
            '/api/v1/auth/passkeys/register/verify',
            { headers: appHeaders, data: registration },
          );
          crossStatus = cross.status();
          await route.continue();
        },
      );
      await page.getByRole('button', { name: '添加 Passkey' }).click();
      await expect(page.getByRole('status')).toContainText('已添加');
      expect(crossStatus).toBeGreaterThanOrEqual(400);
      const replay = await page.request.post(
        '/api/v1/auth/passkeys/register/verify',
        { headers: { Origin: 'http://localhost:3000' }, data: registration },
      );
      expect(replay.status()).toBeGreaterThanOrEqual(400);
      const keysResponse = await page.request.get('/api/v1/auth/passkeys');
      const keys = (await keysResponse.json()).data;
      expect(keys).toHaveLength(1);
      const crossDelete = await request.delete(
        `/api/v1/auth/passkeys/${keys[0].id}`,
        { headers: appHeaders },
      );
      expect(crossDelete.status()).toBeGreaterThanOrEqual(400);
      expect(
        (await (await page.request.get('/api/v1/auth/passkeys')).json()).data,
      ).toHaveLength(1);
      await page.getByRole('button', { name: '退出登录' }).click();
      await expect(page).toHaveURL('/login');
      let assertion: unknown;
      await page.route(
        '**/api/v1/auth/passkeys/login/verify',
        async (route) => {
          assertion = route.request().postDataJSON();
          await route.continue();
        },
      );
      await page.getByRole('button', { name: '使用 Passkey 登录' }).click();
      await expect(page).toHaveURL('/');
      const assertionReplay = await page.request.post(
        '/api/v1/auth/passkeys/login/verify',
        { headers: { Origin: 'http://localhost:3000' }, data: assertion },
      );
      expect(assertionReplay.status()).toBeGreaterThanOrEqual(400);
      await page.goto('/settings');
      page.once('dialog', (dialog) => dialog.accept());
      await page.getByRole('button', { name: /删除/ }).click();
      await expect(page.getByRole('status')).toContainText('已删除');
      expect(
        (await (await page.request.get('/api/v1/auth/passkeys')).json()).data,
      ).toHaveLength(0);
    } finally {
      for (const id of users) await admin.auth.admin.deleteUser(id);
    }
  });
});
