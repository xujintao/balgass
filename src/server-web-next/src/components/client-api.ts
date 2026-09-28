'use client';
export class ClientError extends Error {
  constructor(
    public code: string,
    message: string,
    public nextNicknameChangeAt?: string,
  ) {
    super(message);
  }
}
export async function api<T = unknown>(
  path: string,
  method = 'GET',
  body?: unknown,
  retry = true,
): Promise<T> {
  const response = await fetch(`/api/v1/${path}`, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
    cache: 'no-store',
  });
  if (
    response.status === 401 &&
    retry &&
    ![
      'auth/session/refresh',
      'auth/logout',
      'auth/otp/verify',
      'auth/passkeys/login/verify',
    ].includes(path)
  ) {
    await api('auth/session/refresh', 'POST', undefined, false);
    return api(path, method, body, false);
  }
  const result = await response.json();
  if (!response.ok)
    throw new ClientError(
      result.error.code,
      result.error.message,
      result.error.nextNicknameChangeAt,
    );
  return result.data as T;
}
