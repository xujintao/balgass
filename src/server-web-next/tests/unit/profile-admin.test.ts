import { beforeEach, expect, it, vi } from 'vitest';

const state = vi.hoisted(() => ({ rows: new Map<string, Record<string, unknown>>() }));
vi.mock('../../src/lib/supabase-admin', () => ({
  supabaseAdmin: () => ({
    from: () => ({
      select: () => ({
        eq: (field: string, value: string) => ({
          maybeSingle: async () => ({ data: state.rows.get(value) ?? null, error: null }),
          limit: async () => ({
            data: [...state.rows.values()].filter((row) => row[field] === value).slice(0, 1),
            error: null,
          }),
        }),
      }),
      insert: (values: Record<string, unknown>) => ({
        select: () => ({ single: async () => {
          const id = values.user_id as string;
          if (state.rows.has(id)) return { data: null, error: { code: '23505' } };
          const row = { ...values, created_at: '2026-01-01T00:00:00Z', nickname_changed_at: null };
          state.rows.set(id, row);
          return { data: row, error: null };
        } }),
      }),
    }),
  }),
}));

import { ensureProfile } from '../../src/lib/profile';
beforeEach(() => state.rows.clear());

it('空库中首位已验证用户成为管理员，后续用户为普通用户', async () => {
  await ensureProfile({ id: 'first', email_confirmed_at: '2026-01-01T00:00:00Z' });
  await ensureProfile({ id: 'second', email_confirmed_at: '2026-01-01T00:00:01Z' });
  expect(state.rows.get('first')?.role).toBe('admin');
  expect(state.rows.get('second')?.role).toBe('player');
});

it('未验证邮箱不能建立资料或获得管理员角色', async () => {
  await expect(ensureProfile({ id: 'unverified' })).rejects.toMatchObject({ status: 403 });
  expect(state.rows.size).toBe(0);
});
