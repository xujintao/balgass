import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { findAuthUserByEmail } from '../../src/lib/auth-user-lookup';

const email = 'me@example.com';

beforeEach(() => {
  process.env.SUPABASE_URL = 'https://project.supabase.co';
  process.env.SUPABASE_PUBLISHABLE_KEY = 'public';
  process.env.SUPABASE_SECRET_KEY = 'server-secret';
  process.env.APP_ORIGIN = 'https://example.com';
});
afterEach(() => vi.unstubAllGlobals());

describe('Auth 管理用户查找', () => {
  it.each([
    [[], null],
    [
      [{ email: 'ME@EXAMPLE.COM', email_confirmed_at: '2026-01-01' }],
      '2026-01-01',
    ],
    [[{ email, email_confirmed_at: null }], null],
  ])('完整邮箱匹配并保留验证状态', async (users, confirmedAt) => {
    const fetcher = vi.fn().mockResolvedValue(Response.json({ users }));
    vi.stubGlobal('fetch', fetcher);
    const result = await findAuthUserByEmail(email);
    expect(result?.email_confirmed_at ?? null).toBe(confirmedAt);
    const [url, init] = fetcher.mock.calls[0];
    expect(new URL(url).searchParams.get('filter')).toBe(email);
    expect(init.headers.Authorization).toBe('Bearer server-secret');
    expect(init.cache).toBe('no-store');
  });

  it('过滤到相似邮箱时不误判', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        Response.json({
          users: [
            {
              email: 'other-me@example.com',
              email_confirmed_at: '2026-01-01',
            },
          ],
        }),
      ),
    );
    expect(await findAuthUserByEmail(email)).toBeNull();
  });

  it('首批仅有相似邮箱时继续搜索下一页', async () => {
    const fetcher = vi
      .fn()
      .mockResolvedValueOnce(
        Response.json({
          users: Array.from({ length: 100 }, (_, i) => ({
            email: `other${i}@example.com`,
          })),
        }),
      )
      .mockResolvedValueOnce(
        Response.json({
          users: [
            { email: 'ME@example.com', email_confirmed_at: '2026-01-01' },
          ],
        }),
      );
    vi.stubGlobal('fetch', fetcher);
    expect((await findAuthUserByEmail(email))?.email_confirmed_at).toBe(
      '2026-01-01',
    );
    expect(new URL(fetcher.mock.calls[1][0]).searchParams.get('page')).toBe(
      '2',
    );
  });

  it('搜索失败或响应无效时拒绝继续', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(new Response('', { status: 503 })),
    );
    await expect(findAuthUserByEmail(email)).rejects.toMatchObject({
      status: 502,
    });
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(Response.json({ unexpected: true })),
    );
    await expect(findAuthUserByEmail(email)).rejects.toMatchObject({
      status: 502,
    });
  });
});
