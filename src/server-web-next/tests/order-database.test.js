import assert from 'node:assert/strict';
import { readFile, mkdtemp, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { createServer } from 'node:net';
import { randomUUID } from 'node:crypto';
import { Client } from 'pg';
import EmbeddedPostgres from 'embedded-postgres';
import { test } from 'vitest';

test('订单迁移、约束和 RLS', async () => {
  const socket = createServer();
  await new Promise((resolve) => socket.listen(0, '127.0.0.1', resolve));
  const port = socket.address().port;
  await new Promise((resolve) => socket.close(resolve));
  const directory = await mkdtemp(join(tmpdir(), 'balgass-orders-pg-'));
  const embedded = new EmbeddedPostgres({
    databaseDir: join(directory, 'data'), user: 'postgres', password: 'test-password',
    port, persistent: false, onLog: () => {}, onError: () => {},
  });
  let db;
  const clients = [];
  try {
    await embedded.initialise();
    await embedded.start();
    db = new Client({ connectionString: `postgresql://postgres:test-password@127.0.0.1:${port}/postgres` });
    await db.connect();
    await db.query(`do $$ begin
      if not exists(select 1 from pg_roles where rolname='anon') then create role anon; end if;
      if not exists(select 1 from pg_roles where rolname='authenticated') then create role authenticated; end if;
      if not exists(select 1 from pg_roles where rolname='service_role') then create role service_role bypassrls; end if;
    end $$;
    create schema auth;
    create table auth.users(id uuid primary key);
    create function auth.uid() returns uuid language sql stable as $$ select nullif(current_setting('request.jwt.claim.sub',true),'')::uuid $$;
    grant usage on schema auth to authenticated, anon;
    grant execute on function auth.uid() to authenticated, anon;`);
    for (const file of [
      '20260928000100_create_profiles.sql',
      '20261008000100_create_item_orders.sql',
    ]) await db.query(await readFile(new URL(`../supabase/migrations/${file}`, import.meta.url), 'utf8'));
    const migration = await readFile(new URL('../supabase/migrations/20261008000100_create_item_orders.sql', import.meta.url), 'utf8');
    await db.query(migration);
    assert.equal((await db.query("select count(*)::integer as n from public.migrations where id='20261008000100'")).rows[0].n, 1);
    const owner = randomUUID(), other = randomUUID();
    await db.query('insert into auth.users(id) values($1),($2)', [owner, other]);
    const inserted = await db.query(
      "insert into public.item_orders(user_id,item_section,item_index,item_name_en,item_name_zh,level,excellent,additional) values($1,0,1,'Test Sword','测试剑',15,array['excellent_attack_rate'],16) returning id",
      [owner],
    );
    await db.query(
      "insert into public.item_orders(user_id,item_section,item_index,item_name_en,item_name_zh,level,additional) values($1,0,2,'Short Sword','短剑',0,0)",
      [owner],
    );
    assert.deepEqual(
      (await db.query('select item_name_en,item_name_zh from public.item_orders where id=$1', [inserted.rows[0].id])).rows[0],
      { item_name_en: 'Test Sword', item_name_zh: '测试剑' },
    );
    await assert.rejects(db.query(
      "insert into public.item_orders(user_id,item_section,item_index,item_name_en,level,additional) values($1,0,3,'Partial',0,0)",
      [owner],
    ), /not-null constraint/);
    await assert.rejects(db.query(
      "insert into public.item_orders(user_id,item_section,item_index,item_name_en,item_name_zh,level,additional) values($1,0,3,'','中文',0,0)",
      [owner],
    ), /check constraint/);
    async function asUser(uid, role = 'authenticated') {
      const client = new Client({ connectionString: `postgresql://postgres:test-password@127.0.0.1:${port}/postgres` });
      await client.connect(); clients.push(client);
      await client.query(`set role ${role}`);
      await client.query("select set_config('request.jwt.claim.sub',$1,false)", [uid || '']);
      return client;
    }
    const own = await asUser(owner), otherUser = await asUser(other), anon = await asUser(null, 'anon');
    assert.equal((await own.query('select * from public.item_orders')).rows.length, 2);
    assert.equal((await otherUser.query('select * from public.item_orders')).rows.length, 0);
    await assert.rejects(anon.query('select * from public.item_orders'), /permission denied/);
    await assert.rejects(own.query("insert into public.item_orders(user_id,item_section,item_index,item_name_en,item_name_zh,level,additional) values($1,0,2,'Fake','假',0,0)", [owner]), /permission denied/);
    await assert.rejects(own.query('update public.item_orders set status=$1 where id=$2', ['PAID', inserted.rows[0].id]), /permission denied/);
    await assert.rejects(own.query('delete from public.item_orders where id=$1', [inserted.rows[0].id]), /permission denied/);
    await assert.rejects(db.query("insert into public.item_orders(user_id,item_section,item_index,item_name_en,item_name_zh,level,excellent,additional) values($1,0,1,'Test','测试',0,array['invalid'],0)", [owner]), /check constraint/);
  } finally {
    await Promise.all(clients.map((client) => client.end()));
    if (db) await db.end();
    await embedded.stop();
    await rm(directory, { recursive: true, force: true });
  }
}, 60000);
