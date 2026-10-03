import { test, expect, type Page } from '@playwright/test';
import { randomBytes, randomUUID } from 'node:crypto';
async function mockTurnstile(page: Page) {
  await page.route('**/turnstile/v0/api.js**', (route) =>
    route.fulfill({
      contentType: 'application/javascript',
      body: `window.turnstile = {
        render(element, options) {
          element.textContent = 'Turnstile test widget';
          options.callback('test-captcha-token');
          return 'test-widget';
        },
        remove() {},
      };`,
    }),
  );
}
async function mockManualTurnstile(page: Page) {
  await page.route('**/turnstile/v0/api.js**', (route) =>
    route.fulfill({
      contentType: 'application/javascript',
      body: `const callbacks = new Map();
        window.turnstile = {
          render(element, options) {
            const id = String(callbacks.size + 1);
            element.textContent = 'Turnstile test widget';
            const complete = () => options.callback('test-captcha-token');
            const fail = () => options['error-callback']();
            window.addEventListener('test-captcha-complete', complete);
            window.addEventListener('test-captcha-fail', fail);
            callbacks.set(id, { complete, fail });
            return id;
          },
          remove(id) {
            const handlers = callbacks.get(id);
            window.removeEventListener('test-captcha-complete', handlers.complete);
            window.removeEventListener('test-captcha-fail', handlers.fail);
            callbacks.delete(id);
          },
        };`,
    }),
  );
}
test('公开首页、账号导航与页面标题', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveTitle('r2f2');
  await expect(page.getByRole('link', { name: 'r2f2' })).toBeVisible();
  await expect(page.locator('link[rel="icon"]')).toHaveAttribute(
    'href',
    /\/icon\.svg/,
  );
  await expect(
    page.getByRole('heading', { name: '你的 r2f2，从这里开始。' }),
  ).toBeVisible();
  await expect(
    page
      .getByRole('navigation', { name: '账号导航' })
      .getByRole('link', { name: '登录' }),
  ).toBeVisible();
  await expect(
    page
      .getByRole('navigation', { name: '账号导航' })
      .getByRole('link', { name: '注册' }),
  ).toHaveCount(0);
  await expect(page.locator('body > footer')).toHaveCount(0);
  await page.goto('/login');
  await expect(page).toHaveTitle('登录 · r2f2');
  await expect(page.getByText('还没有账号？')).toBeVisible();
  await expect(page.getByRole('link', { name: '注册' })).toHaveAttribute(
    'href',
    '/signup',
  );
  await expect(
    page.getByRole('button', { name: '使用 Passkey 登录' }),
  ).toBeVisible();
  await expect(page.getByLabel('邮箱')).not.toBeVisible();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(
    page.getByRole('button', { name: '返回 Passkey 登录' }),
  ).toHaveCount(0);
  await expect(page.locator('.footer-note')).toContainText('还没有账号？ 注册');
  await page.getByRole('button', { name: '使用邮件验证码' }).click();
  await expect(page.getByLabel('邮箱')).toBeVisible();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(
    page.getByRole('button', { name: '返回 Passkey 登录' }),
  ).toBeVisible();
  await page.getByRole('button', { name: '返回 Passkey 登录' }).click();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await page.goto('/signup');
  await expect(page).toHaveTitle('注册 · r2f2');
  await expect(page.locator('.footer-note')).toContainText('已有账号？ 登录');
  await expect(page.getByRole('link', { name: '登录' }).last()).toHaveAttribute(
    'href',
    '/login',
  );
  await page.goto('/settings/profile');
  await expect(page).toHaveTitle('个人资料 · r2f2');
  await expect(page.getByRole('button', { name: '添加 Passkey' })).toHaveCount(
    0,
  );
  await page.goto('/settings/security');
  await expect(page).toHaveTitle('账号与安全 · r2f2');
  await expect(page.getByRole('button', { name: '保存昵称' })).toHaveCount(0);
});
test('注册页点击发送后才验证并传递 token', async ({ page }) => {
  await mockTurnstile(page);
  let sentToken: string | undefined;
  await page.route('**/api/v1/auth/otp/send', (route) => {
    sentToken = route.request().postDataJSON().captchaToken;
    return route.fulfill({
      json: { data: { sent: true, resendAfterSeconds: 60 } },
    });
  });
  await page.goto('/signup');
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await page.getByLabel('邮箱').fill('test@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByLabel('邮件验证码')).toBeVisible();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  expect(sentToken).toBe('test-captcha-token');
});
test('取消验证不发送邮件，再次点击可验证', async ({ page }) => {
  await mockManualTurnstile(page);
  let sends = 0;
  await page.route('**/api/v1/auth/otp/send', (route) => {
    sends += 1;
    return route.fulfill({
      json: { data: { sent: true, resendAfterSeconds: 60 } },
    });
  });
  await page.goto('/signup');
  await page.getByLabel('邮箱').fill('test@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByRole('dialog')).toBeVisible();
  await expect(page.getByText('Turnstile test widget')).toBeVisible();
  expect(sends).toBe(0);
  await page.getByRole('dialog').getByRole('button', { name: '取消' }).click();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  expect(sends).toBe(0);
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByRole('dialog')).toBeVisible();
  await page.evaluate(() =>
    window.dispatchEvent(new Event('test-captcha-complete')),
  );
  await expect(page.getByLabel('邮件验证码')).toBeVisible();
  expect(sends).toBe(1);
});
test('验证失败不发送邮件', async ({ page }) => {
  await mockManualTurnstile(page);
  let sends = 0;
  await page.route('**/api/v1/auth/otp/send', (route) => {
    sends += 1;
    return route.fulfill({ json: { data: { sent: true } } });
  });
  await page.goto('/signup');
  await page.getByLabel('邮箱').fill('test@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByRole('dialog')).toBeVisible();
  await page.evaluate(() =>
    window.dispatchEvent(new Event('test-captcha-fail')),
  );
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(page.getByRole('alert')).toContainText('人机验证未完成');
  expect(sends).toBe(0);
});
test('邮件登录仅在发送和重发时验证', async ({ page }) => {
  await mockTurnstile(page);
  await page.clock.install();
  const tokens: string[] = [];
  await page.route('**/api/v1/auth/otp/send', (route) => {
    tokens.push(route.request().postDataJSON().captchaToken);
    return route.fulfill({
      json: { data: { sent: true, resendAfterSeconds: 60 } },
    });
  });
  await page.route('**/api/v1/auth/otp/verify', (route) =>
    route.fulfill({ json: { data: {} } }),
  );
  await page.goto('/login');
  await page.getByRole('button', { name: '使用邮件验证码' }).click();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await page.getByLabel('邮箱').fill('test@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByLabel('邮件验证码')).toBeVisible();
  expect(tokens).toEqual(['test-captcha-token']);
  await page.clock.runFor(60_000);
  await page.getByRole('button', { name: '重新发送' }).click();
  await expect.poll(() => tokens.length).toBe(2);
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await page.getByLabel('邮件验证码').fill('123456');
  await page.getByRole('button', { name: '验证并继续' }).click();
  await expect(page).toHaveURL('/');
  expect(tokens).toHaveLength(2);
});
test('Passkey 登录在点击并验证后才请求挑战', async ({ page }) => {
  await mockManualTurnstile(page);
  let challenges = 0;
  await page.route('**/api/v1/auth/passkeys/login/options', (route) => {
    challenges += 1;
    expect(route.request().postDataJSON().captchaToken).toBe(
      'test-captcha-token',
    );
    return route.fulfill({
      status: 500,
      json: { error: { code: 'TEST', message: '测试响应' } },
    });
  });
  await page.goto('/login');
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await page.getByRole('button', { name: '使用 Passkey 登录' }).click();
  await expect(page.getByRole('dialog')).toBeVisible();
  expect(challenges).toBe(0);
  await page.evaluate(() =>
    window.dispatchEvent(new Event('test-captcha-complete')),
  );
  await expect.poll(() => challenges).toBe(1);
  await expect(page.getByRole('dialog')).toHaveCount(0);
});
test('重复注册显示已有账号提示并保留登录入口', async ({ page }) => {
  await mockTurnstile(page);
  await page.route('**/api/v1/auth/otp/send', (route) =>
    route.fulfill({
      status: 409,
      json: {
        error: {
          code: 'EMAIL_ALREADY_REGISTERED',
          message: '该邮箱已注册，请登录。',
        },
      },
    }),
  );
  await page.goto('/signup');
  await page.getByLabel('邮箱').fill('existing@example.com');
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByRole('alert')).toContainText('该邮箱已注册，请登录。');
  await expect(page.getByLabel('邮件验证码')).toHaveCount(0);
  await expect(page.getByRole('link', { name: '登录' }).last()).toHaveAttribute(
    'href',
    '/login',
  );
});
test('注册验证码后可以跳过 passkey', async ({ page }) => {
  await mockTurnstile(page);
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
  await mockTurnstile(page);
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
            rp: { id: 'localhost', name: 'r2f2' },
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
  await mockTurnstile(page);
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
  await expect(
    page
      .getByRole('navigation', { name: '账号导航' })
      .getByRole('link', { name: '登录' }),
  ).toBeVisible();
  await expect(
    page
      .getByRole('navigation', { name: '账号导航' })
      .getByRole('link', { name: '注册' }),
  ).toHaveCount(0);
  await expect(page.getByLabel('邮箱')).toBeVisible();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await page.getByLabel('邮箱').fill('test@example.com');
  await mockManualTurnstile(page);
  await page.getByRole('button', { name: '发送验证码' }).click();
  await expect(page.getByRole('dialog')).toBeVisible();
  const bounds = await page.getByRole('dialog').boundingBox();
  expect(bounds).not.toBeNull();
  expect(bounds!.x).toBeGreaterThanOrEqual(0);
  expect(bounds!.x + bounds!.width).toBeLessThanOrEqual(375);
  expect(
    await page.evaluate(() => document.documentElement.scrollWidth),
  ).toBeLessThanOrEqual(375);
  await page.screenshot({
    path: info.outputPath('signup-mobile.png'),
    fullPage: true,
  });
});
