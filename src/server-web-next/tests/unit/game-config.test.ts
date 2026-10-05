import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';
import { afterEach, describe, expect, it, vi } from 'vitest';

vi.mock('server-only', () => ({}));
vi.mock('next/cache', () => ({ unstable_cache: (fn: unknown) => fn }));

import { gameWebSocketUrl, getGameMapNames, parseGameMapNames } from '../../src/lib/game-config';

const fixtureUrl = new URL('../fixtures/', import.meta.url);

afterEach(() => {
  vi.unstubAllEnvs();
  vi.unstubAllGlobals();
});

describe('game map configuration', () => {
  it('matches the old map naming rules', async () => {
    const xml = await readFile(fileURLToPath(new URL('IGC_MapList.xml', fixtureUrl)), 'utf8');
    expect(parseGameMapNames(xml)).toEqual(['Lorencia', 'Dungeon', 'BloodCastle', 'BloodCastle2']);
    expect(() => parseGameMapNames('<MapList><DefaultMaps><Map File="x" /></DefaultMaps></MapList>')).toThrow();
    expect(() => parseGameMapNames('<!DOCTYPE MapList><MapList><DefaultMaps /></MapList>')).toThrow();
  });

  it('parses the repository game map configuration', async () => {
    const realUrl = new URL('../../../../config/server-game-common/IGCData/IGC_MapList.xml', import.meta.url);
    const names = parseGameMapNames(await readFile(fileURLToPath(realUrl), 'utf8'));
    expect(names.length).toBeGreaterThan(0);
    expect(names.every((name) => typeof name === 'string' && name.length > 0)).toBe(true);
    expect(names).not.toContain('Null');
  });

  it('reads the configured local file without a network request', async () => {
    vi.stubEnv('GAME_CONFIG_URL', fixtureUrl.href);
    const fetcher = vi.fn();
    vi.stubGlobal('fetch', fetcher);
    await expect(getGameMapNames()).resolves.toContain('Lorencia');
    expect(fetcher).not.toHaveBeenCalled();
  });

  it('reads local changes on the next request without retaining parsed results', async () => {
    const directory = await mkdtemp(path.join(os.tmpdir(), 'game-config-'));
    try {
      vi.stubEnv('GAME_CONFIG_URL', pathToFileURL(directory + path.sep).href);
      const filename = path.join(directory, 'IGC_MapList.xml');
      await writeFile(filename, '<MapList><DefaultMaps><Map File="00_Lorencia.att" /></DefaultMaps></MapList>');
      await expect(getGameMapNames()).resolves.toEqual(['Lorencia']);
      await writeFile(filename, '<MapList><DefaultMaps><Map File="01_Dungeon.att" /></DefaultMaps></MapList>');
      await expect(getGameMapNames()).resolves.toEqual(['Dungeon']);
    } finally {
      await rm(directory, { recursive: true, force: true });
    }
  });

  it('reads the exact HTTPS file with separate Basic Auth credentials', async () => {
    vi.stubEnv('GAME_CONFIG_URL', 'https://game.r2f2.com/config/');
    vi.stubEnv('GAME_CONFIG_BASIC_AUTH_USER', 'nextjs');
    vi.stubEnv('GAME_CONFIG_BASIC_AUTH_PASSWORD', 'read-only-password');
    const xml = await readFile(fileURLToPath(new URL('IGC_MapList.xml', fixtureUrl)), 'utf8');
    const fetcher = vi.fn().mockResolvedValue(new Response(xml));
    vi.stubGlobal('fetch', fetcher);
    await expect(getGameMapNames()).resolves.toHaveLength(4);
    const [url, init] = fetcher.mock.calls[0] as [URL, RequestInit];
    expect(url.href).toBe('https://game.r2f2.com/config/IGC_MapList.xml');
    expect(init.headers).toMatchObject({ Authorization: `Basic ${Buffer.from('nextjs:read-only-password').toString('base64')}` });
    expect(init.cache).toBe('no-store');
  });

  it('fails closed on missing, insecure or malformed sources', async () => {
    await expect(getGameMapNames()).rejects.toThrow('GAME_CONFIG_URL is required');
    vi.stubEnv('GAME_CONFIG_URL', 'http://game.r2f2.com/config/');
    await expect(getGameMapNames()).rejects.toThrow();
    vi.stubEnv('GAME_CONFIG_URL', 'https://game.r2f2.com/config/');
    await expect(getGameMapNames()).rejects.toThrow('credentials');
    vi.stubEnv('GAME_CONFIG_BASIC_AUTH_USER', 'nextjs');
    vi.stubEnv('GAME_CONFIG_BASIC_AUTH_PASSWORD', 'read-only-password');
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('', { status: 401 })));
    await expect(getGameMapNames()).rejects.toThrow('401');
  });

  it('accepts only ws or wss URLs without embedded credentials', () => {
    vi.stubEnv('GAME_WEBSOCKET_URL', 'wss://game.r2f2.com/api/game');
    expect(gameWebSocketUrl()).toBe('wss://game.r2f2.com/api/game');
    vi.stubEnv('GAME_WEBSOCKET_URL', 'https://game.r2f2.com/api/game');
    expect(gameWebSocketUrl).toThrow();
    vi.stubEnv('GAME_WEBSOCKET_URL', 'wss://name:password@game.r2f2.com/api/game');
    expect(gameWebSocketUrl).toThrow();
  });
});
