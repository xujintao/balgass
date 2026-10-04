import { z } from 'zod';
import { ApiError } from './errors';
import type { context } from './session';

type UserContext = Awaited<ReturnType<typeof context>>;
const accountInput = z.object({
  name: z.string().min(1).max(10).regex(/^[\x21-\x7e]+$/),
  password: z.string().min(1).max(10),
  passwordConfirmation: z.string(),
}).strict().refine((value) => value.password === value.passwordConfirmation);
const character = z.object({ name: z.string(), level: z.number(), master_level: z.number() });
const account = z.object({ name: z.string(), user_email: z.string(), characters: z.array(character).nullable() });

function endpoint() {
  const value = process.env.GAME_API_URL;
  if (!value) throw new ApiError(503, 'GAME_UNAVAILABLE', '游戏账号服务未配置。');
  let url: URL;
  try { url = new URL(value); } catch { throw new ApiError(503, 'GAME_UNAVAILABLE', '游戏账号服务配置无效。'); }
  const local = ['localhost', '127.0.0.1', '[::1]'].includes(url.hostname);
  if ((url.protocol !== 'https:' && !(process.env.NODE_ENV !== 'production' && url.protocol === 'http:' && local)) ||
    url.pathname !== '/api/command' || url.username || url.password || url.search || url.hash)
    throw new ApiError(503, 'GAME_UNAVAILABLE', '游戏账号服务配置无效。');
  return url;
}

function gameAPIToken() {
  const token = process.env.GAME_API_TOKEN;
  if (!token || token.length < 32 || /\s/.test(token))
    throw new ApiError(503, 'GAME_UNAVAILABLE', '游戏账号服务未配置。');
  return token;
}

function email(ctx: UserContext) {
  if (!ctx.user.email || !ctx.user.email_confirmed_at)
    throw new ApiError(403, 'EMAIL_UNVERIFIED', '请先验证邮箱。');
  return ctx.user.email;
}

async function gameRequest(action: 'GetAccountList' | 'CreateAccount', input: Record<string, string>) {
  const url = endpoint();
  const token = gameAPIToken();
  let response: Response;
  try {
    response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ action, in: input }),
      cache: 'no-store',
      signal: AbortSignal.timeout(10000),
    });
  } catch {
    throw new ApiError(502, 'GAME_UNAVAILABLE', '游戏账号服务暂不可用。');
  }
  if (!response.ok) {
    if (response.status === 429) throw new ApiError(429, 'RATE_LIMITED', '操作过于频繁，请稍后重试。');
    if (response.status === 409) throw new ApiError(409, 'GAME_ACCOUNT_EXISTS', '游戏账号名已被使用。');
    throw new ApiError(502, 'GAME_UNAVAILABLE', '游戏账号服务暂不可用。');
  }
  try {
    const data = z.object({ action: z.literal(action), out: z.unknown() }).safeParse(await response.json());
    if (data.success) return data.data.out;
    throw new Error('invalid command response');
  } catch {
    throw new ApiError(502, 'GAME_UNAVAILABLE', '游戏账号服务返回了无效数据。');
  }
}

export async function listGameAccounts(ctx: UserContext) {
  const owner = email(ctx);
  const data = z.array(account).nullable().safeParse(await gameRequest('GetAccountList', { user_email: owner }));
  if (!data.success || data.data?.some((item) => item.user_email !== owner))
    throw new ApiError(502, 'GAME_UNAVAILABLE', '游戏账号服务返回了无效数据。');
  return (data.data ?? []).map((item) => ({
    name: item.name,
    characters: (item.characters ?? []).map((entry) => ({ name: entry.name, level: entry.level + entry.master_level })),
  }));
}

export async function createGameAccount(ctx: UserContext, input: unknown) {
  const owner = email(ctx);
  const parsed = accountInput.safeParse(input);
  if (!parsed.success) throw new ApiError(400, 'INVALID_INPUT', '请检查账号名和密码，确认两次密码一致。');
  const result = await gameRequest('CreateAccount', {
    name: parsed.data.name, password: parsed.data.password, user_email: owner,
  });
  const created = z.object({ name: z.string(), user_email: z.string() }).safeParse(result);
  if (!created.success || created.data.user_email !== owner || created.data.name !== parsed.data.name)
    throw new ApiError(502, 'GAME_UNAVAILABLE', '游戏账号服务返回了无效数据。');
  return { name: created.data.name };
}
