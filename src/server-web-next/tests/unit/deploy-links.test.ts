import { afterEach, expect, it, vi } from 'vitest';
import { deployLinks, deployUrl } from '../../src/lib/deploy-links';

afterEach(() => vi.unstubAllEnvs());

it('accepts HTTPS and loopback HTTP console URLs but rejects unsafe URLs', () => {
  expect(deployUrl('https://vercel.com/team/project')).toBe('https://vercel.com/team/project');
  expect(deployUrl('http://localhost:3000/')).toBe('http://localhost:3000/');
  expect(deployUrl('http://127.0.0.1:3000/')).toBe('http://127.0.0.1:3000/');
  expect(deployUrl('http://[::1]:3000/')).toBe('http://[::1]:3000/');
  expect(deployUrl('http://vercel.com/team/project')).toBeUndefined();
  expect(deployUrl('http://localhost.example.com/')).toBeUndefined();
  expect(deployUrl('javascript:alert(1)')).toBeUndefined();
  expect(deployUrl('https://user:password@example.com/')).toBeUndefined();
  expect(deployUrl('not a url')).toBeUndefined();
  expect(deployUrl(undefined)).toBeUndefined();
});

it('reads each configured deployment destination independently', () => {
  vi.stubEnv('DEPLOY_VERCEL_URL', 'http://localhost:3000/');
  vi.stubEnv('DEPLOY_SUPABASE_URL', 'http://supabase.com/dashboard');
  vi.stubEnv('DEPLOY_RESEND_URL', 'https://resend.com/overview');
  vi.stubEnv('DEPLOY_GRAFANA_URL', 'https://example.grafana.net/');
  vi.stubEnv('DEPLOY_PGADMIN_GAME_URL', 'https://game-pgadmin.example/');
  vi.stubEnv('DEPLOY_PGADMIN_US_URL', '');
  expect(deployLinks()).toEqual({
    vercel: 'http://localhost:3000/',
    supabase: undefined,
    resend: 'https://resend.com/overview',
    grafana: 'https://example.grafana.net/',
    pgadminGame: 'https://game-pgadmin.example/',
    pgadminUs: undefined,
  });
  vi.stubEnv('DEPLOY_PGADMIN_GAME_URL', '');
  vi.stubEnv('DEPLOY_PGADMIN_US_URL', 'https://us-pgadmin.example/');
  expect(deployLinks().pgadminGame).toBeUndefined();
  expect(deployLinks().pgadminUs).toBe('https://us-pgadmin.example/');
});
