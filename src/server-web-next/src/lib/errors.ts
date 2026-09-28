export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
    public details?: Record<string, unknown>,
  ) {
    super(message);
  }
}
export function providerError(
  error: { code?: string; status?: number; name?: string } | null,
): void {
  if (!error) return;
  const code = error.code ?? 'AUTH_FAILED';
  if (error.status === 429 || code.includes('rate_limit'))
    throw new ApiError(429, 'RATE_LIMITED', '操作过于频繁，请稍后重试。');
  if (code === 'otp_expired')
    throw new ApiError(400, 'OTP_INVALID', '验证码无效或已过期。');
  if (code === 'captcha_failed')
    throw new ApiError(400, 'CAPTCHA_FAILED', '请重新完成人机验证。');
  if (
    [
      'session_not_found',
      'refresh_token_not_found',
      'refresh_token_already_used',
      'bad_jwt',
    ].includes(code)
  )
    throw new ApiError(401, 'UNAUTHENTICATED', '会话已失效，请重新登录。');
  if (code === 'passkey_disabled')
    throw new ApiError(
      503,
      'PASSKEY_UNAVAILABLE',
      'Passkey 暂不可用，请使用邮件验证码。',
    );
  if (code.startsWith('webauthn_'))
    throw new ApiError(
      400,
      code.toUpperCase(),
      'Passkey 验证失败，请重试或使用邮件验证码。',
    );
  if ((error.status ?? 0) >= 500 || error.name === 'AuthRetryableFetchError')
    throw new ApiError(
      502,
      'AUTH_UNAVAILABLE',
      '认证服务暂不可用，请稍后重试。',
    );
  throw new ApiError(400, 'AUTH_FAILED', '认证失败，请检查输入后重试。');
}
