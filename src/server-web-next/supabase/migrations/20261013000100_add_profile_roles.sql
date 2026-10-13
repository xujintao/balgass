-- 在 Supabase Dashboard SQL Editor 中按文件名顺序执行。
begin;
lock table public.migrations in exclusive mode;

do $$
begin
  if not exists (select 1 from public.migrations where id = '20261013000100') then
    alter table public.profiles add column role text not null default 'player'
      check (role in ('player', 'admin'));

    update public.profiles p set role = 'admin'
    where p.user_id = (
      select p2.user_id from public.profiles p2
      join auth.users u on u.id = p2.user_id
      where u.email_confirmed_at is not null
      order by p2.created_at, p2.user_id limit 1
    );

    create index profiles_created_user_idx on public.profiles (created_at, user_id);
    insert into public.migrations (id, filename)
      values ('20261013000100', '20261013000100_add_profile_roles.sql');
  end if;
end $$;

notify pgrst, 'reload schema';
commit;
