import { ApiError } from './errors';
import type { context } from './session';
import { supabaseAdmin } from './supabase-admin';
import {
  profileService,
  type ProfileRepository,
  type ProfileRow,
  type ProfileUser,
} from './profile-service';

function service() {
  const client = supabaseAdmin();
  const repository: ProfileRepository = {
    async find(userId) {
      const { data, error } = await client
        .from('profiles')
        .select('*')
        .eq('user_id', userId)
        .maybeSingle();
      if (error)
        throw new ApiError(
          503,
          'PROFILE_UNAVAILABLE',
          '暂时无法读取用户资料。',
        );
      return data as ProfileRow | null;
    },
    async insert(userId, nickname, key) {
      const { data, error } = await client
        .from('profiles')
        .insert({ user_id: userId, nickname, nickname_key: key })
        .select()
        .single();
      if (error?.code === '23505') return null;
      if (error)
        throw new ApiError(
          503,
          'PROFILE_UNAVAILABLE',
          '暂时无法创建用户资料。',
        );
      return data as ProfileRow;
    },
    async update(observed, nickname, key, changedAt, cutoff) {
      let query = client
        .from('profiles')
        .update({ nickname, nickname_key: key, nickname_changed_at: changedAt })
        .eq('user_id', observed.user_id)
        .eq('nickname', observed.nickname)
        .or(`nickname_changed_at.is.null,nickname_changed_at.lte.${cutoff}`);
      query =
        observed.nickname_changed_at === null
          ? query.is('nickname_changed_at', null)
          : query.eq('nickname_changed_at', observed.nickname_changed_at);
      const { data, error } = await query.select().maybeSingle();
      if (error?.code === '23505')
        throw new ApiError(409, 'NICKNAME_TAKEN', '昵称已被使用。');
      if (error)
        throw new ApiError(503, 'PROFILE_UNAVAILABLE', '暂时无法修改昵称。');
      return data as ProfileRow | null;
    },
  };
  return profileService(repository);
}
export async function ensureProfile(user: ProfileUser) {
  return service().ensure(user);
}
export async function me(ctx: Awaited<ReturnType<typeof context>>) {
  return service().me(ctx.user);
}
export async function changeNickname(
  ctx: Awaited<ReturnType<typeof context>>,
  value: string,
) {
  return service().change(ctx.user, value);
}
