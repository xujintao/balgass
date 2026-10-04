'use client';
import { useState } from 'react';
import { api, ClientError } from './client-api';
import { useDictionary } from './locale-provider';
import { errorMessage } from '@/lib/i18n';

type Account = { name: string; characters: { name: string; level: number }[] };

export function GameAccounts({ initial }: { initial: Account[] }) {
  const t = useDictionary();
  const [accounts, setAccounts] = useState(initial);
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  async function create(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setError('');
    setBusy(true);
    try {
      await api('game/accounts', 'POST', { name, password, passwordConfirmation: confirmation });
      setAccounts(await api<Account[]>('game/accounts'));
      setName(''); setPassword(''); setConfirmation('');
    } catch (failure) {
      setError(failure instanceof ClientError ? errorMessage(t, failure) : t.actionFailed);
    } finally { setBusy(false); }
  }

  return <section className="game-accounts">
    <h1>{t.gameAccounts}</h1>
    <div className="game-account-layout">
      <div>
        <h2>{t.myGameAccounts}</h2>
        {accounts.length === 0 && <p>{t.noGameAccounts}</p>}
        <ul className="game-account-list">{accounts.map((account) => <li key={account.name}>
          <strong>{account.name}</strong>
          {account.characters.length === 0 ? <p>{t.noGameCharacters}</p> : <ul>{account.characters.map((character) => <li key={character.name}>{character.name} · {t.gameCharacterLevel(character.level)}</li>)}</ul>}
        </li>)}</ul>
      </div>
      <form onSubmit={create} className="game-account-form">
        <h2>{t.createGameAccount}</h2>
        <label htmlFor="game-name">{t.gameAccountName}</label>
        <input id="game-name" value={name} onChange={(event) => setName(event.target.value)} minLength={1} maxLength={10} pattern="[!-~]+" required autoComplete="off" />
        <label htmlFor="game-password">{t.gameAccountPassword}</label>
        <input id="game-password" type="password" value={password} onChange={(event) => setPassword(event.target.value)} minLength={1} maxLength={10} required autoComplete="new-password" />
        <label htmlFor="game-confirmation">{t.gameAccountPasswordConfirmation}</label>
        <input id="game-confirmation" type="password" value={confirmation} onChange={(event) => setConfirmation(event.target.value)} minLength={1} maxLength={10} required autoComplete="new-password" />
        {error && <p role="alert" className="error">{error}</p>}
        <button disabled={busy}>{busy ? t.creatingGameAccount : t.createGameAccount}</button>
      </form>
    </div>
  </section>;
}
