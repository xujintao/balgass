import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';
import { z } from 'zod';
import { ApiError, providerError } from './errors';
import { config, requireCaptcha } from './config';
import { supabase } from './supabase';
import { findAuthUserByEmail } from './auth-user-lookup';
import {
  assertOrigin,
  bearerToken,
  clearSession,
  context,
  isApp,
  REFRESH_COOKIE,
  sessionResponse,
} from './session';
import { me, changeNickname } from './profile';
import { passkeyRequest } from './passkeys';
import { createGameAccount, listGameAccounts } from './game-accounts';
import { listItems, parseItemKind } from './item-catalog';
import { createOrder, listOrders } from './item-orders';
import { itemLocaleFromRequest } from './item-locale';
import {
  otpSend,
  otpVerify,
  refreshInput,
  nicknameInput,
  passkeyStart,
  authenticationVerify,
  registrationVerify,
} from './validation';
export const responseHeaders = {
  'Cache-Control': 'private, no-store',
  Vary: 'Cookie, Authorization, X-Client-Type, Accept-Language',
};
const routes: Record<string, string[]> = {
  'auth/otp/send': ['POST'],
  'auth/otp/verify': ['POST'],
  'auth/passkeys/login/options': ['POST'],
  'auth/passkeys/login/verify': ['POST'],
  'auth/passkeys/register/options': ['POST'],
  'auth/passkeys/register/verify': ['POST'],
  'auth/passkeys': ['GET'],
  'auth/session/refresh': ['POST'],
  'auth/logout': ['POST'],
  me: ['GET'],
  'me/nickname': ['PATCH'],
  'game/accounts': ['GET', 'POST'],
  items: ['GET'],
  orders: ['GET', 'POST'],
};
async function body(request: Request) {
  if (
    request.headers.get('content-type')?.split(';')[0].trim() !==
    'application/json'
  )
    throw new ApiError(415, 'UNSUPPORTED_MEDIA_TYPE', '请使用 JSON 请求。');
  if (Number(request.headers.get('content-length')) > 32768)
    throw new ApiError(413, 'PAYLOAD_TOO_LARGE', '请求过大。');
  const reader = request.body?.getReader();
  if (!reader) throw new ApiError(400, 'INVALID_INPUT', 'JSON 格式无效。');
  const chunks: Uint8Array[] = [];
  let length = 0;
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    length += value.length;
    if (length > 32768) {
      await reader.cancel();
      throw new ApiError(413, 'PAYLOAD_TOO_LARGE', '请求过大。');
    }
    chunks.push(value);
  }
  const bytes = new Uint8Array(length);
  let offset = 0;
  for (const chunk of chunks) {
    bytes.set(chunk, offset);
    offset += chunk.length;
  }
  try {
    return JSON.parse(new TextDecoder().decode(bytes));
  } catch {
    throw new ApiError(400, 'INVALID_INPUT', 'JSON 格式无效。');
  }
}
export async function dispatch(request: Request, path: string) {
  const method = request.method;
  const allowed =
    routes[path] ??
    (/^auth\/passkeys\/[0-9a-f-]+$/i.test(path) ? ['DELETE'] : undefined);
  if (!allowed) throw new ApiError(404, 'NOT_FOUND', '接口不存在。');
  if (!allowed.includes(method))
    throw new ApiError(405, 'METHOD_NOT_ALLOWED', '请求方法不支持。', {
      allowedMethods: allowed,
    });
  const authorization = bearerToken(request.headers);
  const clientType = request.headers.get('x-client-type');
  if (clientType && !['app', 'browser'].includes(clientType))
    throw new ApiError(400, 'INVALID_INPUT', '客户端类型无效。');
  if (authorization && !isApp(request))
    throw new ApiError(
      400,
      'INVALID_INPUT',
      'Bearer 请求需设置 X-Client-Type: app。',
    );
  if (!['GET', 'HEAD'].includes(method)) {
    if (!isApp(request) || request.headers.has('origin'))
      assertOrigin(request.headers, config().origin);
  }
  const client = supabase();
  if (path === 'auth/otp/send' && method === 'POST') {
    const input = otpSend.parse(await body(request));
    requireCaptcha(input.captchaToken);
    if (input.intent === 'signup') {
      const existing = await findAuthUserByEmail(input.email);
      if (existing?.email_confirmed_at)
        throw new ApiError(
          409,
          'EMAIL_ALREADY_REGISTERED',
          '该邮箱已注册，请登录。',
        );
    }
    const { error } = await client.auth.signInWithOtp({
      email: input.email,
      options: {
        shouldCreateUser: input.intent === 'signup',
        captchaToken: input.captchaToken,
      },
    });
    // Do not disclose whether a login address exists.
    if (error?.code !== 'otp_disabled' && error?.code !== 'user_not_found')
      providerError(error);
    return { sent: true, resendAfterSeconds: 60 };
  }
  if (path === 'auth/otp/verify' && method === 'POST') {
    const input = otpVerify.parse(await body(request));
    const { data, error } = await client.auth.verifyOtp({
      email: input.email,
      token: input.code,
      type: 'email',
    });
    providerError(error);
    return sessionResponse(request, data.session);
  }
  if (path === 'auth/passkeys/login/options' && method === 'POST') {
    const input = passkeyStart.parse(await body(request));
    requireCaptcha(input.captchaToken);
    const { data, error } = await client.auth.passkey.startAuthentication({
      options: { captchaToken: input.captchaToken },
    });
    providerError(error);
    return data;
  }
  if (path === 'auth/passkeys/login/verify' && method === 'POST') {
    const input = authenticationVerify.parse(await body(request));
    const { data, error } =
      await client.auth.passkey.verifyAuthentication(input);
    providerError(error);
    return sessionResponse(request, data?.session ?? null);
  }
  if (path === 'auth/session/refresh' && method === 'POST') {
    const refreshToken = isApp(request)
      ? refreshInput.parse(await body(request)).refreshToken
      : (await cookies()).get(REFRESH_COOKIE)?.value;
    if (!refreshToken) {
      if (!isApp(request)) await clearSession();
      throw new ApiError(401, 'UNAUTHENTICATED', '缺少刷新令牌。');
    }
    const { data, error } = await client.auth.refreshSession({
      refresh_token: refreshToken,
    });
    if (error) {
      if (
        (error.status ?? 0) >= 500 ||
        error.name === 'AuthRetryableFetchError'
      )
        providerError(error);
      if (!isApp(request)) await clearSession();
      throw new ApiError(401, 'UNAUTHENTICATED', '会话已失效，请重新登录。');
    }
    return sessionResponse(request, data.session);
  }
  if (path === 'auth/logout' && method === 'POST') {
    if (isApp(request)) {
      const ctx = await context(request);
      providerError(
        (await client.auth.admin.signOut(ctx.token, 'local')).error,
      );
    } else {
      // Refresh if possible to revoke even an expired access session.
      const refresh = (await cookies()).get(REFRESH_COOKIE)?.value;
      if (refresh) {
        const result = await client.auth.refreshSession({
          refresh_token: refresh,
        });
        if (
          result.error &&
          ((result.error.status ?? 0) >= 500 ||
            result.error.name === 'AuthRetryableFetchError')
        )
          providerError(result.error);
        if (result.data.session)
          providerError(
            (
              await client.auth.admin.signOut(
                result.data.session.access_token,
                'local',
              )
            ).error,
          );
      }
      await clearSession();
    }
    return { signedOut: true };
  }
  if (path === 'items' && method === 'GET') {
    const kind = new URL(request.url).searchParams.get('kind') ?? 'sword';
    return listItems(parseItemKind(kind), itemLocaleFromRequest(request));
  }
  const ctx = await context(request);
  if (path === 'orders' && method === 'GET') return listOrders(ctx, itemLocaleFromRequest(request));
  if (path === 'orders' && method === 'POST') return createOrder(ctx, await body(request), itemLocaleFromRequest(request));
  if (path === 'game/accounts' && method === 'GET') return listGameAccounts(ctx);
  if (path === 'game/accounts' && method === 'POST')
    return createGameAccount(ctx, await body(request));
  if (path === 'me' && method === 'GET') return me(ctx);
  if (path === 'me/nickname' && method === 'PATCH')
    return changeNickname(
      ctx,
      nicknameInput.parse(await body(request)).nickname,
    );
  if (path === 'auth/passkeys' && method === 'GET')
    return passkeyRequest('', 'GET', ctx.token);
  if (path === 'auth/passkeys/register/options' && method === 'POST')
    return passkeyRequest('/registration/options', 'POST', ctx.token, {});
  if (path === 'auth/passkeys/register/verify' && method === 'POST') {
    const input = registrationVerify.parse(await body(request));
    return passkeyRequest('/registration/verify', 'POST', ctx.token, {
      challenge_id: input.challengeId,
      credential: input.credential,
    });
  }
  if (/^auth\/passkeys\/[^/]+$/.test(path) && method === 'DELETE') {
    const id = z.uuid().parse(path.split('/').at(-1));
    return passkeyRequest(`/${id}`, 'DELETE', ctx.token);
  }
  throw new ApiError(404, 'NOT_FOUND', '接口不存在。');
}
export async function handle(request: Request, path: string) {
  try {
    return NextResponse.json(
      { data: await dispatch(request, path) },
      { headers: responseHeaders },
    );
  } catch (error) {
    const failure =
      error instanceof ApiError
        ? error
        : error instanceof z.ZodError
          ? new ApiError(400, 'INVALID_INPUT', '输入格式无效。')
          : new ApiError(
              503,
              'SERVICE_UNAVAILABLE',
              '服务暂不可用，请稍后重试。',
            );
    return NextResponse.json(
      {
        error: {
          code: failure.code,
          message: failure.message,
          ...(failure.details ?? {}),
        },
      },
      {
        status: failure.status,
        headers: {
          ...responseHeaders,
          ...(failure.status === 429 ? { 'Retry-After': '60' } : {}),
          ...(failure.status === 405
            ? {
                Allow: (failure.details?.allowedMethods as string[]).join(', '),
              }
            : {}),
        },
      },
    );
  }
}
