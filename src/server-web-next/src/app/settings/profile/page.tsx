import type { Metadata } from 'next';
import { ProfileForm } from '@/components/profile-form';
import { SessionRecovery } from '@/components/session-recovery';
import { pageProfile } from '@/lib/page-profile';

export const metadata: Metadata = { title: '个人资料' };
export const dynamic = 'force-dynamic';

export default async function ProfilePage() {
  const result = await pageProfile();
  if (result.kind === 'unauthenticated') return <SessionRecovery />;
  if (result.kind === 'error')
    return (
      <section className="card">
        <h1>暂时无法读取账号</h1>
        <p>{result.message}</p>
      </section>
    );
  return <ProfileForm initial={result.profile} />;
}
