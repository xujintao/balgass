import { randomBytes } from 'node:crypto';
import { ApiError } from './errors';
import { nickname } from './validation';

export const NICKNAME_COOLDOWN_MS = 30 * 24 * 60 * 60 * 1000;
export type ProfileRow = {
  user_id: string;
  nickname: string;
  nickname_key: string;
  created_at: string;
  nickname_changed_at: string | null;
};
export type ProfileUser = {
  id: string;
  email?: string;
  email_confirmed_at?: string;
};
export interface ProfileRepository {
  find(userId: string): Promise<ProfileRow | null>;
  insert(
    userId: string,
    nickname: string,
    key: string,
  ): Promise<ProfileRow | null>;
  // Return null when another request has changed the observed row.
  update(
    observed: ProfileRow,
    nickname: string,
    key: string,
    changedAt: string,
    cutoff: string,
  ): Promise<ProfileRow | null>;
}
export const nicknameKey = (value: string) => value.toLowerCase();
export const randomNickname = () => `玩家_${randomBytes(8).toString('hex')}`;

export function profileService(
  repository: ProfileRepository,
  generate = randomNickname,
  clock = () => new Date(),
) {
  function authorize(user: ProfileUser) {
    if (!user.email_confirmed_at)
      throw new ApiError(403, 'EMAIL_UNVERIFIED', '请先验证邮箱。');
  }
  function present(user: ProfileUser, row: ProfileRow) {
    return {
      id: user.id,
      email: user.email,
      nickname: row.nickname,
      createdAt: row.created_at,
      nicknameChangedAt: row.nickname_changed_at,
      nextNicknameChangeAt: row.nickname_changed_at
        ? new Date(
            new Date(row.nickname_changed_at).getTime() + NICKNAME_COOLDOWN_MS,
          ).toISOString()
        : null,
    };
  }
  async function ensure(user: ProfileUser) {
    authorize(user);
    const existing = await repository.find(user.id);
    if (existing) return existing;
    for (let attempt = 0; attempt < 10; attempt++) {
      const value = nickname.parse(generate());
      const inserted = await repository.insert(
        user.id,
        value,
        nicknameKey(value),
      );
      if (inserted) return inserted;
      // A simultaneous login might have created the same user's profile.
      const winner = await repository.find(user.id);
      if (winner) return winner;
    }
    throw new ApiError(
      503,
      'PROFILE_UNAVAILABLE',
      '暂时无法分配昵称，请重试。',
    );
  }
  async function me(user: ProfileUser) {
    return present(user, await ensure(user));
  }
  async function change(user: ProfileUser, requested: string) {
    authorize(user);
    const value = nickname.parse(requested);
    const observed = await ensure(user);
    if (observed.nickname === value) return present(user, observed);
    const now = clock();
    const next = observed.nickname_changed_at
      ? new Date(observed.nickname_changed_at).getTime() + NICKNAME_COOLDOWN_MS
      : null;
    if (next !== null && now.getTime() < next)
      throw new ApiError(
        409,
        'NICKNAME_COOLDOWN',
        '每 30 天只能修改一次昵称。',
        { nextNicknameChangeAt: new Date(next).toISOString() },
      );
    const updated = await repository.update(
      observed,
      value,
      nicknameKey(value),
      now.toISOString(),
      new Date(now.getTime() - NICKNAME_COOLDOWN_MS).toISOString(),
    );
    if (updated) return present(user, updated);
    const current = await repository.find(user.id);
    if (current?.nickname === value) return present(user, current);
    if (current?.nickname_changed_at)
      throw new ApiError(
        409,
        'NICKNAME_COOLDOWN',
        '每 30 天只能修改一次昵称。',
        {
          nextNicknameChangeAt: new Date(
            new Date(current.nickname_changed_at).getTime() +
              NICKNAME_COOLDOWN_MS,
          ).toISOString(),
        },
      );
    throw new ApiError(
      409,
      'PROFILE_CHANGED',
      '资料已发生变化，请刷新后重试。',
    );
  }
  return { ensure, me, change };
}
