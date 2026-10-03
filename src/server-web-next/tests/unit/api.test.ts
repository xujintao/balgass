import { beforeEach, describe, expect, it, vi } from 'vitest';
const mock = vi.hoisted(() => ({
  jar: new Map<string, string>(),
  user: {
    id: 'user-a',
    email: 'me@example.com',
    email_confirmed_at: '2026-01-01',
  },
  send: vi.fn(),
  findAuthUser: vi.fn(),
  verify: vi.fn(),
  refresh: vi.fn(),
  getUser: vi.fn(),
  signOut: vi.fn(),
  start: vi.fn(),
  passkeyVerify: vi.fn(),
}));
vi.mock('next/headers', () => ({
  cookies: async () => ({
    get: (key: string) =>
      mock.jar.has(key) ? { value: mock.jar.get(key) } : undefined,
    set: (key: string, value: string) => mock.jar.set(key, value),
    delete: (key: string) => mock.jar.delete(key),
  }),
}));
vi.mock('../../src/lib/supabase', () => ({
  supabase: () => ({
    auth: {
      signInWithOtp: mock.send,
      verifyOtp: mock.verify,
      getUser: mock.getUser,
      refreshSession: mock.refresh,
      admin: { signOut: mock.signOut },
      passkey: {
        startAuthentication: mock.start,
        verifyAuthentication: mock.passkeyVerify,
      },
    },
  }),
}));
vi.mock('../../src/lib/auth-user-lookup', () => ({
  findAuthUserByEmail: mock.findAuthUser,
}));
vi.mock('../../src/lib/profile', () => ({
  ensureProfile: vi.fn().mockResolvedValue(undefined),
  me: vi.fn(),
  changeNickname: vi.fn(),
}));
import { handle } from '../../src/lib/api';
import { context, ACCESS_COOKIE, REFRESH_COOKIE } from '../../src/lib/session';
const session = {
  access_token: 'access',
  refresh_token: 'refresh',
  expires_at: 2000000000,
  user: mock.user,
};
function request(
  path: string,
  body: unknown = {},
  headers: Record<string, string> = {},
) {
  return new Request(`http://localhost:3000/api/v1/${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Origin: 'http://localhost:3000',
      ...headers,
    },
    body: JSON.stringify(body),
  });
}
beforeEach(() => {
  vi.clearAllMocks();
  mock.jar.clear();
  process.env.SUPABASE_URL = 'http://127.0.0.1:54321';
  process.env.SUPABASE_PUBLISHABLE_KEY = 'public';
  process.env.APP_ORIGIN = 'http://localhost:3000';
  mock.send.mockResolvedValue({ error: null });
  mock.findAuthUser.mockResolvedValue(null);
  mock.verify.mockResolvedValue({ data: { session }, error: null });
  mock.refresh.mockResolvedValue({ data: { session }, error: null });
  mock.getUser.mockResolvedValue({ data: { user: mock.user }, error: null });
  mock.signOut.mockResolvedValue({ error: null });
});
describe('认证 API', () => {
  it('使用 r2f2 会话 Cookie 名且不读取旧名称', async () => {
    expect(ACCESS_COOKIE).toBe('r2f2_access');
    expect(REFRESH_COOKIE).toBe('r2f2_refresh');
    mock.jar.set('balgass_access', 'legacy');
    await expect(context()).rejects.toMatchObject({ status: 401 });
    expect(mock.getUser).not.toHaveBeenCalled();
  });
  it.each(['signup', 'login'])(
    '发送用途 %s 决定是否允许建号',
    async (intent) => {
      const res = await handle(
        request('auth/otp/send', {
          email: 'me@example.com',
          intent,
          captchaToken: 'test-captcha-token',
        }),
        'auth/otp/send',
      );
      expect(res.status).toBe(200);
      expect(mock.send.mock.calls[0][0].options.shouldCreateUser).toBe(
        intent === 'signup',
      );
      expect(mock.send.mock.calls[0][0].options.captchaToken).toBe(
        'test-captcha-token',
      );
    },
  );
  it('已验证邮箱在注册页返回 409 且不发送邮件', async () => {
    mock.findAuthUser.mockResolvedValue({
      email: 'me@example.com',
      email_confirmed_at: '2026-01-01T00:00:00Z',
    });
    const res = await handle(
      request('auth/otp/send', {
        email: 'ME@example.com',
        intent: 'signup',
        captchaToken: 'test-captcha-token',
      }),
      'auth/otp/send',
    );
    expect(res.status).toBe(409);
    expect((await res.json()).error.code).toBe('EMAIL_ALREADY_REGISTERED');
    expect(mock.findAuthUser).toHaveBeenCalledWith('me@example.com');
    expect(mock.send).not.toHaveBeenCalled();
  });
  it('未验证邮箱继续注册验证码流程', async () => {
    mock.findAuthUser.mockResolvedValue({
      email: 'me@example.com',
      email_confirmed_at: null,
    });
    const res = await handle(
      request('auth/otp/send', {
        email: 'me@example.com',
        intent: 'signup',
        captchaToken: 'test-captcha-token',
      }),
      'auth/otp/send',
    );
    expect(res.status).toBe(200);
    expect(mock.send).toHaveBeenCalledOnce();
  });
  it('查重失败时不发送邮件，登录请求不查重', async () => {
    mock.findAuthUser.mockRejectedValue(new Error('lookup failed'));
    const signup = await handle(
      request('auth/otp/send', {
        email: 'me@example.com',
        intent: 'signup',
        captchaToken: 'test-captcha-token',
      }),
      'auth/otp/send',
    );
    expect(signup.status).toBe(503);
    expect(mock.send).not.toHaveBeenCalled();
    const login = await handle(
      request('auth/otp/send', {
        email: 'me@example.com',
        intent: 'login',
        captchaToken: 'test-captcha-token',
      }),
      'auth/otp/send',
    );
    expect(login.status).toBe(200);
    expect(mock.findAuthUser).toHaveBeenCalledOnce();
  });
  it('不泄露不存在的登录邮箱', async () => {
    mock.send.mockResolvedValue({
      error: { code: 'otp_disabled', status: 400 },
    });
    expect(
      (
        await handle(
          request('auth/otp/send', {
            email: 'none@example.com',
            intent: 'login',
            captchaToken: 'test-captcha-token',
          }),
          'auth/otp/send',
        )
      ).status,
    ).toBe(200);
  });
  it('Cookie 登录不返回令牌', async () => {
    const res = await handle(
      request('auth/otp/verify', { email: 'me@example.com', code: '123456' }),
      'auth/otp/verify',
    );
    const json = await res.json();
    expect(res.status).toBe(200);
    expect(json.data.accessToken).toBeUndefined();
    expect(json.data.refreshToken).toBeUndefined();
    expect(mock.jar.get(ACCESS_COOKIE)).toBe('access');
    expect(mock.jar.get(REFRESH_COOKIE)).toBe('refresh');
  });
  it('App 登录返回令牌且不设置 Cookie', async () => {
    const req = request(
      'auth/otp/verify',
      { email: 'me@example.com', code: '123456' },
      { 'X-Client-Type': 'app' },
    );
    req.headers.delete('Origin');
    const res = await handle(req, 'auth/otp/verify');
    expect((await res.json()).data.refreshToken).toBe('refresh');
    expect(mock.jar.size).toBe(0);
  });
  it.each(['过期', '重用'])('拒绝%s验证码，不产生会话', async () => {
    mock.verify.mockResolvedValue({
      data: { session: null },
      error: { code: 'otp_expired', status: 403 },
    });
    const res = await handle(
      request('auth/otp/verify', { email: 'me@example.com', code: '123456' }),
      'auth/otp/verify',
    );
    expect(res.status).toBe(400);
    expect((await res.json()).error.code).toBe('OTP_INVALID');
    expect(mock.jar.size).toBe(0);
  });
  it('跨站 Cookie 请求被拒绝', async () => {
    expect(
      (
        await handle(
          request('auth/otp/verify', {}, { Origin: 'https://evil.example' }),
          'auth/otp/verify',
        )
      ).status,
    ).toBe(403);
    expect(mock.verify).not.toHaveBeenCalled();
  });
  it('App 标记不能绕过 Cookie 用户身份验证', async () => {
    mock.jar.set(ACCESS_COOKIE, 'cookie-user');
    const req = new Request('http://localhost:3000', {
      headers: { 'X-Client-Type': 'app' },
    });
    await expect(context(req)).rejects.toMatchObject({ status: 401 });
  });
  it('Bearer 优先且由认证服务验证', async () => {
    mock.jar.set(ACCESS_COOKIE, 'other-user');
    const req = new Request('http://localhost:3000', {
      headers: { Authorization: 'Bearer app-token', 'X-Client-Type': 'app' },
    });
    expect((await context(req)).user.id).toBe('user-a');
    expect(mock.getUser).toHaveBeenCalledWith('app-token');
  });
  it('无身份和无效身份返回 401 JSON', async () => {
    const req = new Request('http://localhost:3000/api/v1/me');
    const res = await handle(req, 'me');
    expect(res.status).toBe(401);
    expect((await res.json()).error.code).toBe('UNAUTHENTICATED');
    expect(res.headers.get('Location')).toBeNull();
    mock.getUser.mockResolvedValue({
      data: { user: null },
      error: { status: 401 },
    });
    await expect(
      context(
        new Request('http://localhost:3000', {
          headers: { Authorization: 'Bearer bad', 'X-Client-Type': 'app' },
        }),
      ),
    ).rejects.toMatchObject({ status: 401 });
  });
  it('未验证邮箱不得读取资料', async () => {
    mock.getUser.mockResolvedValue({
      data: { user: { ...mock.user, email_confirmed_at: null } },
      error: null,
    });
    mock.jar.set(ACCESS_COOKIE, 'access');
    await expect(context()).rejects.toMatchObject({ status: 403 });
  });
  it('刷新 Cookie 会话不接受正文中的令牌', async () => {
    mock.jar.set(REFRESH_COOKIE, 'cookie-refresh');
    const res = await handle(
      request('auth/session/refresh', { refreshToken: 'attacker-token' }),
      'auth/session/refresh',
    );
    expect(res.status).toBe(200);
    expect(mock.refresh).toHaveBeenCalledWith({
      refresh_token: 'cookie-refresh',
    });
  });
  it('刷新失效清除 Cookie', async () => {
    mock.jar.set(ACCESS_COOKIE, 'old');
    mock.jar.set(REFRESH_COOKIE, 'old');
    mock.refresh.mockResolvedValue({
      data: { session: null },
      error: { code: 'refresh_token_not_found', status: 400 },
    });
    expect(
      (await handle(request('auth/session/refresh'), 'auth/session/refresh'))
        .status,
    ).toBe(401);
    expect(mock.jar.size).toBe(0);
  });
  it('退出撤销会话并清除 Cookie', async () => {
    mock.jar.set(ACCESS_COOKIE, 'old');
    mock.jar.set(REFRESH_COOKIE, 'old');
    expect((await handle(request('auth/logout'), 'auth/logout')).status).toBe(
      200,
    );
    expect(mock.signOut).toHaveBeenCalledWith('access', 'local');
    expect(mock.jar.size).toBe(0);
  });
  it('邮件限流具有稳定错误与重试时间', async () => {
    mock.send.mockResolvedValue({ error: { status: 429 } });
    const res = await handle(
      request('auth/otp/send', {
        email: 'me@example.com',
        intent: 'login',
        captchaToken: 'test-captcha-token',
      }),
      'auth/otp/send',
    );
    expect(res.status).toBe(429);
    expect(res.headers.get('Retry-After')).toBe('60');
  });
  it.each([
    ['http://localhost:3000', 'auth/otp/send'],
    ['http://localhost:3000', 'auth/passkeys/login/options'],
    ['https://r2f2.example', 'auth/otp/send'],
    ['https://r2f2.example', 'auth/passkeys/login/options'],
  ])('%s 的 %s 缺少 CAPTCHA 时拒绝请求', async (origin, path) => {
    process.env.APP_ORIGIN = origin;
    const res = await handle(
      request(
        path,
        path === 'auth/otp/send'
          ? { email: 'me@example.com', intent: 'signup' }
          : {},
        { Origin: origin },
      ),
      path,
    );
    expect(res.status).toBe(400);
    expect((await res.json()).error.code).toBe('CAPTCHA_REQUIRED');
    expect(mock.send).not.toHaveBeenCalled();
    expect(mock.start).not.toHaveBeenCalled();
  });
  it('Passkey 登录向 Supabase 传递 CAPTCHA token', async () => {
    mock.start.mockResolvedValue({ data: { challenge: 'test' }, error: null });
    const res = await handle(
      request('auth/passkeys/login/options', {
        captchaToken: 'test-captcha-token',
      }),
      'auth/passkeys/login/options',
    );
    expect(res.status).toBe(200);
    expect(mock.start).toHaveBeenCalledWith({
      options: { captchaToken: 'test-captcha-token' },
    });
  });
  it('不缓存私有响应', async () => {
    const res = await handle(
      new Request('http://localhost:3000/api/v1/me'),
      'me',
    );
    expect(res.headers.get('Cache-Control')).toBe('private, no-store');
    expect(res.headers.get('Vary')).toContain('Authorization');
  });
  it('passkey 重放由认证服务拒绝', async () => {
    mock.passkeyVerify.mockResolvedValue({
      data: null,
      error: { code: 'webauthn_challenge_not_found', status: 400 },
    });
    const res = await handle(
      request('auth/passkeys/login/verify', {
        challengeId: '00000000-0000-4000-8000-000000000001',
        credential: {
          id: 'cred',
          rawId: 'cred',
          type: 'public-key',
          response: {
            clientDataJSON: 'YQ',
            authenticatorData: 'YQ',
            signature: 'YQ',
          },
        },
      }),
      'auth/passkeys/login/verify',
    );
    expect(res.status).toBe(400);
    expect((await res.json()).error.code).toBe('WEBAUTHN_CHALLENGE_NOT_FOUND');
  });
  it('未知路径和错误方法返回稳定错误', async () => {
    const notFound = await handle(
      new Request('http://localhost:3000/api/v1/unknown'),
      'unknown',
    );
    expect(notFound.status).toBe(404);
    const wrong = await handle(
      new Request('http://localhost:3000/api/v1/auth/otp/send'),
      'auth/otp/send',
    );
    expect(wrong.status).toBe(405);
    expect(wrong.headers.get('Allow')).toBe('POST');
  });
  it('拒绝过大正文和非 JSON 内容', async () => {
    const huge = await handle(
      request('auth/otp/send', { extra: 'a'.repeat(33000) }),
      'auth/otp/send',
    );
    expect(huge.status).toBe(413);
    const notJson = request('auth/otp/send', {});
    notJson.headers.set('Content-Type', 'text/plain');
    expect((await handle(notJson, 'auth/otp/send')).status).toBe(415);
  });
  it('App 刷新必须显式提供令牌，不能借用 Cookie', async () => {
    mock.jar.set(REFRESH_COOKIE, 'browser-refresh');
    const req = request('auth/session/refresh', {}, { 'X-Client-Type': 'app' });
    req.headers.delete('Origin');
    expect((await handle(req, 'auth/session/refresh')).status).toBe(400);
    expect(mock.refresh).not.toHaveBeenCalled();
  });
  it('浏览器无刷新令牌时清除残留 Cookie', async () => {
    mock.jar.set(ACCESS_COOKIE, 'old');
    expect(
      (await handle(request('auth/session/refresh'), 'auth/session/refresh'))
        .status,
    ).toBe(401);
    expect(mock.jar.size).toBe(0);
  });
  it('认证服务短暂故障不清除 Cookie 或轮换会话', async () => {
    mock.jar.set(ACCESS_COOKIE, 'access');
    mock.jar.set(REFRESH_COOKIE, 'refresh');
    mock.getUser.mockResolvedValue({
      data: { user: null },
      error: { status: 503 },
    });
    const res = await handle(
      new Request('http://localhost:3000/api/v1/me'),
      'me',
    );
    expect(res.status).toBe(502);
    expect(mock.jar.get(REFRESH_COOKIE)).toBe('refresh');
    expect(mock.refresh).not.toHaveBeenCalled();
    mock.refresh.mockResolvedValue({
      data: { session: null },
      error: { status: 503 },
    });
    const refresh = await handle(
      request('auth/session/refresh'),
      'auth/session/refresh',
    );
    expect(refresh.status).toBe(502);
    expect(mock.jar.get(REFRESH_COOKIE)).toBe('refresh');
  });
});
