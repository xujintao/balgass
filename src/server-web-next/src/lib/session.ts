import { cookies } from 'next/headers';
import type { Session } from '@supabase/supabase-js';
import { ApiError, providerError } from './errors';
import { config } from './config';
import { supabase } from './supabase';
export const ACCESS_COOKIE = 'balgass_access';
export const REFRESH_COOKIE = 'balgass_refresh';
export function bearerToken(headers: Headers): string | undefined {
  const value = headers.get('authorization');
  if (value === null) return;
  const match = /^Bearer ([^\s]+)$/i.exec(value);
  if (!match) throw new ApiError(401, 'UNAUTHENTICATED', '无效的身份令牌。');
  return match[1];
}
export function assertOrigin(headers: Headers, origin: string) {
  if (headers.get('origin') !== origin)
    throw new ApiError(403, 'INVALID_ORIGIN', '请求来源无效。');
}
export function isApp(request: Request) {
  return request.headers.get('x-client-type') === 'app';
}
export async function saveSession(session: Session) {
  const jar = await cookies();
  const opts = {
    httpOnly: true,
    secure: config().secure,
    sameSite: 'lax' as const,
    path: '/',
  };
  // Access cookie outlives the JWT so expired sessions can be refreshed on the next request.
  jar.set(ACCESS_COOKIE, session.access_token, {
    ...opts,
    maxAge: 60 * 60 * 24 * 30,
  });
  jar.set(REFRESH_COOKIE, session.refresh_token, {
    ...opts,
    maxAge: 60 * 60 * 24 * 30,
  });
}
export async function clearSession() {
  const jar = await cookies();
  jar.delete(ACCESS_COOKIE);
  jar.delete(REFRESH_COOKIE);
}
export async function context(request?: Request) {
  const bearer = request ? bearerToken(request.headers) : undefined;
  const jar = await cookies();
  if (request && isApp(request) && !bearer)
    throw new ApiError(401, 'UNAUTHENTICATED', 'App 请求需要 Bearer Token。');
  let token = bearer ?? jar.get(ACCESS_COOKIE)?.value;
  if (!token) throw new ApiError(401, 'UNAUTHENTICATED', '请先登录。');
  let client = supabase(token);
  let result = await client.auth.getUser(token);
  if (
    result.error &&
    ((result.error.status ?? 0) >= 500 ||
      result.error.name === 'AuthRetryableFetchError')
  )
    providerError(result.error);
  if (
    result.error &&
    !bearer &&
    request &&
    request.method !== 'GET' &&
    request.headers.get('origin') === config().origin
  ) {
    const refresh = jar.get(REFRESH_COOKIE)?.value;
    if (refresh) {
      const renewed = await supabase().auth.refreshSession({
        refresh_token: refresh,
      });
      if (
        renewed.error &&
        ((renewed.error.status ?? 0) >= 500 ||
          renewed.error.name === 'AuthRetryableFetchError')
      )
        providerError(renewed.error);
      if (!renewed.error && renewed.data.session) {
        await saveSession(renewed.data.session);
        token = renewed.data.session.access_token;
        client = supabase(token);
        result = await client.auth.getUser(token);
      }
    }
  }
  if (result.error || !result.data.user)
    throw new ApiError(401, 'UNAUTHENTICATED', '会话已失效，请重新登录。');
  if (!result.data.user.email_confirmed_at)
    throw new ApiError(403, 'EMAIL_UNVERIFIED', '请先验证邮箱。');
  return { client, token, user: result.data.user };
}
export async function sessionResponse(
  request: Request,
  session: Session | null,
) {
  if (!session) throw new ApiError(400, 'AUTH_FAILED', '未能建立会话。');
  if (!session.user.email_confirmed_at)
    throw new ApiError(403, 'EMAIL_UNVERIFIED', '请先验证邮箱。');
  if (isApp(request))
    return {
      user: { id: session.user.id, email: session.user.email },
      accessToken: session.access_token,
      refreshToken: session.refresh_token,
      expiresAt: session.expires_at,
      tokenType: 'bearer',
    };
  await saveSession(session);
  return {
    user: { id: session.user.id, email: session.user.email },
    expiresAt: session.expires_at,
  };
}
