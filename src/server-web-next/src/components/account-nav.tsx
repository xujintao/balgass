import Link from 'next/link';
import { cookies } from 'next/headers';
import { pageProfile } from '@/lib/page-profile';
import { ACCESS_COOKIE, REFRESH_COOKIE } from '@/lib/session';
import { AccountMenu } from './account-menu';
import { NavRecovery } from './nav-recovery';
import { siteDictionary } from '@/lib/i18n-server';

export async function GuestNav() {
  const t = await siteDictionary();
  return (
    <nav className="account-nav" aria-label={t.accountNavigation}>
      <Link href="/login">{t.login}</Link>
    </nav>
  );
}

export async function AccountNav() {
  const jar = await cookies();
  const access = jar.has(ACCESS_COOKIE);
  const refresh = jar.has(REFRESH_COOKIE);
  if (!access && !refresh) return <GuestNav />;
  if (!access) return <NavRecovery />;

  const result = await pageProfile();
  if (result.kind === 'ok')
    return <AccountMenu nickname={result.profile.nickname} />;
  if (result.kind === 'unauthenticated' && refresh) return <NavRecovery />;
  return <GuestNav />;
}
