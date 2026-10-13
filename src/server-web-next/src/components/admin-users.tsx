'use client';
import { useState } from 'react';
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
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  async function setRole(user: AdminUser, role: UserRole) {
    if (role === user.role) return;
    setBusy(user.id);
    setMessage('');
    setError('');
    try {
      await api(`admin/users/${user.id}/role`, 'PATCH', { role });
      setUsers((current) => current.map((item) => item.id === user.id ? { ...item, role } : item));
      setMessage(t.saved);
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
    {error && <p role="alert" className={styles.error}>{error}</p>}
    {message && <p role="status">{message}</p>}
    <div className={styles.tableWrap}><table>
      <thead><tr><th>{t.email}</th><th>{t.nickname}</th><th>{t.role}</th></tr></thead>
      <tbody>{users.map((user) => <tr key={user.id}>
        <td>{user.email}</td><td>{user.nickname}</td>
        <td><select aria-label={`${t.role}: ${user.email}`} value={user.role}
          disabled={busy !== null} onChange={(event) => setRole(user, event.target.value as UserRole)}>
          <option value="player">{t.player}</option><option value="admin">{t.admin}</option>
        </select></td>
      </tr>)}</tbody>
    </table></div>
    <nav className={styles.pager} aria-label={t.title}>
      {initial.page > 1 && <Link href={`/admin/users?page=${initial.page - 1}`}>{t.previous}</Link>}
      <span>{t.page(initial.page, Math.max(1, initial.totalPages))}</span>
      {initial.page < initial.totalPages && <Link href={`/admin/users?page=${initial.page + 1}`}>{t.next}</Link>}
    </nav>
  </section>;
}
