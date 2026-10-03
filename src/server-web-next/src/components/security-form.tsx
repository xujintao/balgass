'use client';
import { useEffect, useState } from 'react';
import { api } from './client-api';
import { performPasskey } from './passkey-browser';
import { errorMessage, formatDate } from '@/lib/i18n';
import { useDictionary, useLocale } from './locale-provider';

type Passkey = {
  id: string;
  friendly_name?: string;
  created_at: string;
  last_used_at?: string;
};

export function SecurityForm() {
  const t = useDictionary();
  const locale = useLocale();
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
      .catch((e) => setError(errorMessage(t, e)));
  }, [t]);

  async function run(action: () => Promise<void>) {
    setBusy(true);
    setError('');
    setNotice('');
    try {
      await action();
    } catch (e) {
      setError(errorMessage(t, e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h1>{t.security}</h1>
      <h2>{t.yourPasskeys}</h2>
      <p>{t.passkeyDescription}</p>
      {loaded && passkeys.length === 0 && (
        <p className="hint">{t.noPasskeys}</p>
      )}
      <ul className="passkeys">
        {passkeys.map((key) => (
          <li key={key.id}>
            <div>
              <strong>{key.friendly_name || 'Passkey'}</strong>
              <small>
                {t.addedOn} {formatDate(locale, key.created_at)}
              </small>
            </div>
            <button
              className="text-button danger"
              disabled={busy}
              aria-label={t.deletePasskey(key.friendly_name || 'Passkey')}
              onClick={() => {
                if (window.confirm(t.confirmDelete))
                  void run(async () => {
                    await api(`auth/passkeys/${key.id}`, 'DELETE');
                    await load();
                    setNotice(t.passkeyDeleted);
                  });
              }}
            >
              {t.delete}
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
            setNotice(t.passkeyAdded);
          })
        }
      >
        {busy ? t.pleaseWait : t.addPasskey}
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
