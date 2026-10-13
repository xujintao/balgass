import type { CSSProperties } from 'react';
import type { Metadata } from 'next';
import Image from 'next/image';
import { notFound } from 'next/navigation';
import { assertAdmin } from '@/lib/admin-users';
import { deployLinks, type DeployTarget } from '@/lib/deploy-links';
import { ApiError } from '@/lib/errors';
import { siteDictionary } from '@/lib/i18n-server';
import { context } from '@/lib/session';
import styles from './deploy.module.css';

export const dynamic = 'force-dynamic';

export async function generateMetadata(): Promise<Metadata> {
  return { title: (await siteDictionary()).deploy };
}

const hotspots: { target: DeployTarget; label: string; x: number; y: number; width: number; height: number }[] = [
  { target: 'vercel', label: 'server-web-next', x: 340, y: 80, width: 120, height: 40 },
  { target: 'supabase', label: 'Supabase', x: 540, y: 80, width: 80, height: 40 },
  { target: 'resend', label: 'Resend', x: 640, y: 80, width: 80, height: 40 },
  { target: 'grafana', label: 'grafana', x: 810, y: 60, width: 80, height: 40 },
  { target: 'pgadminGame', label: 'pgadmin (game.r2f2.com)', x: 720, y: 320, width: 80, height: 40 },
  { target: 'pgadminUs', label: 'pgadmin (us.r2f2.com)', x: 720, y: 580, width: 80, height: 40 },
];

function hotspotStyle({ x, y, width, height }: typeof hotspots[number]): CSSProperties {
  return {
    left: `${x / 962 * 100}%`,
    top: `${y / 722 * 100}%`,
    width: `${width / 962 * 100}%`,
    height: `${height / 722 * 100}%`,
  };
}

export default async function DeployPage() {
  const t = await siteDictionary();
  try {
    await assertAdmin(await context());
  } catch (error) {
    if (error instanceof ApiError && [401, 403].includes(error.status)) notFound();
    return <section><h1>{t.deploy}</h1><p role="alert">{t.deployUnavailable}</p></section>;
  }
  const links = deployLinks();
  return <section className={styles.page}>
    <h1>{t.deploy}</h1>
    <div className={styles.diagram}>
      <Image src="/deploy.drawio.svg" alt={t.deployDiagram} width={962} height={722} unoptimized priority />
      {hotspots.map((spot) => links[spot.target] &&
        <a key={spot.target} className={styles.hotspot} style={hotspotStyle(spot)}
          href={links[spot.target]} target="_blank" rel="noopener noreferrer"
          aria-label={`${spot.label} (${t.openConsole})`} title={spot.label} />)}
    </div>
  </section>;
}
