'use client';
import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { api } from './client-api';

export function AccountMenu({ nickname }: { nickname: string }) {
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
      setError('退出失败，请重试。');
      setBusy(false);
    }
  }

  return (
    <nav className="account-nav" aria-label="账号导航">
      <div className="account-menu-root" ref={root}>
        <button
          type="button"
          className="account-trigger"
          aria-label="账号菜单"
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
              个人资料
            </Link>
            <Link href="/settings/security" onClick={() => setOpen(false)}>
              账号与安全
            </Link>
            <button type="button" onClick={logout} disabled={busy}>
              退出登录
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
