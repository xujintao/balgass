import type { Metadata } from 'next';
import Link from 'next/link';
import { ApiError } from '@/lib/errors';
import { formatDateTime } from '@/lib/i18n';
import { itemDictionaries } from '@/lib/item-i18n';
import { listOrders } from '@/lib/item-orders';
import { siteLocale } from '@/lib/i18n-server';
import { context } from '@/lib/session';
import { SessionRecovery } from '@/components/session-recovery';
import styles from '../items/items.module.css';

export const dynamic = 'force-dynamic';
export async function generateMetadata(): Promise<Metadata> {
  return { title: itemDictionaries[await siteLocale()].orders };
}
export default async function OrdersPage() {
  const locale = await siteLocale();
  const t = itemDictionaries[locale];
  let orders: Awaited<ReturnType<typeof listOrders>>;
  try { orders = await listOrders(await context(), locale); }
  catch (error) {
    if (error instanceof ApiError && error.status === 401) return <SessionRecovery />;
    return <section className={styles.page}><h1>{t.orders}</h1><p className="error" role="alert">{t.errors.ORDERS_UNAVAILABLE}</p></section>;
  }
  return <section className={styles.page}>
    <div className={styles.heading}><h1>{t.orders}</h1><Link href="/items">{t.items}</Link></div>
    {orders.length === 0 ? <p>{t.noOrders}</p> : <div className={styles.list}>
      {orders.map((order) => <article className={styles.set} key={order.id}>
        <h2>{order.item.name}</h2>
        <p>{t.order} {order.id} · {formatDateTime(locale, order.createdAt)}</p>
        <p>{t.status}: {order.status === 'PENDING' ? t.pending : order.status}</p>
        <p>{t.level}: +{order.item.level} · {t.additional}: +{order.item.additional}</p>
        {order.item.excellent.length > 0 && <><strong>{t.excellent}</strong><ul>
          {order.item.excellent.map((option) => <li key={option}>
            {t.excellentLabels[option as keyof typeof t.excellentLabels] ?? option}
          </li>)}
        </ul></>}
      </article>)}
    </div>}
  </section>;
}
