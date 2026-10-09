import 'server-only';
import { z } from 'zod';
import { ApiError } from './errors';
import { findItem } from './item-catalog';
import { excellentOptions } from './item-options';
import type { Locale } from './i18n';
import type { context } from './session';
import { supabaseAdmin } from './supabase-admin';

export const orderInput = z.strictObject({
  section: z.number().int().min(0),
  index: z.number().int().min(0),
  level: z.number().int().min(0).max(15),
  excellent: z.array(z.enum(excellentOptions)).max(excellentOptions.length).refine(
    (values) => new Set(values).size === values.length,
  ),
  additional: z.union([z.literal(0), z.literal(4), z.literal(8), z.literal(12), z.literal(16)]),
});
export type OrderInput = z.infer<typeof orderInput>;

type OrderRow = {
  id: string;
  status: string;
  item_section: number;
  item_index: number;
  item_name_en: string;
  item_name_zh: string;
  level: number;
  excellent: string[];
  additional: number;
  created_at: string;
  updated_at: string;
};
function present(row: OrderRow, locale: Locale) {
  const name = locale === 'zh-CN' ? row.item_name_zh : row.item_name_en;
  return {
    id: row.id,
    status: row.status,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
    item: {
      section: row.item_section,
      index: row.item_index,
      name,
      level: row.level,
      excellent: row.excellent,
      additional: row.additional,
      quantity: 1,
    },
  };
}
export type ItemOrder = ReturnType<typeof present>;

export async function listOrders(ctx: Awaited<ReturnType<typeof context>>, locale: Locale = 'en'): Promise<ItemOrder[]> {
  const { data, error } = await ctx.client.from('item_orders')
    .select('id,status,item_section,item_index,item_name_en,item_name_zh,level,excellent,additional,created_at,updated_at')
    .eq('user_id', ctx.user.id)
    .order('created_at', { ascending: false });
  if (error) throw new ApiError(503, 'ORDERS_UNAVAILABLE', '订单暂不可用。');
  return (data as OrderRow[]).map((row) => present(row, locale));
}

export async function createOrder(ctx: Awaited<ReturnType<typeof context>>, value: unknown, locale: Locale = 'en'): Promise<ItemOrder> {
  const input = orderInput.parse(value);
  const item = await findItem(input.section, input.index);
  if (!item.excellent && input.excellent.length)
    throw new ApiError(400, 'INVALID_INPUT', '该道具不支持卓越属性。');
  const { data, error } = await supabaseAdmin().from('item_orders').insert({
    user_id: ctx.user.id,
    item_section: item.section,
    item_index: item.index,
    item_name_en: item.name,
    item_name_zh: item.nameZh,
    level: input.level,
    excellent: input.excellent,
    additional: input.additional,
  }).select('id,status,item_section,item_index,item_name_en,item_name_zh,level,excellent,additional,created_at,updated_at').single();
  if (error || !data) throw new ApiError(503, 'ORDERS_UNAVAILABLE', '暂时无法提交订单。');
  return present(data as OrderRow, locale);
}
