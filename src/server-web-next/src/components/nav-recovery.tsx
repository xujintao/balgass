'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from './client-api';
import { useDictionary } from './locale-provider';

export function NavRecovery() {
  const t = useDictionary();
  const router = useRouter();
  const [failed, setFailed] = useState(false);
  useEffect(() => {
    api('auth/session/refresh', 'POST', undefined, false)
      .then(() => router.refresh())
      .catch(() => setFailed(true));
  }, [router]);
  if (!failed)
    return (
      <nav
        className="account-nav"
        aria-label={t.accountNavigation}
        aria-busy="true"
      />
    );
  return (
    <nav className="account-nav" aria-label={t.accountNavigation}>
      <Link href="/login">{t.login}</Link>
    </nav>
  );
}
