import { LOCALE_COOKIE, resolveLocale, type Locale } from './i18n';

function languageTag(value: string): Locale | null {
  const tag = value.toLowerCase();
  if (tag === 'zh' || tag.startsWith('zh-')) return 'zh-CN';
  if (tag === 'es' || tag.startsWith('es-')) return 'es';
  if (tag === 'en' || tag.startsWith('en-')) return 'en';
  return null;
}

export function itemLocaleFromRequest(request: Request): Locale {
  if (request.headers.get('x-client-type') !== 'app') {
    const cookie = request.headers.get('cookie')?.split(';').map((part) => part.trim())
      .find((part) => part.startsWith(`${LOCALE_COOKIE}=`));
    return resolveLocale(cookie?.slice(LOCALE_COOKIE.length + 1));
  }
  const choices = (request.headers.get('accept-language') ?? '').split(',').map((part, index) => {
    const [tag, ...parameters] = part.trim().split(';');
    const quality = parameters.find((value) => value.trim().startsWith('q='));
    const weight = quality ? Number(quality.trim().slice(2)) : 1;
    return { locale: languageTag(tag.trim()), weight, index };
  }).filter((choice) => choice.locale && Number.isFinite(choice.weight) && choice.weight > 0 && choice.weight <= 1)
    .sort((a, b) => b.weight - a.weight || a.index - b.index);
  return choices[0]?.locale ?? 'en';
}
