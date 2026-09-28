-- 新项目：在 Supabase Dashboard 的 SQL Editor 中执行本文件。
-- 随机昵称、格式校验和 30 天改名限制由 Next.js 服务端处理。
begin;

create table if not exists public.profiles (
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

drop policy if exists profiles_read_self on public.profiles;
create policy profiles_read_self on public.profiles
  for select to authenticated
  using (user_id = (select auth.uid()));

commit;
notify pgrst, 'reload schema';
