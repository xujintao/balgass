import { ApiError } from './errors';
export function config() {
  const url = process.env.SUPABASE_URL;
  const key = process.env.SUPABASE_PUBLISHABLE_KEY;
  const origin = process.env.APP_ORIGIN;
  if (!url || !key || !origin)
    throw new ApiError(503, 'CONFIGURATION_REQUIRED', '请先配置网站认证服务。');
  const parsed = new URL(origin);
  if (
    parsed.origin !== origin ||
    (parsed.protocol !== 'https:' &&
      !['localhost', '127.0.0.1', '[::1]'].includes(parsed.hostname))
  )
    throw new ApiError(503, 'CONFIGURATION_REQUIRED', '网站地址配置无效。');
  return { url, key, origin, secure: parsed.protocol === 'https:' };
}
export function requireCaptcha(token?: string) {
  const disabled = process.env.AUTH_CAPTCHA_DISABLED === 'true';
  const origin = new URL(config().origin);
  if (disabled && ['localhost', '127.0.0.1', '[::1]'].includes(origin.hostname))
    return;
  if (!token) throw new ApiError(400, 'CAPTCHA_REQUIRED', '请完成人机验证。');
}
