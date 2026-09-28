import { ApiError } from './errors';
import { context } from './session';
export async function me(ctx: Awaited<ReturnType<typeof context>>) {
  const { data, error } = await ctx.client
    .from('profiles')
    .select('nickname,created_at,nickname_changed_at')
    .eq('user_id', ctx.user.id)
    .single();
  if (error)
    throw new ApiError(503, 'PROFILE_UNAVAILABLE', '暂时无法读取用户资料。');
  return {
    id: ctx.user.id,
    email: ctx.user.email,
    nickname: data.nickname as string,
    createdAt: data.created_at as string,
    nicknameChangedAt: data.nickname_changed_at as string | null,
    nextNicknameChangeAt: data.nickname_changed_at
      ? new Date(
          new Date(data.nickname_changed_at).getTime() + 30 * 86400000,
        ).toISOString()
      : null,
  };
}
export async function changeNickname(
  ctx: Awaited<ReturnType<typeof context>>,
  value: string,
) {
  const { error } = await ctx.client.rpc('change_nickname', {
    requested_nickname: value,
  });
  if (error) {
    if (error.message === 'NICKNAME_TAKEN')
      throw new ApiError(409, 'NICKNAME_TAKEN', '昵称已被使用。');
    if (error.message === 'NICKNAME_COOLDOWN')
      throw new ApiError(
        409,
        'NICKNAME_COOLDOWN',
        '每 30 天只能修改一次昵称。',
        { nextNicknameChangeAt: error.details },
      );
    if (error.message === 'INVALID_NICKNAME')
      throw new ApiError(400, 'INVALID_NICKNAME', '昵称格式无效。');
    throw new ApiError(503, 'PROFILE_UNAVAILABLE', '暂时无法修改昵称。');
  }
  return me(ctx);
}
