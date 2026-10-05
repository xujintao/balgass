import type { Metadata } from 'next';
import { connection } from 'next/server';
import { GameMapViewer } from '@/components/game-map-viewer';
import { gameWebSocketUrl, getGameMapNames } from '@/lib/game-config';
import { siteDictionary } from '@/lib/i18n-server';
import styles from './game.module.css';

export async function generateMetadata(): Promise<Metadata> {
  const t = await siteDictionary();
  return { title: t.gameMap };
}

export default async function GamePage() {
  await connection();
  const t = await siteDictionary();
  let maps: string[];
  try {
    maps = await getGameMapNames();
    if (maps.length === 0) throw new Error('No game maps configured');
  } catch (error) {
    console.error('Unable to load game maps:', error);
    return (
      <section className={styles.page}>
        <h1>{t.gameMap}</h1>
        <p className="error" role="alert">{t.mapUnavailable}</p>
      </section>
    );
  }

  let websocketUrl: string | null = null;
  try {
    websocketUrl = gameWebSocketUrl();
  } catch (error) {
    console.error('Unable to configure game WebSocket:', error);
  }

  return (
    <section className={styles.page}>
      <h1>{t.gameMap}</h1>
      <GameMapViewer maps={maps} websocketUrl={websocketUrl} />
    </section>
  );
}
