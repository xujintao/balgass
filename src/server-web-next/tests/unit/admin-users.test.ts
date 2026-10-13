import { beforeEach, expect, it, vi } from 'vitest';

vi.mock('server-only', () => ({}));
const state = vi.hoisted(() => ({ roles: new Map<string, 'player' | 'admin'>() }));
vi.mock('../../src/lib/supabase-admin', () => ({
  supabaseAdmin: () => ({
    from: () => ({
      select: (columns: string, options?: { head?: boolean }) => query(columns, options),
      update: (values: { role: 'player' | 'admin' }) => query('', undefined, values.role),
    }),
  }),
}));

function query(columns: string, options?: { head?: boolean }, nextRole?: 'player' | 'admin') {
  const filters = new Map<string, string>();
  const result = {
    eq(column: string, value: string) { filters.set(column, value); return result; },
    select() { return result; },
    async maybeSingle() {
      const id = filters.get('user_id') ?? '';
      const role = state.roles.get(id);
      if (!role || (nextRole && role !== filters.get('role')))
        return { data: null, error: null };
      if (nextRole) state.roles.set(id, nextRole);
      return { data: columns === 'role' ? { role } : { user_id: id, role: nextRole ?? role }, error: null };
    },
    then(resolve: (value: unknown) => unknown) {
      const count = [...state.roles.values()].filter((role) => role === filters.get('role')).length;
      return Promise.resolve(resolve({ count: options?.head ? count : null, error: null }));
    },
  };
  return result;
}

import { assertAdmin, changeUserRole, isAdmin } from '../../src/lib/admin-users';
const ctx = { user: { id: 'actor' } } as Parameters<typeof assertAdmin>[0];

beforeEach(() => {
  state.roles.clear();
  state.roles.set('actor', 'player');
  state.roles.set('target', 'player');
});

it('普通用户不能读取管理数据或修改角色', async () => {
  expect(await isAdmin('actor')).toBe(false);
  await expect(assertAdmin(ctx)).rejects.toMatchObject({ status: 403, code: 'FORBIDDEN' });
  await expect(changeUserRole(ctx, 'target', 'admin')).rejects.toMatchObject({ status: 403 });
  expect(state.roles.get('target')).toBe('player');
});

it('顺序操作时阻止撤销最后一位管理员', async () => {
  state.roles.set('actor', 'admin');
  await expect(changeUserRole(ctx, 'actor', 'player')).rejects.toMatchObject({ status: 409, code: 'LAST_ADMIN' });
  expect(await changeUserRole(ctx, 'target', 'admin')).toEqual({ id: 'target', role: 'admin' });
  expect(await changeUserRole(ctx, 'actor', 'player')).toEqual({ id: 'actor', role: 'player' });
  expect(state.roles.get('actor')).toBe('player');
  await expect(changeUserRole(ctx, 'target', 'player')).rejects.toMatchObject({ status: 403, code: 'FORBIDDEN' });
});

it('目标不存在时返回 404', async () => {
  state.roles.set('actor', 'admin');
  await expect(changeUserRole(ctx, 'missing', 'admin')).rejects.toMatchObject({ status: 404, code: 'USER_NOT_FOUND' });
});
