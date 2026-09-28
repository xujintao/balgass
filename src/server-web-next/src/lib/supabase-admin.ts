import 'server-only';
import { createClient } from '@supabase/supabase-js';
import { config } from './config';
import { ApiError } from './errors';

// This privileged client is used only after server-side identity verification.
export function supabaseAdmin() {
  const { url } = config();
  const key =
    process.env.SUPABASE_SECRET_KEY ?? process.env.SUPABASE_SERVICE_ROLE_KEY;
  if (!key)
    throw new ApiError(
      503,
      'CONFIGURATION_REQUIRED',
      '请配置服务端 Supabase Secret Key。',
    );
  return createClient(url, key, {
    auth: {
      persistSession: false,
      autoRefreshToken: false,
      detectSessionInUrl: false,
    },
    global: {
      fetch: (input, init) =>
        fetch(input, {
          ...init,
          cache: 'no-store',
          signal: init?.signal ?? AbortSignal.timeout(15000),
        }),
    },
  });
}
