import type { Metadata } from 'next';
import { SecurityForm } from '@/components/security-form';
import { SessionRecovery } from '@/components/session-recovery';
import { pageProfile } from '@/lib/page-profile';
import { siteDictionary } from '@/lib/i18n-server';

export async function generateMetadata(): Promise<Metadata> {
  return { title: (await siteDictionary()).security };
}
export const dynamic = 'force-dynamic';

export default async function SecurityPage() {
  const t = await siteDictionary();
  const result = await pageProfile();
  if (result.kind === 'unauthenticated') return <SessionRecovery />;
  if (result.kind === 'error')
    return (
      <section className="card">
        <h1>{t.accountUnavailable}</h1>
        <p>{t.serviceUnavailable}</p>
      </section>
    );
  return <SecurityForm />;
}
