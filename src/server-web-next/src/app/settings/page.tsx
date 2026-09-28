import { pageProfile } from '@/lib/page-profile';
import { SettingsForm } from '@/components/settings-form';
import { SessionRecovery } from '@/components/session-recovery';
export const dynamic = 'force-dynamic';
export default async function Settings() {
  const result = await pageProfile();
  if (result.kind === 'unauthenticated') return <SessionRecovery />;
  if (result.kind === 'error')
    return (
      <section className="card">
        <h1>暂时无法读取账号</h1>
        <p>{result.message}</p>
      </section>
    );
  return <SettingsForm initial={result.profile} />;
}
