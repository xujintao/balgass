import Link from 'next/link';
import { cookies } from 'next/headers';
import { isAdmin } from '@/lib/admin-users';
import { context, ACCESS_COOKIE } from '@/lib/session';
import { siteDictionary } from '@/lib/i18n-server';

export async function AdminNavLink() {
  if (!(await cookies()).has(ACCESS_COOKIE)) return null;
  let visible = false;
  try {
    const ctx = await context();
    visible = await isAdmin(ctx.user.id);
  } catch {
    // Authentication and profile errors do not expose the admin link.
  }
  return visible ? <Link href="/admin/users">{(await siteDictionary()).adminUsers}</Link> : null;
}
