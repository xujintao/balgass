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
import { test } from 'vitest';
import { profileService, randomNickname } from '../src/lib/profile-service';

test('真实 PostgreSQL schema、Next.js 业务规则及升级', async () => {
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
    await db.query(`do $$ begin if not exists(select 1 from pg_roles where rolname='anon') then create role anon; end if; if not exists(select 1 from pg_roles where rolname='authenticated') then create role authenticated; end if; if not exists(select 1 from pg_roles where rolname='service_role') then create role service_role bypassrls; end if; end $$;
 create schema auth; create table auth.users(id uuid primary key, email_confirmed_at timestamptz);
 create function auth.uid() returns uuid language sql stable as $$ select nullif(current_setting('request.jwt.claim.sub',true),'')::uuid $$;
 grant usage on schema auth to authenticated, anon; grant execute on function auth.uid() to authenticated, anon;`);
    await db.query(
      await readFile(
        new URL('../supabase/schema/000_schema.sql', import.meta.url),
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

    function row(value) {
      if (!value) return null;
      return {
        ...value,
        created_at: value.created_at.toISOString(),
        nickname_changed_at: value.nickname_changed_at?.toISOString() ?? null,
      };
    }
    const repository = {
      async find(uid) {
        return row(
          (await db.query('select * from profiles where user_id=$1', [uid]))
            .rows[0],
        );
      },
      async insert(uid, nickname, key) {
        try {
          return row(
            (
              await db.query(
                'insert into profiles(user_id,nickname,nickname_key) values($1,$2,$3) returning *',
                [uid, nickname, key],
              )
            ).rows[0],
          );
        } catch (e) {
          if (e.code === '23505') return null;
          throw e;
        }
      },
      async update(observed, nickname, key, changedAt, cutoff) {
        const c = new Client({ connectionString: testUrl.toString() });
        await c.connect();
        try {
          return row(
            (
              await c.query(
                'update profiles set nickname=$1,nickname_key=$2,nickname_changed_at=$3 where user_id=$4 and nickname=$5 and nickname_changed_at is not distinct from $6::timestamptz and (nickname_changed_at is null or nickname_changed_at <= $7::timestamptz) returning *',
                [
                  nickname,
                  key,
                  changedAt,
                  observed.user_id,
                  observed.nickname,
                  observed.nickname_changed_at,
                  cutoff,
                ],
              )
            ).rows[0],
          );
        } catch (e) {
          if (e.code === '23505') {
            e.status = 409;
            e.code = 'NICKNAME_TAKEN';
          }
          throw e;
        } finally {
          await c.end();
        }
      },
    };
    const now = new Date('2026-09-28T12:00:00Z');
    const service = profileService(repository, randomNickname, () => now);
    const a = {
      id: await user(),
      email: 'a@example.com',
      email_confirmed_at: now.toISOString(),
    };
    const b = {
      id: await user(),
      email: 'b@example.com',
      email_confirmed_at: now.toISOString(),
    };
    const unverified = {
      id: await user(false),
      email: 'unverified@example.com',
    };
    const ca = await asUser(a.id),
      cu = await asUser(unverified.id),
      anon = await asUser(null, 'anon');
    await check(
      '注册 Auth 用户不再触发创建资料，Next.js 首次创建随机昵称',
      async () => {
        assert.equal((await db.query('select * from profiles')).rows.length, 0);
        const profile = await service.me(a);
        assert.match(profile.nickname, /^玩家_[a-f0-9]{16}$/);
        assert.equal(profile.nicknameChangedAt, null);
        assert.equal((await service.me(a)).nickname, profile.nickname);
      },
    );
    await check('并发首次登录只创建一条资料', async () => {
      const results = await Promise.all([service.me(b), service.me(b)]);
      assert.equal(results[0].nickname, results[1].nickname);
      assert.equal(
        (await db.query('select * from profiles where user_id=$1', [b.id])).rows
          .length,
        1,
      );
    });
    await check('未验证邮箱不能创建资料或修改昵称', async () => {
      await assert.rejects(
        service.me(unverified),
        (e) => e.code === 'EMAIL_UNVERIFIED',
      );
      await assert.rejects(
        service.change(unverified, '未验证用户'),
        (e) => e.code === 'EMAIL_UNVERIFIED',
      );
      assert.equal(await repository.find(unverified.id), null);
    });
    await check('RLS 只允许本人读取，客户端不能直接写表', async () => {
      assert.equal((await ca.query('select * from profiles')).rows.length, 1);
      assert.equal((await cu.query('select * from profiles')).rows.length, 0);
      await assert.rejects(
        anon.query('select * from profiles'),
        /permission denied/,
      );
      for (const sql of [
        "update profiles set nickname='绕过规则'",
        'update profiles set nickname_changed_at=null',
        'delete from profiles',
        `insert into profiles(user_id,nickname,nickname_key) values('${id()}','绕过','绕过')`,
      ])
        await assert.rejects(ca.query(sql), /permission denied/);
    });
    await check('首次立即改名，规范化空白且不误删 v', async () => {
      const p = await service.change(a, '\u3000vav\u000b');
      assert.equal(p.nickname, 'vav');
      assert.equal(p.nicknameChangedAt, now.toISOString());
    });
    await check('同名不消耗次数，大小写调整仍受冷却限制', async () => {
      const before = await service.me(a);
      assert.deepEqual(await service.change(a, 'vav'), before);
      await assert.rejects(
        service.change(a, 'VaV'),
        (e) => e.code === 'NICKNAME_COOLDOWN',
      );
    });
    await check('大小写判重且失败不消耗修改机会', async () => {
      await assert.rejects(
        service.change(b, 'VAV'),
        (e) => e.code === 'NICKNAME_TAKEN',
      );
      assert.equal((await service.me(b)).nicknameChangedAt, null);
    });
    await check('格式校验统一在 TypeScript 中执行', async () => {
      for (const n of [
        'a',
        'a'.repeat(25),
        '名字🐈',
        'a-b',
        'ＡＢ',
        'name space',
      ])
        await assert.rejects(service.change(b, n));
      assert.equal((await service.me(b)).nicknameChangedAt, null);
    });
    await check('未满 30 天拒绝，刚好满 30 天允许，旧昵称释放', async () => {
      await db.query(
        'update profiles set nickname_changed_at=$1 where user_id=$2',
        [new Date(now.getTime() - 30 * 86400000 + 1), a.id],
      );
      await assert.rejects(
        service.change(a, '新昵称'),
        (e) => e.code === 'NICKNAME_COOLDOWN',
      );
      await db.query(
        'update profiles set nickname_changed_at=$1 where user_id=$2',
        [new Date(now.getTime() - 30 * 86400000), a.id],
      );
      await service.change(a, '新昵称');
      await service.change(b, 'vav');
    });
    await check('登录不覆盖用户修改后的昵称或冷却时间', async () => {
      const before = await service.me(a);
      assert.deepEqual(await service.me(a), before);
    });
    await check('同一用户并发改名仅一次成功', async () => {
      await db.query(
        'update profiles set nickname_changed_at=null where user_id=$1',
        [a.id],
      );
      const results = await Promise.allSettled([
        service.change(a, '并发一'),
        service.change(a, '并发二'),
      ]);
      assert.equal(results.filter((r) => r.status === 'fulfilled').length, 1);
      assert.equal(
        results.filter(
          (r) =>
            r.status === 'rejected' && r.reason.code === 'NICKNAME_COOLDOWN',
        ).length,
        1,
      );
    });
    await check('两个用户争抢昵称仅一个成功', async () => {
      await db.query(
        'update profiles set nickname_changed_at=null where user_id=any($1)',
        [[a.id, b.id]],
      );
      const results = await Promise.allSettled([
        service.change(a, '唯一共享名'),
        service.change(b, '唯一共享名'),
      ]);
      assert.equal(results.filter((r) => r.status === 'fulfilled').length, 1);
      assert.equal(
        results.filter(
          (r) => r.status === 'rejected' && r.reason.code === 'NICKNAME_TAKEN',
        ).length,
        1,
      );
    });
    await check('随机昵称碰撞后在 Next.js 中重新生成', async () => {
      const c = { ...a, id: await user() };
      let attempt = 0;
      const existing = (await service.me(a)).nickname;
      const collision = profileService(
        repository,
        () => (++attempt === 1 ? existing : '随机新名字'),
        () => now,
      );
      assert.equal((await collision.me(c)).nickname, '随机新名字');
      assert.equal(attempt, 2);
    });
    await check('随机昵称持续碰撞时有限重试', async () => {
      const c = { ...a, id: await user() };
      let attempts = 0;
      const existing = (await service.me(a)).nickname;
      const collision = profileService(
        repository,
        () => {
          attempts++;
          return existing;
        },
        () => now,
      );
      await assert.rejects(
        collision.me(c),
        (e) => e.code === 'PROFILE_UNAVAILABLE',
      );
      assert.equal(attempts, 10);
    });
    await check('升级旧表保留数据并删除旧函数和触发器', async () => {
      await db.query('drop table profiles');
      await db.query(
        await readFile(
          new URL('./fixtures/legacy-profiles.sql', import.meta.url),
          'utf8',
        ),
      );
      const uid = await user();
      await db.query(
        "update profiles set nickname='旧版昵称',nickname_changed_at=$1 where user_id=$2",
        [now, uid],
      );
      await db.query(
        await readFile(
          new URL(
            '../supabase/upgrade/001_nextjs_profiles.sql',
            import.meta.url,
          ),
          'utf8',
        ),
      );
      const saved = await repository.find(uid);
      assert.equal(saved.nickname, '旧版昵称');
      assert.equal(saved.nickname_changed_at, now.toISOString());
      assert.equal(
        (
          await db.query(
            "select * from pg_proc p join pg_namespace n on n.oid=p.pronamespace where n.nspname='public' and p.proname in ('change_nickname','create_profile','has_verified_email','normalize_nickname','valid_nickname')",
          )
        ).rows.length,
        0,
      );
      const fresh = { ...a, id: await user() };
      assert.equal(await repository.find(fresh.id), null);
      await service.me(fresh);
      await service.change(fresh, '升级后昵称');
      const legacy = { ...a, id: uid };
      await assert.rejects(
        service.change(legacy, '不应绕过冷却'),
        (e) => e.code === 'NICKNAME_COOLDOWN',
      );
      await db.query(
        await readFile(
          new URL('../supabase/schema/000_schema.sql', import.meta.url),
          'utf8',
        ),
      );
      assert.equal((await repository.find(uid)).nickname, '旧版昵称');
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
}, 60000);
