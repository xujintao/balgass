'use client';
import { useState } from 'react';
import type { me } from '@/lib/profile';
import { api, ClientError } from './client-api';
import { errorMessage, formatDateTime } from '@/lib/i18n';
import { useDictionary, useLocale } from './locale-provider';

type Profile = Awaited<ReturnType<typeof me>>;

export function ProfileForm({ initial }: { initial: Profile }) {
  const t = useDictionary();
  const locale = useLocale();
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
      setNotice(t.nicknameSaved);
    } catch (e) {
      setError(errorMessage(t, e));
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
      <h1>{t.profile}</h1>
      <label>{t.loginEmail}</label>
      <p className="email-value">{profile.email}</p>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          void saveNickname();
        }}
      >
        <label htmlFor="nickname">{t.nickname}</label>
        <input
          id="nickname"
          value={nickname}
          required
          onChange={(event) => setNickname(event.target.value)}
          disabled={busy || cooling}
        />
        <p className="hint">{t.nicknameHint}</p>
        {cooling && (
          <p className="hint">
            {t.nextNicknameChange}{' '}
            <time dateTime={next!}>{formatDateTime(locale, next!)}</time>
          </p>
        )}
        <button disabled={busy || cooling}>{t.saveNickname}</button>
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
