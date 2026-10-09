import { beforeEach, describe, expect, it, vi } from 'vitest';
vi.mock('server-only', () => ({}));
const mock = vi.hoisted(() => ({
  findItem: vi.fn(),
  insert: vi.fn(),
  select: vi.fn(),
  single: vi.fn(),
  listSelect: vi.fn(),
  eq: vi.fn(),
  order: vi.fn(),
}));
vi.mock('../../src/lib/item-catalog', () => ({ findItem: mock.findItem }));
vi.mock('../../src/lib/supabase-admin', () => ({
  supabaseAdmin: () => ({ from: () => ({ insert: mock.insert }) }),
}));
import { createOrder, listOrders } from '../../src/lib/item-orders';

const user = { id: 'user-a', email_confirmed_at: '2026-01-01' };
const ctx = { user, client: { from: () => ({ select: mock.listSelect }) } } as never;
const input = { section: 0, index: 1, level: 15, excellent: [], additional: 16 };
beforeEach(() => {
  vi.clearAllMocks();
  mock.findItem.mockResolvedValue({ section: 0, index: 1, name: 'Test Sword', nameZh: '测试剑', excellent: false });
  mock.insert.mockReturnValue({ select: mock.select });
  mock.select.mockReturnValue({ single: mock.single });
  mock.single.mockResolvedValue({ data: {
    id: 'order-1', status: 'PENDING', item_section: 0, item_index: 1,
    item_name_en: 'Test Sword', item_name_zh: '测试剑', level: 15, excellent: [], additional: 16,
    created_at: '2026-10-08T00:00:00Z', updated_at: '2026-10-08T00:00:00Z',
  }, error: null });
  mock.listSelect.mockReturnValue({ eq: mock.eq });
  mock.eq.mockReturnValue({ order: mock.order });
  mock.order.mockResolvedValue({ data: [], error: null });
});
describe('item orders', () => {
  it('creates a pending single item order owned by the verified user', async () => {
    const order = await createOrder(ctx, input, 'zh-CN');
    expect(order).toMatchObject({ id: 'order-1', status: 'PENDING', item: { name: '测试剑', quantity: 1 } });
    expect(mock.insert).toHaveBeenCalledWith(expect.objectContaining({
      user_id: user.id, item_name_en: 'Test Sword', item_name_zh: '测试剑', level: 15, additional: 16,
    }));
  });
  it('rejects forged ownership, invalid options and excellent options on ordinary items', async () => {
    await expect(createOrder(ctx, { ...input, user_id: 'other' })).rejects.toThrow();
    await expect(createOrder(ctx, { ...input, level: 16 })).rejects.toThrow();
    await expect(createOrder(ctx, { ...input, excellent: ['excellent_attack_rate'] })).rejects.toMatchObject({ status: 400 });
    expect(mock.insert).not.toHaveBeenCalled();
  });
  it('filters reads by the authenticated user', async () => {
    expect(await listOrders(ctx)).toEqual([]);
    expect(mock.eq).toHaveBeenCalledWith('user_id', user.id);
    expect(mock.order).toHaveBeenCalledWith('created_at', { ascending: false });
  });
  it('uses stored translations for each requested language', async () => {
    const row = (await mock.single()).data;
    mock.order.mockResolvedValue({ data: [row], error: null });
    expect((await listOrders(ctx, 'es'))[0].item.name).toBe('Test Sword');
    expect((await listOrders(ctx, 'zh-CN'))[0].item.name).toBe('测试剑');
    expect((await listOrders(ctx, 'en'))[0].item.name).toBe('Test Sword');
  });
});
