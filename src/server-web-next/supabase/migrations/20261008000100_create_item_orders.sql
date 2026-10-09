-- 在 Supabase Dashboard SQL Editor 中按文件名顺序执行。
begin;

lock table public.migrations in exclusive mode;

do $$
begin
  if not exists (select 1 from public.migrations where id = '20261008000100') then
    create table public.item_orders (
      id uuid primary key default gen_random_uuid(),
      user_id uuid not null references auth.users(id) on delete cascade,
      status text not null default 'PENDING' check (status in ('PENDING', 'PAID', 'SHIPPING', 'SHIPPED', 'COMPLETED', 'CANCELLED')),
      item_section integer not null check (item_section >= 0),
      item_index integer not null check (item_index >= 0),
      item_name_en text not null check (length(item_name_en) > 0),
      item_name_zh text not null check (length(item_name_zh) > 0),
      level integer not null check (level between 0 and 15),
      excellent text[] not null default '{}',
      additional integer not null check (additional in (0, 4, 8, 12, 16)),
      created_at timestamptz not null default now(),
      updated_at timestamptz not null default now(),
      check (excellent <@ array[
        'excellent_attack_rate', 'excellent_attack_level', 'excellent_attack_percent',
        'excellent_attack_speed', 'excellent_attack_hp', 'excellent_attack_mp'
      ]::text[])
    );

    create index item_orders_user_created_idx on public.item_orders (user_id, created_at desc);
    alter table public.item_orders enable row level security;
    revoke all on public.item_orders from public, anon, authenticated;
    grant select on public.item_orders to authenticated;
    grant select, insert, update, delete on public.item_orders to service_role;
    create policy item_orders_read_self on public.item_orders
      for select to authenticated using (user_id = (select auth.uid()));

    insert into public.migrations (id, filename)
      values ('20261008000100', '20261008000100_create_item_orders.sql');
  end if;
end $$;

notify pgrst, 'reload schema';
commit;
