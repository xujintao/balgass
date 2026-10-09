import { describe, expect, it } from 'vitest';
import { itemLocaleFromRequest } from '../../src/lib/item-locale';

function request(headers: Record<string, string>) {
  return new Request('http://localhost:3000/api/v1/items', { headers });
}
describe('item API locale', () => {
  it('uses the browser language cookie and ignores browser Accept-Language', () => {
    expect(itemLocaleFromRequest(request({ Cookie: 'other=1; r2f2-locale=zh-CN', 'Accept-Language': 'es' }))).toBe('zh-CN');
    expect(itemLocaleFromRequest(request({ 'Accept-Language': 'zh-CN' }))).toBe('en');
  });
  it('uses App Accept-Language quality and ignores cookies', () => {
    expect(itemLocaleFromRequest(request({ 'X-Client-Type': 'app', Cookie: 'r2f2-locale=zh-CN', 'Accept-Language': 'es-MX;q=0.5, zh;q=0.9' }))).toBe('zh-CN');
    expect(itemLocaleFromRequest(request({ 'X-Client-Type': 'app', 'Accept-Language': 'es-MX' }))).toBe('es');
    expect(itemLocaleFromRequest(request({ 'X-Client-Type': 'app', 'Accept-Language': 'fr,zh;q=0' }))).toBe('en');
  });
});
