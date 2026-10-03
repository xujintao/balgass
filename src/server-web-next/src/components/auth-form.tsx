'use client';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from 'react';
import { api } from './client-api';
import { performPasskey } from './passkey-browser';
import { Captcha } from './captcha';
import styles from './auth-form.module.css';
export function AuthForm({ signup = false }: { signup?: boolean }) {
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
      setError(e instanceof Error ? e.message : '操作失败，请重试。');
    } finally {
      setBusy(false);
    }
  }
  function requestCaptcha(action: 'passkey' | 'send') {
    if (busy || captchaAction) return;
    if (!process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY) {
      setError('人机验证未配置，请稍后重试。');
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
      setError('人机验证未完成，请重试。');
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
        <span className="eyebrow">邮箱已验证</span>
        <h1>设置你的 Passkey</h1>
        <p>下次使用指纹、面容或设备 PIN 即可登录。</p>
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
          {busy ? '正在设置…' : '创建 Passkey'}
        </button>
        <Link className="secondary-link" href="/">
          稍后设置
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
      <h1>{signup ? '创建你的账号' : '登录'}</h1>
      {!signup && !emailMode && (
        <>
          <button disabled={busy} onClick={() => requestCaptcha('passkey')}>
            {busy ? '正在登录…' : '使用 Passkey 登录'}
          </button>
          <button
            className="secondary"
            disabled={busy}
            onClick={() => {
              setEmailMode(true);
              setError('');
            }}
          >
            使用邮件验证码
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
          <label htmlFor="email">邮箱</label>
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
              <label htmlFor="code">邮件验证码</label>
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
              <p className="hint">
                若邮箱可以接收此请求，验证码将发送到 {email}，10 分钟内有效。
              </p>
            </>
          )}
          <button disabled={busy} type="submit">
            {busy ? '请稍候…' : sent ? '验证并继续' : '发送验证码'}
          </button>
          {sent && (
            <div className="row">
              <button
                type="button"
                className="text-button"
                disabled={busy || remaining > 0}
                onClick={() => requestCaptcha('send')}
              >
                {remaining ? `${remaining} 秒后可重发` : '重新发送'}
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
                更换邮箱
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
          <h2 id="captcha-dialog-title">请完成人机验证</h2>
          <Captcha onToken={completeCaptcha} />
          <button type="button" className="secondary" onClick={closeCaptcha}>
            取消
          </button>
        </dialog>
      )}
      {error && (
        <p role="alert" className="error">
          {error}
        </p>
      )}
      <p className="footer-note">
        {signup ? '已有账号？' : '还没有账号？'}{' '}
        <Link href={signup ? '/login' : '/signup'}>
          {signup ? '登录' : '注册'}
        </Link>
      </p>
      {!signup && emailMode && (
        <button
          className="text-button"
          onClick={() => {
            setEmailMode(false);
            setError('');
          }}
        >
          返回 Passkey 登录
        </button>
      )}
    </section>
  );
}
