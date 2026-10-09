'use client';
import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { api } from './client-api';
import { useDictionary } from './locale-provider';
import { useLocale } from './locale-provider';
import { itemDictionaries } from '@/lib/item-i18n';

export function AccountMenu({ nickname }: { nickname: string }) {
  const t = useDictionary();
  const shop = itemDictionaries[useLocale()];
  const root = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!open) return;
    const closeOnOutside = (event: PointerEvent) => {
      if (!root.current?.contains(event.target as Node)) setOpen(false);
    };
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setOpen(false);
        root.current?.querySelector('button')?.focus();
      }
    };
    document.addEventListener('pointerdown', closeOnOutside);
    document.addEventListener('keydown', closeOnEscape);
    return () => {
      document.removeEventListener('pointerdown', closeOnOutside);
      document.removeEventListener('keydown', closeOnEscape);
    };
  }, [open]);

  async function logout() {
    setBusy(true);
    setError('');
    try {
      await api('auth/logout', 'POST');
      window.location.assign('/');
    } catch {
      setError(t.logoutFailed);
      setBusy(false);
    }
  }

  return (
    <nav className="account-nav" aria-label={t.accountNavigation}>
      <div className="account-menu-root" ref={root}>
        <button
          type="button"
          className="account-trigger"
          aria-label={t.accountMenu}
          aria-expanded={open}
          aria-controls="account-menu"
          onClick={() => setOpen((value) => !value)}
        >
          <svg
            viewBox="0 0 24 24"
            width="22"
            height="22"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.8"
            aria-hidden="true"
          >
            <circle cx="12" cy="8" r="3.5" />
            <path d="M4.5 20c.5-4 3-6 7.5-6s7 2 7.5 6" />
          </svg>
        </button>
        {open && (
          <div className="account-popover" id="account-menu">
            <div className="account-nickname">{nickname}</div>
            <Link href="/settings/profile" onClick={() => setOpen(false)}>
              {t.profile}
            </Link>
            <Link href="/settings/security" onClick={() => setOpen(false)}>
              {t.security}
            </Link>
            <Link href="/game/accounts" onClick={() => setOpen(false)}>
              {t.gameAccounts}
            </Link>
            <Link href="/orders" onClick={() => setOpen(false)}>
              {shop.orders}
            </Link>
            <button type="button" onClick={logout} disabled={busy}>
              {t.logout}
            </button>
            {error && (
              <p role="alert" className="error">
                {error}
              </p>
            )}
          </div>
        )}
      </div>
    </nav>
  );
}
