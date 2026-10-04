import type { Metadata } from 'next';
import { context } from '@/lib/session';
import { listGameAccounts } from '@/lib/game-accounts';
import { ApiError } from '@/lib/errors';
import { SessionRecovery } from '@/components/session-recovery';
import { GameAccounts } from '@/components/game-accounts';
import { siteDictionary } from '@/lib/i18n-server';
import { errorMessage } from '@/lib/i18n';

export const dynamic = 'force-dynamic';
export async function generateMetadata(): Promise<Metadata> {
  return { title: (await siteDictionary()).gameAccounts };
}

export default async function GameAccountsPage() {
  const t = await siteDictionary();
  let result: { kind: 'ok'; accounts: Awaited<ReturnType<typeof listGameAccounts>> } | { kind: 'unauthenticated' } | { kind: 'error'; message: string };
  try {
    result = { kind: 'ok', accounts: await listGameAccounts(await context()) };
  } catch (error) {
    result = error instanceof ApiError && error.status === 401
      ? { kind: 'unauthenticated' }
      : { kind: 'error', message: error instanceof ApiError ? errorMessage(t, error) : t.serviceUnavailable };
  }
  if (result.kind === 'unauthenticated') return <SessionRecovery />;
  if (result.kind === 'error') return <section className="game-accounts"><h1>{t.gameAccounts}</h1><p role="alert" className="error">{result.message}</p></section>;
  return <GameAccounts initial={result.accounts} />;
}
