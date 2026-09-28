// Run exact production migration against an isolated real PostgreSQL 17 database.
// Minimal auth schema emulates Supabase's auth.uid(), not its authentication server.
import assert from 'node:assert/strict';
import { readFile, mkdtemp, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { createServer } from 'node:net';
import { randomUUID } from 'node:crypto';
import { Client } from 'pg';
import EmbeddedPostgres from 'embedded-postgres';
const socket = createServer();
await new Promise((resolve) => socket.listen(0, '127.0.0.1', resolve));
const port = socket.address().port;
await new Promise((resolve) => socket.close(resolve));
const directory = await mkdtemp(join(tmpdir(), 'balgass-pg-'));
const embedded = process.env.TEST_DATABASE_URL
  ? null
  : new EmbeddedPostgres({
      databaseDir: join(directory, 'data'),
      user: 'postgres',
      password: 'test-password',
      port,
      persistent: false,
      onLog: () => {},
      onError: () => {},
    });
let admin, db;
const clients = [];
const dbName = `test_${randomUUID().replaceAll('-', '')}`;
let passed = 0;
async function check(name, test) {
  await test();
  passed++;
  console.log(`✓ ${name}`);
}
const id = () => randomUUID();
try {
  if (embedded) {
    await embedded.initialise();
    await embedded.start();
  }
  const url =
    process.env.TEST_DATABASE_URL ||
    `postgresql://postgres:test-password@127.0.0.1:${port}/postgres`;
  admin = new Client({ connectionString: url });
  await admin.connect();
  await admin.query(`create database ${dbName}`);
  const testUrl = new URL(url);
  testUrl.pathname = `/${dbName}`;
  db = new Client({ connectionString: testUrl.toString() });
  await db.connect();
  await db.query(`do $$ begin if not exists(select 1 from pg_roles where rolname='anon') then create role anon; end if; if not exists(select 1 from pg_roles where rolname='authenticated') then create role authenticated; end if; end $$;
 create schema auth; create table auth.users(id uuid primary key, email_confirmed_at timestamptz);
 create function auth.uid() returns uuid language sql stable as $$ select nullif(current_setting('request.jwt.claim.sub',true),'')::uuid $$;
 grant usage on schema auth to authenticated, anon; grant execute on function auth.uid() to authenticated, anon;`);
  await db.query(
    await readFile(
      new URL(
        '../supabase/migrations/202609280001_profiles.sql',
        import.meta.url,
      ),
      'utf8',
    ),
  );
  async function user(verified = true) {
    const uid = id();
    await db.query(
      'insert into auth.users(id,email_confirmed_at) values($1,$2)',
      [uid, verified ? new Date() : null],
    );
    return uid;
  }
  async function asUser(uid, role = 'authenticated') {
    const c = new Client({ connectionString: testUrl.toString() });
    await c.connect();
    clients.push(c);
    await c.query(`set role ${role}`);
    await c.query("select set_config('request.jwt.claim.sub',$1,false)", [
      uid || '',
    ]);
    return c;
  }
  const a = await user(),
    b = await user(),
    unverified = await user(false);
  const ca = await asUser(a),
    cb = await asUser(b),
    cu = await asUser(unverified),
    anon = await asUser(null, 'anon');
  const rename = (c, n) =>
    c.query('select (public.change_nickname($1)).*', [n]);
  await check('每个用户得到随机且唯一的昵称', async () => {
    const { rows } = await db.query('select * from profiles');
    assert.equal(rows.length, 3);
    assert.equal(new Set(rows.map((r) => r.nickname)).size, 3);
    assert(rows.every((r) => /^玩家_[a-f0-9]{16}$/.test(r.nickname)));
    assert(rows.every((r) => r.nickname_changed_at === null));
  });
  await check('RLS 只允许本人读取，未验证邮箱不可读取', async () => {
    assert.equal((await ca.query('select * from profiles')).rows.length, 1);
    assert.equal((await ca.query('select * from profiles')).rows[0].user_id, a);
    assert.equal((await cu.query('select * from profiles')).rows.length, 0);
    await assert.rejects(
      anon.query('select * from profiles'),
      /permission denied/,
    );
  });
  await check('无法直接修改昵称、用户关联和时间或新增删除资料', async () => {
    for (const query of [
      "update profiles set nickname='绕过规则'",
      'update profiles set nickname_changed_at=null',
      `update profiles set user_id='${b}'`,
      'delete from profiles',
      `insert into profiles(user_id,nickname) values('${id()}','绕过插入')`,
    ])
      await assert.rejects(ca.query(query), /permission denied/);
  });
  await check('空白规范化不误删英文字母 v', async () => {
    assert.equal(
      (await db.query("select public.normalize_nickname('vav') as nickname"))
        .rows[0].nickname,
      'vav',
    );
    assert.equal(
      (
        await db.query(
          "select public.normalize_nickname(chr(11)||'vav'||chr(11)) as nickname",
        )
      ).rows[0].nickname,
      'vav',
    );
  });
  await check('首次可立即修改，统一处理 Unicode 空白', async () => {
    const { rows } = await rename(ca, '\u3000玩家_Ab\u00a0');
    assert.equal(rows[0].nickname, '玩家_Ab');
    assert(rows[0].nickname_changed_at);
  });
  await check('相同昵称不消耗次数，大小写变更受到冷却限制', async () => {
    const previous = (
      await db.query(
        'select nickname_changed_at from profiles where user_id=$1',
        [a],
      )
    ).rows[0].nickname_changed_at;
    await rename(ca, '玩家_Ab');
    assert.equal(
      (
        await db.query(
          'select nickname_changed_at from profiles where user_id=$1',
          [a],
        )
      ).rows[0].nickname_changed_at.getTime(),
      previous.getTime(),
    );
    await assert.rejects(rename(ca, '玩家_ab'), /NICKNAME_COOLDOWN/);
  });
  await check('大小写判重，失败不消耗第一次修改机会', async () => {
    await assert.rejects(rename(cb, '玩家_aB'), /NICKNAME_TAKEN/);
    assert.equal(
      (
        await db.query(
          'select nickname_changed_at from profiles where user_id=$1',
          [b],
        )
      ).rows[0].nickname_changed_at,
      null,
    );
    await rename(cb, '另一位玩家');
  });
  await check('非法昵称无法通过 RPC 写入', async () => {
    for (const n of [
      '',
      'a',
      'a'.repeat(25),
      '名字🐈',
      'name space',
      'a-b',
      'ＡＢ',
      'éé',
    ])
      await assert.rejects(rename(cb, n), /INVALID_NICKNAME/);
  });
  await check('未验证与匿名用户不能修改昵称', async () => {
    await assert.rejects(rename(cu, '未验证用户'), /UNAUTHENTICATED/);
    await assert.rejects(rename(anon, '匿名用户'), /permission denied/);
  });
  await check('30 天内拒绝修改并返回明确 UTC 时间', async () => {
    await db.query(
      "update profiles set nickname_changed_at=clock_timestamp()-interval '719 hours' where user_id=$1",
      [a],
    );
    await assert.rejects(
      rename(ca, '未满30天'),
      (e) =>
        e.message === 'NICKNAME_COOLDOWN' &&
        /^\d{4}-\d{2}-\d{2}T.*Z$/.test(e.detail),
    );
  });
  await check('满 30 天可再次修改，旧昵称被释放', async () => {
    await db.query(
      "update profiles set nickname_changed_at=clock_timestamp()-interval '720 hours' where user_id=$1",
      [a],
    );
    await rename(ca, '新的昵称');
    await db.query(
      'update profiles set nickname_changed_at=null where user_id=$1',
      [b],
    );
    await rename(cb, '玩家_Ab');
  });
  await check('并发修改同一用户仅一项成功', async () => {
    await db.query(
      'update profiles set nickname_changed_at=null where user_id=$1',
      [a],
    );
    const ca2 = await asUser(a);
    const results = await Promise.allSettled([
      rename(ca, '并发昵称1'),
      rename(ca2, '并发昵称2'),
    ]);
    assert.equal(results.filter((r) => r.status === 'fulfilled').length, 1);
    assert.equal(
      results.filter(
        (r) =>
          r.status === 'rejected' && r.reason.message === 'NICKNAME_COOLDOWN',
      ).length,
      1,
    );
  });
  await check('两个用户争抢同一昵称仅一项成功', async () => {
    await db.query(
      'update profiles set nickname_changed_at=null where user_id=any($1)',
      [[a, b]],
    );
    const results = await Promise.allSettled([
      rename(ca, '共享唯一名'),
      rename(cb, '共享唯一名'),
    ]);
    assert.equal(results.filter((r) => r.status === 'fulfilled').length, 1);
    assert.equal(
      results.filter(
        (r) => r.status === 'rejected' && r.reason.message === 'NICKNAME_TAKEN',
      ).length,
      1,
    );
  });
  await check('汉字 Unicode 范围与 Web 验证一致', async () => {
    const { rows } = await db.query(
      "select cp, public.valid_nickname(chr(cp)||'a') as valid from generate_series(1,1114111) cp where cp not between 55296 and 57343",
    );
    for (const r of rows) {
      const expected = /^[\p{Script=Han}A-Za-z0-9_]$/u.test(
        String.fromCodePoint(r.cp),
      );
      assert.equal(r.valid, expected, `U+${r.cp.toString(16)}`);
    }
  });
  await check('随机昵称碰撞时重试', async () => {
    await db.query('begin');
    try {
      await db.query(
        "create sequence extensions.test_random_counter; create or replace function extensions.gen_random_bytes(integer) returns bytea language plpgsql as $$ begin if nextval('extensions.test_random_counter') <= 2 then return decode(repeat('00',$1),'hex'); end if; return decode(repeat('01',$1),'hex'); end $$;",
      );
      const first = await user(),
        second = await user();
      const { rows } = await db.query(
        'select nickname from profiles where user_id=any($1)',
        [[first, second]],
      );
      assert.equal(new Set(rows.map((r) => r.nickname)).size, 2);
      assert.equal(
        (
          await db.query(
            'select last_value from extensions.test_random_counter',
          )
        ).rows[0].last_value,
        '3',
      );
    } finally {
      await db.query('rollback');
    }
  });
  console.log(
    `${passed} database checks passed (PostgreSQL ${(await db.query('show server_version')).rows[0].server_version}).`,
  );
} finally {
  await Promise.all(clients.map((c) => c.end()));
  if (db) await db.end();
  if (admin) {
    await admin.query(`drop database if exists ${dbName} with (force)`);
    await admin.end();
  }
  if (embedded) await embedded.stop();
  await rm(directory, { recursive: true, force: true });
}
