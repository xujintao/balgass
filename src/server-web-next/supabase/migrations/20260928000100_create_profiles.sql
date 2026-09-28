-- 在 Supabase Dashboard SQL Editor 中按文件名顺序执行。
-- 昵称业务规则由 Next.js 服务端处理；DO 是一次性执行块，不创建存储函数。
begin;

create table if not exists public.migrations (
  id text primary key,
  filename text not null unique,
  executed_at timestamptz not null default now()
);

alter table public.migrations enable row level security;
revoke all on public.migrations from public, anon, authenticated, service_role;
-- 迁移历史只由 SQL Editor 的数据库管理员维护。
lock table public.migrations in exclusive mode;

do $$
begin
  if not exists (
    select 1 from public.migrations where id = '20260928000100'
  ) then
    create table public.profiles (
      user_id uuid primary key references auth.users(id) on delete cascade,
      nickname text not null,
      nickname_key text not null unique,
      created_at timestamptz not null default now(),
      nickname_changed_at timestamptz
    );

    comment on column public.profiles.user_id is 'Supabase 认证用户 ID';
    comment on column public.profiles.nickname is '用户显示昵称';
    comment on column public.profiles.nickname_key is 'Next.js 生成的小写昵称，用于唯一性检查';
    comment on column public.profiles.nickname_changed_at is '最后一次用户改名时间；初始随机昵称不计入改名';

    -- 客户端只能读取自己的资料，不能直接插入、改名或修改时间。
    alter table public.profiles enable row level security;
    revoke all on public.profiles from public, anon, authenticated;
    grant select on public.profiles to authenticated;
    grant select, insert, update, delete on public.profiles to service_role;

    create policy profiles_read_self on public.profiles
      for select to authenticated
      using (user_id = (select auth.uid()));

    insert into public.migrations (id, filename)
      values ('20260928000100', '20260928000100_create_profiles.sql');
  end if;
end $$;

notify pgrst, 'reload schema';
commit;
