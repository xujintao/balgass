'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api } from './client-api';
import { useDictionary } from './locale-provider';
export function SessionRecovery() {
  const t = useDictionary();
  const [failed, setFailed] = useState(false);
  useEffect(() => {
    api('auth/session/refresh', 'POST', undefined, false)
      .then(() => window.location.reload())
      .catch(() => setFailed(true));
  }, []);
  return (
    <section className="card">
      <h1>{failed ? t.loginFirst : t.restoringSession}</h1>
      {failed && (
        <>
          <p>{t.sessionEnded}</p>
          <Link className="button-link" href="/login">
            {t.goToLogin}
          </Link>
        </>
      )}
    </section>
  );
}
