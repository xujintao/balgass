import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { describe, expect, it, vi } from 'vitest';

vi.mock('server-only', () => ({}));
vi.mock('next/cache', () => ({ unstable_cache: (fn: unknown) => fn }));
import { getItemCatalog, listItems, parseItemCatalog, type CatalogItem } from '../../src/lib/item-catalog';

const base = new URL('../../../../config/server-game-common/IGCData/', import.meta.url);
async function gameXml() {
  const files = [
    'Skills/IGC_SkillList.xml', 'Items/IGC_ItemList.xml',
    'Items/IGC_ItemSetType.xml', 'Items/IGC_ItemSetOption.xml',
  ];
  const [skills, items, setTypes, setOptions] = await Promise.all(
    files.map((file) => readFile(fileURLToPath(new URL(file, base)), 'utf8')),
  );
  return { skills, items, setTypes, setOptions };
}
describe('item catalog', () => {
  it('parses the real game configuration into catalog entries and calculated levels', async () => {
    const catalog = parseItemCatalog(await gameXml());
    expect(catalog.items.length).toBeGreaterThan(500);
    expect(catalog.sets.length).toBeGreaterThan(20);
    expect(catalog.items.every((item) => item.detail.length === 16)).toBe(true);
    expect(catalog.items.some((item) => item.kind === 'sword' && item.excellent)).toBe(true);
    expect(catalog.items.some((item) => item.kind === 'wing')).toBe(true);
    expect(catalog.sets.some((set) => set.items.length > 0)).toBe(true);
    const kris = catalog.items.find((item) => item.section === 0 && item.index === 0)!;
    expect(kris.name).toBe('Kris');
    expect(kris.nameZh).toBe('波刃剑');
    expect(kris.detail[0]).toMatchObject({
      damage: '6-11', excellentDamage: '36-41',
      requiredStrength: 27, excellentRequiredStrength: 57,
    });
    expect(kris.detail[15]).toMatchObject({ damage: '72-77', excellentDamage: '102-107' });
    const sword = catalog.items.find((item) => item.kind === 'sword' && item.excellent)!;
    expect(sword.detail[15].level).toBe(15);
    expect(sword.detail[15].damage).not.toBe(sword.detail[0].damage);
  });
  it('loads local files and lists a category', async () => {
    vi.stubEnv('GAME_CONFIG_URL', base.href);
    try {
      expect((await getItemCatalog()).items.length).toBeGreaterThan(500);
      const swords = await listItems('sword');
      expect(swords.kind).toBe('sword');
      expect(swords.items.length).toBeGreaterThan(0);
      const english = swords.items.find((item) => 'section' in item && item.section === 0 && item.index === 1);
      const chinese = (await listItems('sword', 'zh-CN')).items.find((item) => 'section' in item && item.section === 0 && item.index === 1);
      const spanish = (await listItems('sword', 'es')).items.find((item) => 'section' in item && item.section === 0 && item.index === 1);
      expect(english?.name).toBe('Short Sword');
      expect(chinese?.name).toBe('短剑');
      expect(spanish?.name).toBe('Short Sword');
      const setsEn = await listItems('set', 'en');
      const setsZh = await listItems('set', 'zh-CN');
      expect(setsEn.items[0].name).not.toBe(setsZh.items[0].name);
      expect((setsEn.items[0] as { items: string[] }).items[0]).not.toBe((setsZh.items[0] as { items: string[] }).items[0]);
      const skilled = (await getItemCatalog()).items.find((item) => item.skill !== '-' && item.skillZh !== '-')!;
      const sameItem = (locale: 'en' | 'zh-CN') => listItems(skilled.kind, locale).then((result) =>
        (result.items as CatalogItem[]).find((item) => item.section === skilled.section && item.index === skilled.index));
      expect((await sameItem('en'))?.skill).toBe(skilled.skill);
      expect((await sameItem('zh-CN'))?.skill).toBe(skilled.skillZh);
    } finally { vi.unstubAllEnvs(); }
  });
  it('rejects malformed XML', async () => {
    const xml = await gameXml();
    expect(() => parseItemCatalog({ ...xml, items: '<!DOCTYPE ItemList><ItemList />' })).toThrow();
  });
});
