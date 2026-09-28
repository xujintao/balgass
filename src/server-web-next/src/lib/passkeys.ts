// Isolate experimental Supabase wire endpoints here. Authenticated endpoints accept
// an access token directly, so App clients never need to surrender a refresh token.
import { config } from './config';
import { ApiError, providerError } from './errors';
export async function passkeyRequest(
  path: string,
  method: string,
  token: string,
  body?: unknown,
) {
  const { url, key } = config();
  const response = await fetch(`${url}/auth/v1/passkeys${path}`, {
    method,
    headers: {
      apikey: key,
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: body === undefined ? undefined : JSON.stringify(body),
    cache: 'no-store',
    signal: AbortSignal.timeout(15000),
  });
  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    providerError({
      code: error.code ?? error.error_code,
      status: response.status,
    });
  }
  if (method === 'DELETE') return { deleted: true };
  const result = await response.json();
  if (!result)
    throw new ApiError(502, 'AUTH_UNAVAILABLE', '认证服务响应无效。');
  return result;
}
