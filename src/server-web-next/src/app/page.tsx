import { siteDictionary } from '@/lib/i18n-server';
export default async function Home() {
  const t = await siteDictionary();
  return (
    <section className="card hero">
      <span className="eyebrow">{t.homeEyebrow}</span>
      <h1>{t.homeTitle}</h1>
      <p>{t.homeText}</p>
    </section>
  );
}
