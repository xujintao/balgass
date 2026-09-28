-- 仅在已经执行旧版 profiles SQL 的项目上运行。
-- 保留现有用户、昵称和改名时间；停止数据库业务函数和注册触发器。
-- 新版 Next.js 需要配置服务端 Secret Key，首次登录会补齐缺少的资料。
begin;

drop trigger if exists auth_user_create_profile on auth.users;
drop policy if exists profiles_read_self on public.profiles;
alter table public.profiles drop constraint if exists profiles_nickname_check;
alter table public.profiles alter column nickname_key drop expression if exists;

drop function if exists public.change_nickname(text);
drop function if exists public.create_profile();
drop function if exists public.has_verified_email();
drop function if exists public.normalize_nickname(text);
drop function if exists public.valid_nickname(text);

alter table public.profiles enable row level security;
revoke all on public.profiles from public, anon, authenticated;
grant select on public.profiles to authenticated;
grant select, insert, update, delete on public.profiles to service_role;
create policy profiles_read_self on public.profiles
  for select to authenticated
  using (user_id = (select auth.uid()));

commit;
notify pgrst, 'reload schema';
