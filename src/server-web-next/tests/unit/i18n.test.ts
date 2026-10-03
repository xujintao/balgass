import { describe, expect, it } from 'vitest';
import {
  dictionaries,
  errorMessage,
  formatDate,
  formatDateTime,
  resolveLocale,
} from '@/lib/i18n';

describe('site localization', () => {
  it('defaults to English for missing or unsupported cookies', () => {
    expect(resolveLocale(undefined)).toBe('en');
    expect(resolveLocale('invalid')).toBe('en');
    expect(resolveLocale('es')).toBe('es');
    expect(resolveLocale('zh-CN')).toBe('zh-CN');
  });

  it('maps API and browser passkey errors without showing provider text', () => {
    expect(
      errorMessage(dictionaries.en, {
        code: 'EMAIL_ALREADY_REGISTERED',
        message: '中文',
      }),
    ).toBe('This email is already registered. Please log in.');
    expect(
      errorMessage(dictionaries.es, { code: 'PASSKEY_CANCELLED' }),
    ).toContain('canceló');
    expect(
      errorMessage(dictionaries.en, { code: 'WEBAUTHN_VERIFY_FAILED' }),
    ).toContain('Passkey verification failed');
    expect(
      errorMessage(dictionaries.en, { code: 'UNKNOWN', message: '中文' }),
    ).toBe(dictionaries.en.actionFailed);
  });

  it('uses the selected locale to format dates and times', () => {
    const instant = '2026-10-03T12:30:00.000Z';
    expect(formatDate('en', instant)).not.toBe(formatDate('es', instant));
    expect(formatDateTime('en', instant)).not.toBe(
      formatDateTime('zh-CN', instant),
    );
  });
});
