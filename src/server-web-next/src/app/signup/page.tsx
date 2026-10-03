import { AuthForm } from '@/components/auth-form';
import type { Metadata } from 'next';
import { siteDictionary } from '@/lib/i18n-server';
export async function generateMetadata(): Promise<Metadata> {
  return { title: (await siteDictionary()).signup };
}
export default function Signup() {
  return <AuthForm signup />;
}
