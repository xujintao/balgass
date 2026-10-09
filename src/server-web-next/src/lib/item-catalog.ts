import 'server-only';
import { unstable_cache } from 'next/cache';
import { XMLParser, XMLValidator } from 'fast-xml-parser';
import { ApiError } from './errors';
import { readGameConfigFile } from './game-config';
import type { Locale } from './i18n';

export const itemKindGroups = [
  [
    'sword', 'claw', 'axe', 'mace', 'scepter', 'spear', 'bow', 'crossbow',
    'staff', 'stick', 'book', 'shield',
  ],
  [
    'helmet', 'armor', 'pants', 'gloves', 'boots', 'wing',
    'pendant', 'ring', 'set',
  ],
] as const;
export const itemKinds = [...itemKindGroups[0], ...itemKindGroups[1]] as const;
export type ItemKind = (typeof itemKinds)[number];
const kindByCode: Record<number, Exclude<ItemKind, 'set'>> = {
  1: 'sword', 2: 'sword', 3: 'claw', 4: 'axe', 5: 'mace', 6: 'scepter',
  7: 'spear', 8: 'bow', 9: 'crossbow', 12: 'staff', 13: 'stick',
  14: 'book', 15: 'shield', 16: 'helmet', 17: 'armor', 18: 'pants',
  19: 'gloves', 20: 'boots', 23: 'wing', 24: 'wing', 25: 'wing',
  26: 'wing', 27: 'wing', 28: 'wing', 29: 'pendant', 30: 'pendant', 31: 'ring',
};
type Attributes = Record<string, unknown>;
export type ItemDetail = {
  level: number;
  damage: string;
  excellentDamage: string;
  magicPower: string;
  excellentMagicPower: string;
  defense: number;
  excellentDefense: number;
  defenseRate: number;
  excellentDefenseRate: number;
  requiredStrength: number | '-';
  requiredDexterity: number | '-';
  excellentRequiredStrength: number | '-';
  excellentRequiredDexterity: number | '-';
};
export type CatalogItem = {
  section: number;
  index: number;
  kind: Exclude<ItemKind, 'set'>;
  name: string;
  skill: string;
  twoHand: boolean;
  excellent: boolean;
  dropLevel: number;
  attackSpeed: number;
  moveSpeed: number;
  defense: number;
  defenseRate: number;
  wingKind?: string;
  pendantAttack?: boolean;
  pendantMagicAttack?: boolean;
  detail: ItemDetail[];
};
export type CatalogSet = { index: number; name: string; items: string[] };
export type CatalogItemSource = CatalogItem & { nameZh: string; skillZh: string };
export type CatalogSetSource = CatalogSet & { nameZh: string; itemsZh: string[] };
export type Catalog = { items: CatalogItemSource[]; sets: CatalogSetSource[] };

const parser = new XMLParser({
  ignoreAttributes: false,
  attributeNamePrefix: '@_',
  processEntities: false,
});
function xmlRoot(xml: string, root: string): Attributes {
  if (XMLValidator.validate(xml) !== true || /<!DOCTYPE/i.test(xml))
    throw new Error('Invalid game item XML');
  const value = parser.parse(xml)?.[root];
  if (!value || typeof value !== 'object') throw new Error(`Invalid ${root}`);
  return value as Attributes;
}
function entries(value: unknown): Attributes[] {
  if (value === undefined) return [];
  const values = Array.isArray(value) ? value : [value];
  if (values.some((entry) => !entry || typeof entry !== 'object'))
    throw new Error('Invalid game item entry');
  return values as Attributes[];
}
function str(node: Attributes, key: string): string {
  const value = node[`@_${key}`];
  return value === undefined || value === null ? '' : String(value);
}
function num(node: Attributes, key: string, fallback?: number): number {
  const value = str(node, key);
  if (!value && fallback !== undefined) return fallback;
  if (!/^-?\d+$/.test(value)) throw new Error(`Invalid ${key}`);
  return Number(value);
}
function names(node: Attributes, annotationKey = 'annotation') {
  const base = str(node, 'Name');
  const annotation = str(node, annotationKey);
  return base && annotation ? { en: base, zh: annotation } : null;
}
function itemKey(section: number, index: number) { return `${section}:${index}`; }
function range(a: number, b: number): string { return `${Math.trunc(a)}-${Math.trunc(b)}`; }

function details(node: Attributes, kindCode: number, dropLevel: number): ItemDetail[] {
  const min = num(node, 'DamageMin', 0);
  const max = num(node, 'DamageMax', 0);
  const magic = num(node, 'MagicPower', 0);
  const baseDefense = num(node, 'Defense', 0);
  const baseRate = num(node, 'SuccessfulBlocking', 0);
  const strength = num(node, 'ReqStrength', 0);
  const dexterity = num(node, 'ReqDexterity', 0);
  const levelForDivision = dropLevel || 1;
  return Array.from({ length: 16 }, (_, level) => {
    let damage = '-';
    let excellentDamage = '-';
    let magicPower = '-';
    let excellentMagicPower = '-';
    if ((min > 0 && max > 0) || magic > 0) {
      let delta = level * 3;
      let excellentDelta = delta + min * 25 / levelForDivision + 5;
      if (level >= 10) {
        const extra = (level - 9) * (level - 8) / 2;
        delta += extra;
        excellentDelta += extra;
      }
      if (min > 0 && max > 0) {
        damage = range(min + delta, max + delta);
        excellentDamage = range(min + excellentDelta, max + excellentDelta);
      }
      if (magic > 0) {
        magicPower = `${Math.trunc((magic + delta) / 2 + level * 2)}%`;
        excellentMagicPower = `${Math.trunc((magic + excellentDelta) / 2 + level * 2)}%`;
      }
    }
    let defense = baseDefense;
    let excellentDefense = baseDefense;
    if (baseDefense > 0) {
      let delta: number;
      let excellentDelta: number;
      if (kindCode === 15) {
        delta = level; excellentDelta = delta;
      } else if ([23, 28].includes(kindCode)) {
        delta = level * 3;
        if (level >= 10) delta += (level - 9) * (level - 8) / 2;
        excellentDelta = delta;
      } else if ([24, 27].includes(kindCode)) {
        delta = level * 2;
        if (level >= 10) delta += (level - 9) * (level - 8) / 2 + level - 9;
        excellentDelta = delta;
      } else if (kindCode === 26) {
        delta = level * 2 + 15;
        if (level >= 10) delta += (level - 9) * (level - 8) / 2 + level - 9;
        excellentDelta = delta;
      } else if (kindCode === 25) {
        delta = level * 4;
        if (level >= 10) delta += (level - 9) * (level - 8) / 2;
        excellentDelta = delta;
      } else {
        delta = level * 3;
        excellentDelta = delta + baseDefense * 12 / levelForDivision + dropLevel / 5 + 4;
        if (level >= 10) {
          const extra = (level - 9) * (level - 8) / 2;
          delta += extra; excellentDelta += extra;
        }
      }
      defense = Math.trunc(baseDefense + delta);
      excellentDefense = Math.trunc(baseDefense + excellentDelta);
    }
    let defenseRate = baseRate;
    let excellentDefenseRate = baseRate;
    if (baseRate > 0) {
      let delta = level * 3;
      let excellentDelta = delta + baseRate * 25 / levelForDivision + 5;
      if (level >= 10) {
        const extra = (level - 9) * (level - 8) / 2;
        delta += extra; excellentDelta += extra;
      }
      defenseRate = Math.trunc(baseRate + delta);
      excellentDefenseRate = Math.trunc(baseRate + excellentDelta);
    }
    const required = (base: number, excellent: boolean): number | '-' =>
      base > 0 ? Math.trunc(base * (dropLevel + (excellent ? 25 : 0) + level * 3) * 3 / 100 + 20) : '-';
    return {
      level, damage, excellentDamage, magicPower, excellentMagicPower,
      defense, excellentDefense, defenseRate, excellentDefenseRate,
      requiredStrength: required(strength, false),
      requiredDexterity: required(dexterity, false),
      excellentRequiredStrength: required(strength, true),
      excellentRequiredDexterity: required(dexterity, true),
    };
  });
}

export function parseItemCatalog(xml: {
  skills: string; items: string; setTypes: string; setOptions: string;
}): Catalog {
  const skillsRoot = xmlRoot(xml.skills, 'SkillList');
  const itemsRoot = xmlRoot(xml.items, 'ItemList');
  const typesRoot = xmlRoot(xml.setTypes, 'SetItemType');
  const optionsRoot = xmlRoot(xml.setOptions, 'SetItemOption');
  const skills = new Map(entries(skillsRoot.Skill).map((skill) => [num(skill, 'Index'), names(skill, 'anotation')]));
  const items: CatalogItemSource[] = [];
  const byKey = new Map<string, CatalogItemSource>();
  for (const section of entries(itemsRoot.Section)) {
    const sectionIndex = num(section, 'Index');
    for (const node of entries(section.Item)) {
      const type = num(node, 'Type');
      if (type === 4) continue;
      const itemNames = names(node);
      if (!itemNames) continue;
      const kindCode = num(node, 'KindB');
      const kind = kindByCode[kindCode];
      if (!kind) continue;
      const index = num(node, 'Index');
      const dropLevel = num(node, 'DropLevel', 0) || num(node, 'ReqLevel', 0);
      const itemSkill = skills.get(num(node, 'SkillIndex'));
      const item: CatalogItemSource = {
        section: sectionIndex, index, kind,
        name: itemNames.en, nameZh: itemNames.zh,
        skill: itemSkill?.en ?? '-', skillZh: itemSkill?.zh ?? '-',
        twoHand: str(node, 'TwoHand') === '1',
        excellent: str(node, 'Option') === '1' && ![2, 6, 7].includes(type),
        dropLevel,
        attackSpeed: num(node, 'AttackSpeed', 0),
        moveSpeed: num(node, 'WalkSpeed', 0),
        defense: num(node, 'Defense', 0),
        defenseRate: num(node, 'SuccessfulBlocking', 0),
        detail: details(node, kindCode, dropLevel),
      };
      if (kind === 'wing')
        item.wingKind = ({ 23: '1D', 24: '2D', 25: '3D', 26: '2D', 27: '2D', 28: '2D5' } as Record<number, string>)[kindCode];
      if (kind === 'pendant') {
        item.pendantAttack = kindCode === 29;
        item.pendantMagicAttack = kindCode === 30;
      }
      items.push(item);
      byKey.set(itemKey(sectionIndex, index), item);
    }
  }
  const setMembers = new Map<number, Set<string>>();
  for (const section of entries(typesRoot.Section)) {
    const sectionIndex = num(section, 'Index');
    for (const node of entries(section.Item)) {
      const key = itemKey(sectionIndex, num(node, 'Index'));
      for (const tier of ['TierI', 'TierII']) {
        const id = num(node, tier, 0);
        if (id === 0) continue;
        const members = setMembers.get(id) ?? new Set<string>();
        members.add(key);
        setMembers.set(id, members);
      }
    }
  }
  const sets: CatalogSetSource[] = [];
  for (const node of entries(optionsRoot.SetItem)) {
    const index = num(node, 'Index');
    if (index > 38) break;
    const setNames = names(node);
    if (!setNames) continue;
    const members = [...(setMembers.get(index) ?? [])].sort((a, b) => {
      const [as, ai] = a.split(':').map(Number);
      const [bs, bi] = b.split(':').map(Number);
      return as - bs || ai - bi;
    }).map((key) => byKey.get(key)).filter((value): value is CatalogItemSource => Boolean(value));
    sets.push({
      index, name: setNames.en, nameZh: setNames.zh,
      items: members.map((item) => item.name),
      itemsZh: members.map((item) => item.nameZh),
    });
  }
  return { items, sets };
}

async function loadCatalog(): Promise<Catalog> {
  const [skills, items, setTypes, setOptions] = await Promise.all([
    readGameConfigFile('Skills/IGC_SkillList.xml'),
    readGameConfigFile('Items/IGC_ItemList.xml'),
    readGameConfigFile('Items/IGC_ItemSetType.xml'),
    readGameConfigFile('Items/IGC_ItemSetOption.xml'),
  ]);
  return parseItemCatalog({ skills, items, setTypes, setOptions });
}
const cachedCatalog = unstable_cache(loadCatalog, ['item-catalog'], { revalidate: 60 });
export async function getItemCatalog(): Promise<Catalog> {
  try {
    return process.env.GAME_CONFIG_URL?.startsWith('file:')
      ? await loadCatalog()
      : await cachedCatalog();
  } catch (error) {
    console.error('Unable to load item catalog:', error);
    throw new ApiError(503, 'ITEM_CATALOG_UNAVAILABLE', '道具目录暂不可用。');
  }
}
export function parseItemKind(value: string): ItemKind {
  if (!itemKinds.includes(value as ItemKind))
    throw new ApiError(400, 'INVALID_INPUT', '道具分类无效。');
  return value as ItemKind;
}
export async function listItems(kind: ItemKind = 'sword', locale: Locale = 'en'): Promise<{ kind: ItemKind; items: CatalogItem[] | CatalogSet[] }> {
  const catalog = await getItemCatalog();
  if (kind === 'set') return { kind, items: catalog.sets.map(({ nameZh, itemsZh, ...set }) => ({
    ...set, name: locale === 'zh-CN' ? nameZh : set.name,
    items: locale === 'zh-CN' ? itemsZh : set.items,
  })) };
  const items = catalog.items.filter((item) => item.kind === kind);
  items.sort((a, b) => kind === 'wing'
    ? (a.wingKind ?? '').localeCompare(b.wingKind ?? '')
    : a.dropLevel - b.dropLevel);
  return { kind, items: items.map(({ nameZh, skillZh, ...item }) => ({
    ...item, name: locale === 'zh-CN' ? nameZh : item.name,
    skill: locale === 'zh-CN' ? skillZh : item.skill,
  })) };
}
export async function findItem(section: number, index: number): Promise<CatalogItemSource> {
  const item = (await getItemCatalog()).items.find((entry) => entry.section === section && entry.index === index);
  if (!item) throw new ApiError(404, 'ITEM_NOT_FOUND', '道具不存在。');
  return item;
}
