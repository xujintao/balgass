create schema if not exists extensions;
create extension if not exists pgcrypto with schema extensions;

-- Han ranges generated from Unicode Script=Han (Node 24); checked against web validation.
create function public.normalize_nickname(value text) returns text
language sql immutable strict set search_path = '' as $$
  select normalize(btrim(value, E' \t\n\r\f' || chr(11) || chr(160) || chr(5760) || chr(8192) || chr(8193) || chr(8194) || chr(8195) || chr(8196) || chr(8197) || chr(8198) || chr(8199) || chr(8200) || chr(8201) || chr(8202) || chr(8232) || chr(8233) || chr(8239) || chr(8287) || chr(12288) || chr(65279)), NFC)
$$;
create function public.valid_nickname(value text) returns boolean
language sql immutable strict set search_path = '' as $$
 select char_length(value) between 2 and 24 and value !~ '[^A-Za-z0-9_⺀-⺙⺛-⻳⼀-⿕々〇〡-〩〸-〻㐀-䶿一-鿿豈-舘並-龎𖿢-𖿣𖿰-𖿶𠀀-𪛟𪜀-𫠝𫠠-𬺭𬺰-𮯠𮯰-𮹝丽-𪘀𰀀-𱍊𱍐-𳑹]'
$$;
create table public.profiles (
 user_id uuid primary key references auth.users(id) on delete cascade,
 nickname text not null check (nickname = public.normalize_nickname(nickname) and public.valid_nickname(nickname)),
 nickname_key text generated always as (translate(nickname,'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')) stored,
 created_at timestamptz not null default now(),
 nickname_changed_at timestamptz,
 constraint profiles_nickname_unique unique (nickname_key)
);
alter table public.profiles enable row level security;
revoke all on public.profiles from public, anon, authenticated;
grant select on public.profiles to authenticated;
-- auth.users is private; the helper exposes only the caller verification status.
create function public.has_verified_email() returns boolean
language sql stable security definer set search_path = '' as $$
 select exists(select 1 from auth.users where id = auth.uid() and email_confirmed_at is not null)
$$;
revoke all on function public.has_verified_email() from public;
grant execute on function public.has_verified_email() to authenticated;
create policy profiles_read_self on public.profiles for select to authenticated
 using (user_id = (select auth.uid()) and (select public.has_verified_email()));

create function public.create_profile() returns trigger
language plpgsql security definer set search_path = '' as $$
begin
 loop
  begin
   insert into public.profiles(user_id,nickname) values(new.id, '玩家_' || encode(extensions.gen_random_bytes(8),'hex'));
   exit;
  exception when unique_violation then
   -- Only retry nickname collisions; never mask a duplicate user ID.
   if exists(select 1 from public.profiles where user_id = new.id) then raise; end if;
  end;
 end loop;
 return new;
end;
$$;
revoke all on function public.create_profile() from public, anon, authenticated;
create trigger auth_user_create_profile after insert on auth.users for each row execute function public.create_profile();

create function public.change_nickname(requested_nickname text) returns public.profiles
language plpgsql security definer set search_path = '' as $$
declare
 profile public.profiles;
 value text := public.normalize_nickname(requested_nickname);
 changed_at timestamptz;
begin
 if auth.uid() is null or not public.has_verified_email() then raise exception using message = 'UNAUTHENTICATED', errcode = 'P0001'; end if;
 if value is null or not public.valid_nickname(value) then raise exception using message = 'INVALID_NICKNAME', errcode = 'P0001'; end if;
 select * into strict profile from public.profiles where user_id = auth.uid() for update;
 if profile.nickname = value then return profile; end if;
 changed_at := clock_timestamp();
 if profile.nickname_changed_at is not null and changed_at < profile.nickname_changed_at + interval '720 hours' then
  raise exception using message = 'NICKNAME_COOLDOWN', detail = to_char((profile.nickname_changed_at + interval '720 hours') at time zone 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'), errcode = 'P0001';
 end if;
 begin
  update public.profiles set nickname = value, nickname_changed_at = changed_at where user_id = auth.uid() returning * into profile;
 exception when unique_violation then raise exception using message = 'NICKNAME_TAKEN', errcode = 'P0001';
 end;
 return profile;
end;
$$;
revoke all on function public.change_nickname(text) from public, anon;
grant execute on function public.change_nickname(text) to authenticated;
