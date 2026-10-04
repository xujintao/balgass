import { afterEach, describe, expect, it, vi } from 'vitest';
import { createGameAccount, listGameAccounts } from '../../src/lib/game-accounts';

const owner = 'player@example.com';
const ctx = { user: { email: owner, email_confirmed_at: '2026-10-01' } } as Parameters<typeof listGameAccounts>[0];
const token = 'a'.repeat(32);

afterEach(() => {
  vi.unstubAllEnvs();
  vi.unstubAllGlobals();
});

function configured() {
  vi.stubEnv('GAME_API_URL', 'https://game.r2f2.com/api/command');
  vi.stubEnv('GAME_API_TOKEN', token);
}

describe('game account command service', () => {
  it('sends an authenticated list command for the verified owner and narrows the result', async () => {
    configured();
    const fetcher = vi.fn().mockResolvedValue(Response.json({
      action: 'GetAccountList',
      out: [{ name: 'hero', user_email: owner, password: 'secret', characters: [{ name: 'Mage', level: 20, master_level: 3 }] }],
    }));
    vi.stubGlobal('fetch', fetcher);

    await expect(listGameAccounts(ctx)).resolves.toEqual([{ name: 'hero', characters: [{ name: 'Mage', level: 23 }] }]);
    const [url, init] = fetcher.mock.calls[0] as [URL, RequestInit];
    expect(url.toString()).toBe('https://game.r2f2.com/api/command');
    expect(init.headers).toMatchObject({ Authorization: `Bearer ${token}` });
    expect(JSON.parse(init.body as string)).toEqual({ action: 'GetAccountList', in: { user_email: owner } });
  });

  it('sends creation with session ownership and without confirmation password', async () => {
    configured();
    const fetcher = vi.fn().mockResolvedValue(Response.json({ action: 'CreateAccount', out: { name: 'hero', user_email: owner } }));
    vi.stubGlobal('fetch', fetcher);

    await expect(createGameAccount(ctx, { name: 'hero', password: 'secret', passwordConfirmation: 'secret' })).resolves.toEqual({ name: 'hero' });
    expect(JSON.parse(fetcher.mock.calls[0][1].body)).toEqual({
      action: 'CreateAccount', in: { name: 'hero', password: 'secret', user_email: owner },
    });
  });

  it('rejects mismatched actions and another owner in the upstream result', async () => {
    configured();
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(Response.json({ action: 'DeleteAccount', out: [] })));
    await expect(listGameAccounts(ctx)).rejects.toMatchObject({ code: 'GAME_UNAVAILABLE' });
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(Response.json({
      action: 'GetAccountList', out: [{ name: 'other', user_email: 'other@example.com', characters: null }],
    })));
    await expect(listGameAccounts(ctx)).rejects.toMatchObject({ code: 'GAME_UNAVAILABLE' });
  });

  it('treats a nil upstream account slice as an empty list', async () => {
    configured();
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(Response.json({ action: 'GetAccountList', out: null })));
    await expect(listGameAccounts(ctx)).resolves.toEqual([]);
  });

  it('fails closed without a service token or with an insecure production URL', async () => {
    configured();
    vi.stubEnv('GAME_API_TOKEN', '');
    await expect(listGameAccounts(ctx)).rejects.toMatchObject({ code: 'GAME_UNAVAILABLE' });
    vi.stubEnv('GAME_API_TOKEN', token);
    vi.stubEnv('GAME_API_URL', 'http://game.r2f2.com/api/command');
    await expect(listGameAccounts(ctx)).rejects.toMatchObject({ code: 'GAME_UNAVAILABLE' });
  });
});
