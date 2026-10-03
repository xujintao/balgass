'use client';
import { useState } from 'react';
import type { me } from '@/lib/profile';
import { api, ClientError } from './client-api';

type Profile = Awaited<ReturnType<typeof me>>;

export function ProfileForm({ initial }: { initial: Profile }) {
  const [profile, setProfile] = useState(initial);
  const [nickname, setNickname] = useState(initial.nickname);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [busy, setBusy] = useState(false);
  const next = profile.nextNicknameChangeAt;
  const cooling = next !== null && new Date(next).getTime() > Date.now();

  async function saveNickname() {
    setBusy(true);
    setError('');
    setNotice('');
    try {
      const updated = await api<Profile>('me/nickname', 'PATCH', { nickname });
      setProfile(updated);
      setNickname(updated.nickname);
      setNotice('昵称已保存。');
    } catch (e) {
      setError(e instanceof Error ? e.message : '操作失败。');
      if (e instanceof ClientError && e.nextNicknameChangeAt)
        setProfile((value) => ({
          ...value,
          nextNicknameChangeAt: e.nextNicknameChangeAt!,
        }));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h1>个人资料</h1>
      <label>登录邮箱</label>
      <p className="email-value">{profile.email}</p>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          void saveNickname();
        }}
      >
        <label htmlFor="nickname">昵称</label>
        <input
          id="nickname"
          value={nickname}
          required
          onChange={(event) => setNickname(event.target.value)}
          disabled={busy || cooling}
        />
        <p className="hint">
          2–24 个汉字、英文字母、数字或下划线。昵称唯一，每 30 天可修改一次。
        </p>
        {cooling && (
          <p className="hint">
            下次可修改时间：
            <time dateTime={next!}>
              {new Date(next!).toLocaleString('zh-CN')}
            </time>
          </p>
        )}
        <button disabled={busy || cooling}>保存昵称</button>
      </form>
      {error && (
        <p role="alert" className="error">
          {error}
        </p>
      )}
      {notice && (
        <p role="status" className="success">
          {notice}
        </p>
      )}
    </section>
  );
}
