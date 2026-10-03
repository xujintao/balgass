'use client';
import { useRouter } from 'next/navigation';
import { LOCALE_COOKIE, type Locale } from '@/lib/i18n';
import { useDictionary, useLocale } from './locale-provider';

export function LanguageSwitcher() {
  const locale = useLocale();
  const t = useDictionary();
  const router = useRouter();
  function change(value: string) {
    if (value !== 'en' && value !== 'zh-CN' && value !== 'es') return;
    document.cookie = `${LOCALE_COOKIE}=${value}; Path=/; Max-Age=31536000; SameSite=Lax${location.protocol === 'https:' ? '; Secure' : ''}`;
    router.refresh();
  }
  return (
    <label className="language-switcher">
      <span className="sr-only">{t.language}</span>
      <select
        aria-label={t.language}
        value={locale}
        onChange={(event) => change(event.target.value as Locale)}
      >
        <option value="en">English</option>
        <option value="zh-CN">简体中文</option>
        <option value="es">Español</option>
      </select>
    </label>
  );
}
