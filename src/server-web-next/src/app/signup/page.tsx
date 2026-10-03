import { AuthForm } from '@/components/auth-form';
import type { Metadata } from 'next';
export const metadata: Metadata = { title: '注册' };
export default function Signup() {
  return <AuthForm signup />;
}
