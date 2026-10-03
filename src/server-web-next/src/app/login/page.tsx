import { AuthForm } from '@/components/auth-form';
import type { Metadata } from 'next';
export const metadata: Metadata = { title: '登录' };
export default function Login() {
  return <AuthForm />;
}
