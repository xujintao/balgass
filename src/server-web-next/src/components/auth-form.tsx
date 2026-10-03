'use client';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { api } from './client-api';
import { performPasskey } from './passkey-browser';
import { Captcha } from './captcha';
export function AuthForm({ signup = false }: { signup?: boolean }) {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [emailMode, setEmailMode] = useState(signup);
  const [sent, setSent] = useState(false);
  const [verified, setVerified] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [captchaToken, setCaptchaToken] = useState('');
  const [captchaKey, setCaptchaKey] = useState(0);
  const [remaining, setRemaining] = useState(0);
  useEffect(() => {
    if (!remaining) return;
    const timer = setTimeout(() => setRemaining((v) => v - 1), 1000);
    return () => clearTimeout(timer);
  }, [remaining]);
  async function run(action: () => Promise<void>, consumeCaptcha = false) {
    setBusy(true);
    setError('');
    try {
      await action();
    } catch (e) {
      setError(e instanceof Error ? e.message : '操作失败，请重试。');
    } finally {
      setBusy(false);
      if (consumeCaptcha) {
        setCaptchaToken('');
        setCaptchaKey((v) => v + 1);
      }
    }
  }
  async function send() {
    await api('auth/otp/send', 'POST', {
      email: email.trim(),
      intent: signup ? 'signup' : 'login',
      captchaToken: captchaToken || undefined,
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
      <span className="eyebrow">r2f2 账号</span>
      <h1>{signup ? '创建你的账号' : '欢迎回来'}</h1>
      <p>
        {signup
          ? '使用邮箱注册，无需设置密码。'
          : '使用 Passkey 安全快捷地登录。'}
      </p>
      {!signup && !emailMode && (
        <>
          <button
            disabled={busy}
            onClick={() =>
              run(async () => {
                await performPasskey('login', captchaToken);
                router.replace('/');
                router.refresh();
              }, true)
            }
          >
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
            void run(async () => {
              if (!sent) return send();
              await api('auth/otp/verify', 'POST', {
                email: email.trim(),
                code,
              });
              if (signup) setVerified(true);
              else {
                router.replace('/');
                router.refresh();
              }
            }, !sent);
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
                onClick={() => run(send, true)}
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
      {(!sent || remaining === 0) && (
        <Captcha key={captchaKey} onToken={setCaptchaToken} />
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
