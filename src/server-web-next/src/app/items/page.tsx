import type { Metadata } from 'next';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import { ItemBrowser } from '@/components/item-browser';
import { itemKindGroups, listItems, parseItemKind } from '@/lib/item-catalog';
import { itemDictionaries } from '@/lib/item-i18n';
import { siteLocale } from '@/lib/i18n-server';
import styles from './items.module.css';

export const dynamic = 'force-dynamic';
export async function generateMetadata(): Promise<Metadata> {
  return { title: itemDictionaries[await siteLocale()].items };
}
export default async function ItemsPage({ searchParams }: {
  searchParams: Promise<{ kind?: string | string[] }>;
}) {
  const locale = await siteLocale();
  const t = itemDictionaries[locale];
  const raw = (await searchParams).kind;
  if (Array.isArray(raw)) notFound();
  let kind;
  try { kind = parseItemKind(raw ?? 'sword'); } catch { notFound(); }
  let result: Awaited<ReturnType<typeof listItems>> | null = null;
  try { result = await listItems(kind, locale); } catch { /* show a stable error state */ }
  return <section className={styles.page}>
    <div className={styles.heading}><h1>{t.items}</h1><Link href="/orders">{t.orders}</Link></div>
    <nav className={styles.categories} aria-label={t.items}>
      {itemKindGroups.map((group) => <div className={styles.categoryGroup} key={group[0]}>
        {group.map((entry) => <Link key={entry} href={`/items?kind=${entry}`} aria-current={kind === entry ? 'page' : undefined}>
          {t.categories[entry]}
        </Link>)}
      </div>)}
    </nav>
    {result ? <ItemBrowser result={result} /> : <p className="error" role="alert">{t.unavailable}</p>}
  </section>;
}
