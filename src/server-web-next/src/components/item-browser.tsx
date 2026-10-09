'use client';
import { useState } from 'react';
import Link from 'next/link';
import { api, ClientError } from './client-api';
import { useDictionary, useLocale } from './locale-provider';
import { errorMessage } from '@/lib/i18n';
import { itemDictionaries } from '@/lib/item-i18n';
import { excellentOptions } from '@/lib/item-options';
import type { CatalogItem, CatalogSet, ItemKind } from '@/lib/item-catalog';
import styles from '@/app/items/items.module.css';

type Result = { kind: ItemKind; items: CatalogItem[] | CatalogSet[] };
function ItemCard({ item }: { item: CatalogItem }) {
  const t = itemDictionaries[useLocale()];
  const common = useDictionary();
  const [level, setLevel] = useState(0);
  const [additional, setAdditional] = useState(0);
  const [excellent, setExcellent] = useState<string[]>([]);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const showAttack = ['sword', 'claw', 'axe', 'mace', 'scepter', 'spear', 'bow', 'crossbow', 'staff', 'stick', 'book'].includes(item.kind);
  const showDefense = ['shield', 'helmet', 'armor', 'pants', 'gloves', 'boots', 'wing'].includes(item.kind);
  async function submit(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(''); setSuccess(''); setBusy(true);
    try {
      await api('orders', 'POST', { section: item.section, index: item.index, level, excellent, additional });
      setSuccess(t.submitted);
    } catch (failure) {
      if (failure instanceof ClientError && failure.code in t.errors)
        setError(t.errors[failure.code as keyof typeof t.errors]);
      else setError(failure instanceof ClientError ? errorMessage(common, failure) : common.actionFailed);
    } finally { setBusy(false); }
  }
  return <details className={styles.item}>
    <summary><strong>{item.name}</strong><span>{t.level} {item.dropLevel}</span></summary>
    <div className={styles.itemBody}>
      <dl className={styles.facts}>
        {item.skill !== '-' && <><dt>{t.skill}</dt><dd>{item.skill}</dd></>}
        {item.twoHand && <><dt>{t.twoHand}</dt><dd>{t.yes}</dd></>}
        {item.attackSpeed > 0 && <><dt>{t.attackSpeed}</dt><dd>{item.attackSpeed}</dd></>}
        {item.wingKind && <><dt>{t.wingKind}</dt><dd>{item.wingKind}</dd></>}
        {item.kind === 'pendant' && <><dt>{t.attack}</dt><dd>{item.pendantAttack ? t.yes : t.no}</dd><dt>{t.magicPower}</dt><dd>{item.pendantMagicAttack ? t.yes : t.no}</dd></>}
      </dl>
      {(showAttack || showDefense) && <div className={styles.tableScroll}><table>
        <thead><tr><th>{t.level}</th>
          {showAttack && <><th>{t.attack}</th><th>{t.magicPower}</th></>}
          {showDefense && <><th>{t.defense}</th>{item.kind === 'shield' && <th>{t.defenseRate}</th>}</>}
          {item.kind !== 'wing' && <th>{t.requirements}</th>}
        </tr></thead>
        <tbody>{item.detail.map((detail) => <tr key={detail.level}>
          <th>{detail.level}</th>
          {showAttack && <><td>{detail.damage}{item.excellent && ` (${detail.excellentDamage})`}</td><td>{detail.magicPower}{item.excellent && ` (${detail.excellentMagicPower})`}</td></>}
          {showDefense && <><td>{detail.defense}{item.excellent && ` (${detail.excellentDefense})`}</td>{item.kind === 'shield' && <td>{detail.defenseRate}{item.excellent && ` (${detail.excellentDefenseRate})`}</td>}</>}
          {item.kind !== 'wing' && <td>{detail.requiredStrength}/{detail.requiredDexterity}{item.excellent && ` (${detail.excellentRequiredStrength}/${detail.excellentRequiredDexterity})`}</td>}
        </tr>)}</tbody>
      </table></div>}
      <form onSubmit={submit} className={styles.orderForm}>
        <p>{t.pendingNotice}</p>
        <div className={styles.formRow}>
          <label>{t.level}<select value={level} onChange={(event) => setLevel(Number(event.target.value))}>
            {Array.from({ length: 16 }, (_, value) => <option key={value} value={value}>+{value}</option>)}
          </select></label>
          <label>{t.additional}<select value={additional} onChange={(event) => setAdditional(Number(event.target.value))}>
            {[0, 4, 8, 12, 16].map((value) => <option key={value} value={value}>+{value}</option>)}
          </select></label>
        </div>
        {item.excellent && <fieldset><legend>{t.excellent}</legend><div className={styles.checks}>
          {excellentOptions.map((option) => <label key={option}><input type="checkbox" checked={excellent.includes(option)} onChange={(event) => setExcellent((current) => event.target.checked ? [...current, option] : current.filter((value) => value !== option))} />{t.excellentLabels[option]}</label>)}
        </div></fieldset>}
        <button disabled={busy}>{busy ? t.submitting : t.submit}</button>
        {error && <p className="error" role="alert">{error}</p>}
        {success && <p className="success" role="status">{success} <Link href="/orders">{t.orders}</Link></p>}
      </form>
    </div>
  </details>;
}
export function ItemBrowser({ result }: { result: Result }) {
  const t = itemDictionaries[useLocale()];
  if (result.items.length === 0) return <p>{t.noItems}</p>;
  if (result.kind === 'set') return <div className={styles.list}>{(result.items as CatalogSet[]).map((set) => <article className={styles.set} key={set.index}>
    <h2>{set.name}</h2><strong>{t.setPieces}</strong><ul>{set.items.map((name, index) => <li key={`${name}-${index}`}>{name}</li>)}</ul>
  </article>)}</div>;
  return <div className={styles.list}>{(result.items as CatalogItem[]).map((item) => <ItemCard key={`${item.section}-${item.index}`} item={item} />)}</div>;
}
