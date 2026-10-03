import type { Metadata } from 'next';
import Link from 'next/link';
import { AccountNav } from '@/components/account-nav';
import './globals.css';
export const metadata: Metadata = {
  title: { default: 'r2f2', template: '%s · r2f2' },
  description: '邮箱注册，使用 Passkey 安全登录。',
};
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-CN">
      <body>
        <header>
          <Link className="brand" href="/">
            r2f2
          </Link>
          <AccountNav />
        </header>
        <main>{children}</main>
      </body>
    </html>
  );
}
