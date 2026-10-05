import 'server-only';
import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { unstable_cache } from 'next/cache';
import { XMLParser, XMLValidator } from 'fast-xml-parser';

const MAX_CONFIG_BYTES = 2 * 1024 * 1024;
const mapListFile = 'IGC_MapList.xml';

export function parseGameMapNames(xml: string): string[] {
  if (XMLValidator.validate(xml) !== true || /<!DOCTYPE/i.test(xml)) {
    throw new Error('Invalid game map XML');
  }
  const document = new XMLParser({
    ignoreAttributes: false,
    attributeNamePrefix: '@_',
    processEntities: false,
  }).parse(xml);
  const maps = document?.MapList?.DefaultMaps?.Map;
  const entries = Array.isArray(maps) ? maps : maps && typeof maps === 'object' ? [maps] : null;
  if (!entries) throw new Error('Invalid game map list');

  const counts = new Map<string, number>();
  const names: string[] = [];
  for (const map of entries) {
    const file = map?.['@_File'];
    if (typeof file !== 'string' || !file.endsWith('.att')) {
      throw new Error('Invalid game map entry');
    }
    const base = /^\d+_(.+)\.att$/.exec(file)?.[1] ?? file.slice(0, -4);
    if (base === 'Null') continue;
    const count = counts.get(base) ?? 0;
    names.push(count === 0 ? base : `${base}${count + 1}`);
    counts.set(base, count + 1);
  }
  return names;
}

function configUrl(): URL {
  const value = process.env.GAME_CONFIG_URL;
  if (!value) throw new Error('GAME_CONFIG_URL is required');
  let url: URL;
  try {
    url = new URL(value);
  } catch {
    throw new Error('GAME_CONFIG_URL is invalid');
  }
  if (
    !['file:', 'https:'].includes(url.protocol) ||
    !url.pathname.endsWith('/') ||
    url.username || url.password || url.search || url.hash
  ) {
    throw new Error('GAME_CONFIG_URL must be a file or HTTPS directory URL');
  }
  return url;
}

async function readGameConfig(base: URL, filename: typeof mapListFile): Promise<string> {
  const url = new URL(filename, base);
  let xml: string;
  if (base.protocol === 'file:') {
    const contents = await readFile(fileURLToPath(url));
    if (contents.byteLength > MAX_CONFIG_BYTES) throw new Error('Game config is too large');
    xml = contents.toString('utf8');
  } else {
    const user = process.env.GAME_CONFIG_BASIC_AUTH_USER;
    const password = process.env.GAME_CONFIG_BASIC_AUTH_PASSWORD;
    if (!user || !password) throw new Error('Game config credentials are required');
    const authorization = Buffer.from(`${user}:${password}`).toString('base64');
    const response = await fetch(url, {
      headers: { Authorization: `Basic ${authorization}` },
      cache: 'no-store',
      redirect: 'error',
      signal: AbortSignal.timeout(10000),
    });
    if (!response.ok) throw new Error(`Game config request failed: ${response.status}`);
    const length = Number(response.headers.get('content-length'));
    if (length > MAX_CONFIG_BYTES) throw new Error('Game config is too large');
    const contents = await response.arrayBuffer();
    if (contents.byteLength > MAX_CONFIG_BYTES) throw new Error('Game config is too large');
    xml = new TextDecoder().decode(contents);
  }
  return xml;
}

async function loadGameMapNames(baseUrl: string): Promise<string[]> {
  return parseGameMapNames(await readGameConfig(new URL(baseUrl), mapListFile));
}

const cachedGameMapNames = unstable_cache(loadGameMapNames, ['game-map-names'], {
  revalidate: 60,
});

export async function getGameMapNames(): Promise<string[]> {
  const base = configUrl();
  return base.protocol === 'file:'
    ? loadGameMapNames(base.href)
    : cachedGameMapNames(base.href);
}

export function gameWebSocketUrl(): string {
  const value = process.env.GAME_WEBSOCKET_URL;
  if (!value) throw new Error('GAME_WEBSOCKET_URL is required');
  let url: URL;
  try {
    url = new URL(value);
  } catch {
    throw new Error('GAME_WEBSOCKET_URL is invalid');
  }
  if (
    !['ws:', 'wss:'].includes(url.protocol) ||
    url.username || url.password || url.search || url.hash
  ) {
    throw new Error('GAME_WEBSOCKET_URL is invalid');
  }
  return url.href;
}
