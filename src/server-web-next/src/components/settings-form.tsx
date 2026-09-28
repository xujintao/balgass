'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import type { me } from '@/lib/profile';
import { api, ClientError } from './client-api';
import { performPasskey } from './passkey-browser';
type Profile = Awaited<ReturnType<typeof me>>;
type Passkey = {
  id: string;
  friendly_name?: string;
  created_at: string;
  last_used_at?: string;
};
export function SettingsForm({ initial }: { initial: Profile }) {
  const router = useRouter();
  const [profile, setProfile] = useState(initial);
  const [nickname, setNickname] = useState(initial.nickname);
  const [passkeys, setPasskeys] = useState<Passkey[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [busy, setBusy] = useState(false);
  const next = profile.nextNicknameChangeAt;
  const cooling = next !== null && new Date(next).getTime() > Date.now();
  async function load() {
    setPasskeys(await api<Passkey[]>('auth/passkeys'));
    setLoaded(true);
  }
  useEffect(() => {
    api<Passkey[]>('auth/passkeys')
      .then((v) => {
        setPasskeys(v);
        setLoaded(true);
      })
      .catch((e) => setError(e.message));
  }, []);
  async function run(action: () => Promise<void>) {
    setBusy(true);
    setError('');
    setNotice('');
    try {
      await action();
    } catch (e) {
      setError(e instanceof Error ? e.message : '操作失败。');
      if (e instanceof ClientError && e.nextNicknameChangeAt)
        setProfile((v) => ({
          ...v,
          nextNicknameChangeAt: e.nextNicknameChangeAt!,
        }));
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="settings-grid">
      <section className="card">
        <span className="eyebrow">个人资料</span>
        <h1>账号设置</h1>
        <label>登录邮箱</label>
        <p className="email-value">{profile.email}</p>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            void run(async () => {
              const updated = await api<Profile>('me/nickname', 'PATCH', {
                nickname,
              });
              setProfile(updated);
              setNickname(updated.nickname);
              setNotice('昵称已保存。');
            });
          }}
        >
          <label htmlFor="nickname">昵称</label>
          <input
            id="nickname"
            value={nickname}
            required
            onChange={(e) => setNickname(e.target.value)}
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
      </section>
      <section className="card">
        <span className="eyebrow">登录与安全</span>
        <h2>你的 Passkey</h2>
        <p>为常用设备添加 Passkey。也可以随时使用邮件验证码登录。</p>
        {loaded && passkeys.length === 0 && (
          <p className="hint">尚未添加 Passkey。</p>
        )}
        <ul className="passkeys">
          {passkeys.map((key) => (
            <li key={key.id}>
              <div>
                <strong>{key.friendly_name || 'Passkey'}</strong>
                <small>
                  添加于 {new Date(key.created_at).toLocaleDateString('zh-CN')}
                </small>
              </div>
              <button
                className="text-button danger"
                disabled={busy}
                aria-label={`删除 ${key.friendly_name || 'Passkey'}`}
                onClick={() => {
                  if (
                    window.confirm('删除此 Passkey？你仍可使用邮件验证码登录。')
                  )
                    void run(async () => {
                      await api(`auth/passkeys/${key.id}`, 'DELETE');
                      await load();
                      setNotice('Passkey 已删除。');
                    });
                }}
              >
                删除
              </button>
            </li>
          ))}
        </ul>
        <button
          disabled={busy}
          onClick={() =>
            run(async () => {
              await performPasskey('register');
              await load();
              setNotice('Passkey 已添加。');
            })
          }
        >
          {busy ? '请稍候…' : '添加 Passkey'}
        </button>
        <button
          className="secondary"
          disabled={busy}
          onClick={() =>
            run(async () => {
              await api('auth/logout', 'POST');
              router.replace('/login');
              router.refresh();
            })
          }
        >
          退出登录
        </button>
      </section>
      {error && (
        <p role="alert" className="error full-width">
          {error}
        </p>
      )}
      {notice && (
        <p role="status" className="success full-width">
          {notice}
        </p>
      )}
    </div>
  );
}
