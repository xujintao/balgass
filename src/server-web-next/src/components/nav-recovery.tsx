'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from './client-api';

export function NavRecovery() {
  const router = useRouter();
  const [failed, setFailed] = useState(false);
  useEffect(() => {
    api('auth/session/refresh', 'POST', undefined, false)
      .then(() => router.refresh())
      .catch(() => setFailed(true));
  }, [router]);
  if (!failed)
    return (
      <nav className="account-nav" aria-label="账号导航" aria-busy="true" />
    );
  return (
    <nav className="account-nav" aria-label="账号导航">
      <Link href="/login">登录</Link>
    </nav>
  );
}
