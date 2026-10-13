'use client';
import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { api, ClientError } from './client-api';
import { useLocale } from './locale-provider';
import { adminText } from '@/lib/admin-users-i18n';
import type { AdminUser, UserRole } from '@/lib/admin-users';
import styles from '@/app/admin/users/users.module.css';

type List = { users: AdminUser[]; page: number; pageSize: number; total: number; totalPages: number };

export function AdminUsers({ initial, selfId }: { initial: List; selfId: string }) {
  const t = adminText[useLocale()];
  const router = useRouter();
  const [users, setUsers] = useState(initial.users);
  const [busy, setBusy] = useState<string | null>(null);
  const [announcement, setAnnouncement] = useState('');
  const [error, setError] = useState('');
  const announcementTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => () => {
    if (announcementTimer.current) clearTimeout(announcementTimer.current);
  }, []);

  async function setRole(user: AdminUser, role: UserRole) {
    if (role === user.role) return;
    setBusy(user.id);
    if (announcementTimer.current) clearTimeout(announcementTimer.current);
    setAnnouncement('');
    setError('');
    try {
      await api(`admin/users/${user.id}/role`, 'PATCH', { role });
      setUsers((current) => current.map((item) => item.id === user.id ? { ...item, role } : item));
      setAnnouncement(t.saved);
      announcementTimer.current = setTimeout(() => {
        setAnnouncement('');
        announcementTimer.current = null;
      }, 3000);
      if (user.id === selfId && role === 'player') {
        router.replace('/');
        router.refresh();
      } else router.refresh();
    } catch (cause) {
      setError(cause instanceof ClientError && cause.code === 'LAST_ADMIN' ? t.lastAdmin : t.unavailable);
    } finally {
      setBusy(null);
    }
  }

  return <section className={styles.page}>
    <h1>{t.title}</h1>
    <div className={styles.toast} data-visible={Boolean(announcement)} role="status">{announcement}</div>
    {error && <p role="alert" className={styles.error}>{error}</p>}
    <div className={styles.tableWrap}><table>
      <thead><tr><th>{t.email}</th><th>{t.nickname}</th><th>{t.role}</th><th>{t.action}</th></tr></thead>
      <tbody>{users.map((user) => <tr key={user.id}>
        <td>{user.email}</td><td>{user.nickname}</td>
        <td>{user.role === 'admin' ? t.admin : t.player}</td>
        <td><div className={styles.actionCell}>
          <button type="button" disabled={busy !== null}
            onClick={() => setRole(user, user.role === 'admin' ? 'player' : 'admin')}>
            {user.role === 'admin' ? t.demote : t.promote}
          </button>
        </div></td>
      </tr>)}</tbody>
    </table></div>
    <nav className={styles.pager} aria-label={t.title}>
      {initial.page > 1 && <Link href={`/admin/users?page=${initial.page - 1}`}>{t.previous}</Link>}
      <span>{t.page(initial.page, Math.max(1, initial.totalPages))}</span>
      {initial.page < initial.totalPages && <Link href={`/admin/users?page=${initial.page + 1}`}>{t.next}</Link>}
    </nav>
  </section>;
}
