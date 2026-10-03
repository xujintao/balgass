'use client';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from 'react';
import { api } from './client-api';
import { performPasskey } from './passkey-browser';
import { Captcha } from './captcha';
import styles from './auth-form.module.css';
import { errorMessage } from '@/lib/i18n';
import { useDictionary } from './locale-provider';
export function AuthForm({ signup = false }: { signup?: boolean }) {
  const t = useDictionary();
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [emailMode, setEmailMode] = useState(signup);
  const [sent, setSent] = useState(false);
  const [verified, setVerified] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [captchaAction, setCaptchaAction] = useState<'passkey' | 'send' | null>(
    null,
  );
  const captchaPending = useRef(false);
  const captchaDialog = useRef<HTMLDialogElement>(null);
  const [remaining, setRemaining] = useState(0);
  useEffect(() => {
    if (!remaining) return;
    const timer = setTimeout(() => setRemaining((v) => v - 1), 1000);
    return () => clearTimeout(timer);
  }, [remaining]);
  useEffect(() => {
    if (!captchaAction) return;
    if (!captchaDialog.current?.open) captchaDialog.current?.showModal();
  }, [captchaAction]);
  async function run(action: () => Promise<void>) {
    setBusy(true);
    setError('');
    try {
      await action();
    } catch (e) {
      setError(errorMessage(t, e));
    } finally {
      setBusy(false);
    }
  }
  function requestCaptcha(action: 'passkey' | 'send') {
    if (busy || captchaAction) return;
    if (!process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY) {
      setError(t.captchaUnconfigured);
      return;
    }
    setError('');
    captchaPending.current = true;
    setCaptchaAction(action);
  }
  function closeCaptcha() {
    captchaPending.current = false;
    setCaptchaAction(null);
  }
  function completeCaptcha(token: string) {
    if (!captchaPending.current || !captchaAction) return;
    if (!token) {
      closeCaptcha();
      setError(t.captchaIncomplete);
      return;
    }
    captchaPending.current = false;
    const action = captchaAction;
    setCaptchaAction(null);
    void run(async () => {
      if (action === 'send') return send(token);
      await performPasskey('login', token);
      router.replace('/');
      router.refresh();
    });
  }
  async function send(captchaToken: string) {
    await api('auth/otp/send', 'POST', {
      email: email.trim(),
      intent: signup ? 'signup' : 'login',
      captchaToken,
    });
    setSent(true);
    setRemaining(60);
  }
  if (verified)
    return (
      <section className="card">
        <span className="eyebrow">{t.emailVerified}</span>
        <h1>{t.setUpPasskey}</h1>
        <p>{t.passkeyNextTime}</p>
        <button
          disabled={busy}
          onClick={() =>
            run(async () => {
              await performPasskey('register');
              router.replace('/');
              router.refresh();
            })
          }
        >
          {busy ? t.settingUp : t.createPasskey}
        </button>
        <Link className="secondary-link" href="/">
          {t.later}
        </Link>
        {error && (
          <p role="alert" className="error">
            {error}
          </p>
        )}
      </section>
    );
  return (
    <section className="card">
      <h1>{signup ? t.createAccount : t.login}</h1>
      {!signup && !emailMode && (
        <>
          <button disabled={busy} onClick={() => requestCaptcha('passkey')}>
            {busy ? t.loggingIn : t.passkeyLogin}
          </button>
          <button
            className="secondary"
            disabled={busy}
            onClick={() => {
              setEmailMode(true);
              setError('');
            }}
          >
            {t.emailLogin}
          </button>
        </>
      )}
      {emailMode && (
        <form
          onSubmit={(e) => {
            e.preventDefault();
            if (!sent) {
              requestCaptcha('send');
              return;
            }
            void run(async () => {
              await api('auth/otp/verify', 'POST', {
                email: email.trim(),
                code,
              });
              if (signup) setVerified(true);
              else {
                router.replace('/');
                router.refresh();
              }
            });
          }}
        >
          <label htmlFor="email">{t.email}</label>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            required
            maxLength={254}
            value={email}
            disabled={busy || sent}
            onChange={(e) => setEmail(e.target.value)}
          />
          {sent && (
            <>
              <label htmlFor="code">{t.emailCode}</label>
              <input
                id="code"
                name="code"
                inputMode="numeric"
                autoComplete="one-time-code"
                pattern="[0-9]{6}"
                minLength={6}
                maxLength={6}
                required
                value={code}
                onChange={(e) => setCode(e.target.value)}
              />
              <p className="hint">{t.codeSent(email)}</p>
            </>
          )}
          <button disabled={busy} type="submit">
            {busy ? t.pleaseWait : sent ? t.verifyContinue : t.sendCode}
          </button>
          {sent && (
            <div className="row">
              <button
                type="button"
                className="text-button"
                disabled={busy || remaining > 0}
                onClick={() => requestCaptcha('send')}
              >
                {remaining ? t.resendIn(remaining) : t.resend}
              </button>
              <button
                type="button"
                className="text-button"
                disabled={busy}
                onClick={() => {
                  setSent(false);
                  setCode('');
                }}
              >
                {t.changeEmail}
              </button>
            </div>
          )}
        </form>
      )}
      {captchaAction && (
        <dialog
          ref={captchaDialog}
          className={styles.dialog}
          aria-labelledby="captcha-dialog-title"
          onCancel={closeCaptcha}
        >
          <h2 id="captcha-dialog-title">{t.captchaTitle}</h2>
          <Captcha onToken={completeCaptcha} />
          <button type="button" className="secondary" onClick={closeCaptcha}>
            {t.cancel}
          </button>
        </dialog>
      )}
      {error && (
        <p role="alert" className="error">
          {error}
        </p>
      )}
      <div className="footer-actions">
        {!signup && emailMode && (
          <button
            className="text-button"
            onClick={() => {
              setEmailMode(false);
              setError('');
            }}
          >
            {t.backToPasskey}
          </button>
        )}
        <p className="footer-note">
          {signup ? t.haveAccount : t.needAccount}{' '}
          <Link href={signup ? '/login' : '/signup'}>
            {signup ? t.login : t.signup}
          </Link>
        </p>
      </div>
    </section>
  );
}
