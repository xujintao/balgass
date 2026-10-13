import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { ApiError } from '@/lib/errors';
import { listAdminUsers } from '@/lib/admin-users';
import { context } from '@/lib/session';
import { siteLocale } from '@/lib/i18n-server';
import { AdminUsers } from '@/components/admin-users';
import { adminText } from '@/lib/admin-users-i18n';
import styles from './users.module.css';

export const dynamic = 'force-dynamic';
export async function generateMetadata(): Promise<Metadata> {
  return { title: adminText[await siteLocale()].title };
}

export default async function AdminUsersPage({ searchParams }: {
  searchParams: Promise<{ page?: string }>;
}) {
  const t = adminText[await siteLocale()];
  let result: Awaited<ReturnType<typeof listAdminUsers>>;
  let selfId: string;
  try {
    const ctx = await context();
    const page = Number((await searchParams).page ?? '1');
    result = await listAdminUsers(ctx, page);
    selfId = ctx.user.id;
  } catch (error) {
    if (error instanceof ApiError && [401, 403].includes(error.status)) notFound();
    return <section className={styles.page}><h1>{t.title}</h1>
      <p role="alert">{error instanceof ApiError && error.status === 400 ? t.invalidPage : t.unavailable}</p>
    </section>;
  }
  return <AdminUsers key={result.page} initial={result} selfId={selfId} />;
}
