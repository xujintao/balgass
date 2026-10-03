import Link from 'next/link';
import { cookies } from 'next/headers';
import { ACCESS_COOKIE } from '@/lib/session';
import { pageProfile } from '@/lib/page-profile';
import { SessionRecovery } from '@/components/session-recovery';
export const dynamic = 'force-dynamic';
export default async function Home() {
  if (!(await cookies()).has(ACCESS_COOKIE))
    return (
      <section className="card hero">
        <span className="eyebrow">欢迎来到 r2f2</span>
        <h1>从你的账号开始。</h1>
        <p>使用邮箱注册，以 Passkey 登录。无需记住密码。</p>
        <Link href="/signup" className="button-link">
          创建账号
        </Link>
        <Link href="/login" className="secondary-link">
          已有账号？登录
        </Link>
      </section>
    );
  const result = await pageProfile();
  if (result.kind === 'unauthenticated') return <SessionRecovery />;
  if (result.kind === 'error')
    return (
      <section className="card">
        <h1>服务暂不可用</h1>
        <p>{result.message}</p>
      </section>
    );
  return (
    <section className="card hero">
      <span className="eyebrow">已登录</span>
      <h1>你好，{result.profile.nickname}</h1>
      <p>{result.profile.email}</p>
      <Link className="button-link" href="/settings">
        管理账号与 Passkey
      </Link>
    </section>
  );
}
