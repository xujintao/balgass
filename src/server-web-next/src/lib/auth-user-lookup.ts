import { ApiError } from './errors';
import { config } from './config';

type AuthUser = { email?: string | null; email_confirmed_at?: string | null };

// This Auth admin endpoint is called only by the server. Never return its user list.
export async function findAuthUserByEmail(
  email: string,
): Promise<AuthUser | null> {
  const { url } = config();
  const key =
    process.env.SUPABASE_SECRET_KEY ?? process.env.SUPABASE_SERVICE_ROLE_KEY;
  if (!key)
    throw new ApiError(
      503,
      'CONFIGURATION_REQUIRED',
      '请配置服务端 Supabase Secret Key。',
    );

  const pageSize = 100;
  for (let page = 1; page <= 20; page++) {
    const endpoint = new URL(`${url}/auth/v1/admin/users`);
    endpoint.searchParams.set('filter', email);
    endpoint.searchParams.set('page', String(page));
    endpoint.searchParams.set('per_page', String(pageSize));
    let response: Response;
    try {
      response = await fetch(endpoint, {
        headers: { apikey: key, Authorization: `Bearer ${key}` },
        cache: 'no-store',
        signal: AbortSignal.timeout(15000),
      });
    } catch {
      throw new ApiError(
        502,
        'AUTH_UNAVAILABLE',
        '认证服务暂不可用，请稍后重试。',
      );
    }
    if (!response.ok)
      throw new ApiError(
        502,
        'AUTH_UNAVAILABLE',
        '认证服务暂不可用，请稍后重试。',
      );
    const data: unknown = await response.json().catch(() => null);
    if (
      !data ||
      typeof data !== 'object' ||
      !('users' in data) ||
      !Array.isArray(data.users)
    )
      throw new ApiError(502, 'AUTH_UNAVAILABLE', '认证服务响应无效。');
    const users = data.users as AuthUser[];
    const match = users.find((user) => user.email?.toLowerCase() === email);
    if (match) return match;
    if (users.length < pageSize) return null;
  }
  throw new ApiError(
    502,
    'AUTH_UNAVAILABLE',
    '认证服务搜索结果过多，请稍后重试。',
  );
}
