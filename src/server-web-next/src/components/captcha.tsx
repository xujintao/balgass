'use client';
import Script from 'next/script';
import { useEffect, useRef } from 'react';
declare global {
  interface Window {
    turnstile?: {
      render: (
        element: HTMLElement,
        options: Record<string, unknown>,
      ) => string;
      remove: (id: string) => void;
    };
  }
}
export function Captcha({ onToken }: { onToken: (token: string) => void }) {
  const container = useRef<HTMLDivElement>(null);
  const id = useRef<string | null>(null);
  const sitekey = process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY;
  useEffect(
    () => () => {
      if (id.current && window.turnstile) {
        window.turnstile.remove(id.current);
        id.current = null;
      }
    },
    [],
  );
  function mount() {
    if (!container.current || !window.turnstile || id.current) return;
    id.current = window.turnstile.render(container.current, {
      sitekey,
      callback: onToken,
      'expired-callback': () => onToken(''),
      'error-callback': () => onToken(''),
    });
  }
  // Each instance has a fresh container; the provider handles expiration.
  if (!sitekey) return null;
  return (
    <>
      <Script
        src="https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit"
        strategy="afterInteractive"
        onReady={mount}
      />
      <div ref={container} className="captcha" />
    </>
  );
}
