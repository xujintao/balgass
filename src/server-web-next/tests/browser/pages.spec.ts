import { test, expect } from '@playwright/test';
import { randomBytes, randomUUID } from 'node:crypto';
test('首页与默认登录入口', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('link', { name: '创建账号' })).toBeVisible();
  await page.goto('/login');
  await expect(
    page.getByRole('button', { name: '使用 Passkey 登录' }),
  ).toBeVisible();
  await expect(page.getByLabel('邮箱')).not.toBeVisible();
  await page.getByRole('button', { name: '使用邮件验证码' }).click();
  await expect(page.getByLabel('邮箱')).toBeVisible();
});
test('注册验证码后可以跳过 passkey', async ({ page }) => {
  await page.route('**/api/v1/auth/otp/send', (route) =>
    route.fulfill({ json: { data: { sent: true, resendAfterSeconds: 60 } } }),
  );
  await page.route('**/api/v1/auth/otp/verify', (route) =>
    route.fulfill({
      json: { data: { user: { id: 'test', email: 'test@example.com' } } },
    }),
  );
  await page.goto('/signup');
  await page.getByLabel('邮箱').fill('test@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await page.getByLabel('邮件验证码').fill('123456');
  await page.getByRole('button', { name: '验证并继续' }).click();
  await expect(
    page.getByRole('heading', { name: '设置你的 Passkey' }),
  ).toBeVisible();
  await page.getByRole('link', { name: '稍后设置' }).click();
  await expect(page).toHaveURL('/');
});
test('真实浏览器虚拟认证器完成注册与登录 ceremony（API 模拟）', async ({
  page,
  context,
}) => {
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
  const registrationChallenge = randomBytes(32).toString('base64url'),
    loginChallenge = randomBytes(32).toString('base64url');
  let credentialId = '';
  await page.route('**/api/v1/auth/otp/send', (r) =>
    r.fulfill({ json: { data: { sent: true, resendAfterSeconds: 60 } } }),
  );
  await page.route('**/api/v1/auth/otp/verify', (r) =>
    r.fulfill({ json: { data: {} } }),
  );
  await page.route('**/api/v1/auth/passkeys/register/options', (r) =>
    r.fulfill({
      json: {
        data: {
          challenge_id: randomUUID(),
          options: {
            rp: { id: 'localhost', name: 'Balgass' },
            user: {
              id: randomBytes(16).toString('base64url'),
              name: 'virtual@example.com',
              displayName: '玩家',
            },
            challenge: registrationChallenge,
            pubKeyCredParams: [{ type: 'public-key', alg: -7 }],
            authenticatorSelection: {
              residentKey: 'required',
              userVerification: 'required',
            },
            attestation: 'none',
          },
        },
      },
    }),
  );
  await page.route('**/api/v1/auth/passkeys/register/verify', async (r) => {
    const input = r.request().postDataJSON();
    const clientData = JSON.parse(
      Buffer.from(
        input.credential.response.clientDataJSON,
        'base64url',
      ).toString(),
    );
    expect(clientData).toMatchObject({
      type: 'webauthn.create',
      origin: 'http://localhost:3000',
      challenge: registrationChallenge,
    });
    expect(input.credential.response.attestationObject).toBeTruthy();
    credentialId = input.credential.id;
    await r.fulfill({ json: { data: { id: randomUUID() } } });
  });
  await page.goto('/signup');
  await page.getByLabel('邮箱').fill('virtual@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await page.getByLabel('邮件验证码').fill('123456');
  await page.getByRole('button', { name: '验证并继续' }).click();
  await page.getByRole('button', { name: '创建 Passkey' }).click();
  await expect(page).toHaveURL('/');
  expect(credentialId).toBeTruthy();
  await page.route('**/api/v1/auth/passkeys/login/options', (r) =>
    r.fulfill({
      json: {
        data: {
          challenge_id: randomUUID(),
          options: {
            challenge: loginChallenge,
            rpId: 'localhost',
            userVerification: 'required',
          },
        },
      },
    }),
  );
  await page.route('**/api/v1/auth/passkeys/login/verify', async (r) => {
    const input = r.request().postDataJSON();
    expect(input.credential.id).toBe(credentialId);
    const clientData = JSON.parse(
      Buffer.from(
        input.credential.response.clientDataJSON,
        'base64url',
      ).toString(),
    );
    expect(clientData).toMatchObject({
      type: 'webauthn.get',
      origin: 'http://localhost:3000',
      challenge: loginChallenge,
    });
    expect(input.credential.response.signature).toBeTruthy();
    await r.fulfill({ json: { data: {} } });
  });
  await page.goto('/login');
  await page.getByRole('button', { name: '使用 Passkey 登录' }).click();
  await expect(page).toHaveURL('/');
});
test('取消 passkey 操作后邮件入口仍可用', async ({ page }) => {
  await page.addInitScript(() => {
    Object.defineProperty(navigator.credentials, 'get', {
      value: () =>
        Promise.reject(new DOMException('cancelled', 'NotAllowedError')),
    });
  });
  await page.route('**/api/v1/auth/passkeys/login/options', (r) =>
    r.fulfill({
      json: {
        data: {
          challenge_id: randomUUID(),
          options: {
            challenge: randomBytes(32).toString('base64url'),
            rpId: 'localhost',
          },
        },
      },
    }),
  );
  await page.goto('/login');
  await page.getByRole('button', { name: '使用 Passkey 登录' }).click();
  await expect(
    page.getByRole('alert').filter({ hasText: '取消' }),
  ).toContainText('取消');
  await page.getByRole('button', { name: '使用邮件验证码' }).click();
  await expect(page.getByLabel('邮箱')).toBeVisible();
});

// Check the actual Next.js handler rather than route mocks.
test('私有 API 返回 JSON 401，浏览器跨站写操作被拒绝', async ({ request }) => {
  const me = await request.get('/api/v1/me');
  expect(me.status()).toBe(401);
  expect((await me.json()).error.code).toBe('UNAUTHENTICATED');
  expect(me.headers()['cache-control']).toBe('private, no-store');
  expect(me.headers()['location']).toBeUndefined();
  const csrf = await request.post('/api/v1/auth/logout', {
    headers: { Origin: 'https://evil.example' },
  });
  expect(csrf.status()).toBe(403);
  expect((await csrf.json()).error.code).toBe('INVALID_ORIGIN');
});

test('移动端注册表单没有横向溢出', async ({ page }, info) => {
  await page.setViewportSize({ width: 375, height: 812 });
  await page.goto('/signup');
  await expect(page.getByLabel('邮箱')).toBeVisible();
  expect(
    await page.evaluate(() => document.documentElement.scrollWidth),
  ).toBeLessThanOrEqual(375);
  await page.screenshot({
    path: info.outputPath('signup-mobile.png'),
    fullPage: true,
  });
});
