import 'server-only';
import { ApiError } from './errors';
import { supabaseAdmin } from './supabase-admin';
import type { context } from './session';

type Context = Awaited<ReturnType<typeof context>>;
export type UserRole = 'player' | 'admin';
export type AdminUser = {
  id: string;
  email: string;
  nickname: string;
  role: UserRole;
  createdAt: string;
};
type Profile = { user_id: string; nickname: string; role: UserRole; created_at: string };
export const ADMIN_PAGE_SIZE = 20;

export async function isAdmin(userId: string): Promise<boolean> {
  const { data, error } = await supabaseAdmin().from('profiles')
    .select('role').eq('user_id', userId).maybeSingle();
  if (error) throw new ApiError(503, 'USERS_UNAVAILABLE', '用户资料暂不可用。');
  return data?.role === 'admin';
}

export async function assertAdmin(ctx: Context) {
  if (!await isAdmin(ctx.user.id))
    throw new ApiError(403, 'FORBIDDEN', '仅管理员可以访问。');
}

export async function listAdminUsers(ctx: Context, page: number) {
  await assertAdmin(ctx);
  if (!Number.isSafeInteger(page) || page < 1 || page > 100000)
    throw new ApiError(400, 'INVALID_INPUT', '页码无效。');
  const client = supabaseAdmin();
  const start = (page - 1) * ADMIN_PAGE_SIZE;
  const { data, error, count } = await client.from('profiles')
    .select('user_id,nickname,role,created_at', { count: 'exact' })
    .order('created_at', { ascending: true })
    .order('user_id', { ascending: true })
    .range(start, start + ADMIN_PAGE_SIZE - 1);
  if (error || count === null)
    throw new ApiError(503, 'USERS_UNAVAILABLE', '用户列表暂不可用。');
  const users = await Promise.all((data as Profile[]).map(async (profile) => {
    const { data: auth, error: authError } = await client.auth.admin.getUserById(profile.user_id);
    if (authError) throw new ApiError(502, 'AUTH_UNAVAILABLE', '认证服务暂不可用。');
    const user = auth.user;
    if (!user?.email || !user.email_confirmed_at) return null;
    return {
      id: profile.user_id,
      email: user.email,
      nickname: profile.nickname,
      role: profile.role,
      createdAt: profile.created_at,
    } satisfies AdminUser;
  }));
  return { users: users.filter((user): user is AdminUser => user !== null), page,
    pageSize: ADMIN_PAGE_SIZE, total: count, totalPages: Math.ceil(count / ADMIN_PAGE_SIZE) };
}

export async function changeUserRole(ctx: Context, userId: string, role: UserRole) {
  await assertAdmin(ctx);
  const client = supabaseAdmin();
  const { data: target, error: lookupError } = await client.from('profiles')
    .select('user_id,role').eq('user_id', userId).maybeSingle();
  if (lookupError) throw new ApiError(503, 'USERS_UNAVAILABLE', '角色修改失败，请重试。');
  if (!target) throw new ApiError(404, 'USER_NOT_FOUND', '用户不存在。');
  if (target.role === role) return { id: userId, role };
  if (target.role === 'admin' && role === 'player') {
    const { count, error: countError } = await client.from('profiles')
      .select('user_id', { count: 'exact', head: true }).eq('role', 'admin');
    if (countError || count === null)
      throw new ApiError(503, 'USERS_UNAVAILABLE', '角色修改失败，请重试。');
    if (count <= 1)
      throw new ApiError(409, 'LAST_ADMIN', '不能撤销最后一位管理员。');
  }
  const { data, error } = await client.from('profiles').update({ role })
    .eq('user_id', userId).eq('role', target.role).select('user_id,role').maybeSingle();
  if (error) throw new ApiError(503, 'USERS_UNAVAILABLE', '角色修改失败，请重试。');
  if (!data) throw new ApiError(409, 'ROLE_CONFLICT', '角色已变化，请刷新后重试。');
  return { id: data.user_id as string, role: data.role as UserRole };
}
