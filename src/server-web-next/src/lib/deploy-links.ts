export type DeployTarget = 'vercel' | 'supabase' | 'resend' | 'grafana' | 'pgadminGame' | 'pgadminUs';

export function deployUrl(value: string | undefined): string | undefined {
  if (!value) return;
  try {
    const url = new URL(value);
    const localHttp = url.protocol === 'http:' &&
      ['localhost', '127.0.0.1', '[::1]'].includes(url.hostname);
    if ((!localHttp && url.protocol !== 'https:') || !url.hostname || url.username || url.password)
      return;
    return url.href;
  } catch {
    return;
  }
}

export function deployLinks(): Record<DeployTarget, string | undefined> {
  return {
    vercel: deployUrl(process.env.DEPLOY_VERCEL_URL),
    supabase: deployUrl(process.env.DEPLOY_SUPABASE_URL),
    resend: deployUrl(process.env.DEPLOY_RESEND_URL),
    grafana: deployUrl(process.env.DEPLOY_GRAFANA_URL),
    pgadminGame: deployUrl(process.env.DEPLOY_PGADMIN_GAME_URL),
    pgadminUs: deployUrl(process.env.DEPLOY_PGADMIN_US_URL),
  };
}
