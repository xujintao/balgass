'use client';
import { useEffect, useState } from 'react';
import { api } from './client-api';
import { performPasskey } from './passkey-browser';

type Passkey = {
  id: string;
  friendly_name?: string;
  created_at: string;
  last_used_at?: string;
};

export function SecurityForm() {
  const [passkeys, setPasskeys] = useState<Passkey[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [busy, setBusy] = useState(false);

  async function load() {
    setPasskeys(await api<Passkey[]>('auth/passkeys'));
    setLoaded(true);
  }

  useEffect(() => {
    api<Passkey[]>('auth/passkeys')
      .then((value) => {
        setPasskeys(value);
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
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h1>账号与安全</h1>
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
