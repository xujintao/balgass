'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api } from './client-api';
export function SessionRecovery() {
  const [failed, setFailed] = useState(false);
  useEffect(() => {
    api('auth/session/refresh', 'POST', undefined, false)
      .then(() => window.location.reload())
      .catch(() => setFailed(true));
  }, []);
  return (
    <section className="card">
      <h1>{failed ? '请先登录' : '正在恢复会话…'}</h1>
      {failed && (
        <>
          <p>会话已结束，请使用 Passkey 或邮件验证码重新登录。</p>
          <Link className="button-link" href="/login">
            前往登录
          </Link>
        </>
      )}
    </section>
  );
}
