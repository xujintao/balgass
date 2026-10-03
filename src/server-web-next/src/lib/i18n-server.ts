import 'server-only';
import { cookies } from 'next/headers';
import { dictionaries, LOCALE_COOKIE, resolveLocale } from './i18n';

export async function siteLocale() {
  return resolveLocale((await cookies()).get(LOCALE_COOKIE)?.value);
}

export async function siteDictionary() {
  return dictionaries[await siteLocale()];
}
