import type { Metadata } from 'next';
import Link from 'next/link';
import './globals.css';
export const metadata: Metadata = {
  title: 'Balgass · 账号',
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
            B<span>·</span> BALGASS
          </Link>
          <nav>
            <Link href="/settings">账号设置</Link>
          </nav>
        </header>
        <main>{children}</main>
        <footer>你的账号，你的世界。</footer>
      </body>
    </html>
  );
}
