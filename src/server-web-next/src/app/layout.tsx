import type { Metadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';
import { AccountNav } from '@/components/account-nav';
import { LanguageSwitcher } from '@/components/language-switcher';
import { LocaleProvider } from '@/components/locale-provider';
import { siteDictionary, siteLocale } from '@/lib/i18n-server';
import './globals.css';
export async function generateMetadata(): Promise<Metadata> {
  const t = await siteDictionary();
  return {
    title: { default: 'r2f2', template: '%s · r2f2' },
    description: t.description,
  };
}
export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = await siteLocale();
  const t = await siteDictionary();
  return (
    <html lang={locale}>
      <body>
        <LocaleProvider key={locale} locale={locale}>
          <header>
            <div className="header-primary">
              <Link className="brand" href="/">
                <Image src="/icon.svg" alt="r2f2" width={36} height={36} />
              </Link>
              <nav className="site-nav" aria-label={t.primaryNavigation}>
                <Link href="/game">{t.gameMap}</Link>
              </nav>
            </div>
            <div className="header-actions">
              <LanguageSwitcher />
              <AccountNav />
            </div>
          </header>
          <main>{children}</main>
        </LocaleProvider>
      </body>
    </html>
  );
}
