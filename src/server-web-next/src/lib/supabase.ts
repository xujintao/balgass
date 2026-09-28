import { createClient } from '@supabase/supabase-js';
import { config } from './config';
export function supabase(token?: string) {
  const { url, key } = config();
  return createClient(url, key, {
    auth: {
      persistSession: false,
      autoRefreshToken: false,
      detectSessionInUrl: false,
      experimental: { passkey: true },
    },
    global: {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      fetch: (input, init) =>
        fetch(input, {
          ...init,
          cache: 'no-store',
          signal: init?.signal ?? AbortSignal.timeout(15000),
        }),
    },
  });
}
