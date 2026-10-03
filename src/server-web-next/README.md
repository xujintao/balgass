# r2f2 账号网站

Next.js App Router + Supabase Auth/PostgreSQL。邮箱为登录标识，默认 Passkey 登录，备用邮件验证码；昵称由 Next.js 服务端随机生成，唯一且每 30 天可修改一次。

## 本地启动

需要 Node.js >=22.18 和 npm。提交的 `package-lock.json` 用于复现依赖。

```sh
npm ci
cp .env.example .env.local
npm run dev
```

本地 Next.js 直接连接云端开发 Supabase，不需要 Supabase CLI 或 Docker。通过 Dashboard SQL Editor 执行 `supabase/migrations/20260928000100_create_profiles.sql`，在 `.env.local` 配置云端 URL、publishable key 和**服务端 secret key**。浏览器仍访问 http://localhost:3000。

首次邮箱验证成功后，Next.js 创建 `profiles` 并随机生成昵称；以后邮件登录、passkey 登录及读取本人资料时也会补齐缺失记录。已有昵称和修改时间保持不变，不在每次登录时重置。生成规则集中在 `src/lib/profile-service.ts` 的 `randomNickname()`，当前为「玩家_」加 16 位随机十六进制字符。

`.env.example` 中的 CAPTCHA 关闭开关只适用于 localhost/loopback。正式环境需要 Turnstile；客户端只接收公开 site key。没有认证配置时 API 返回 `503 CONFIGURATION_REQUIRED`。

## 数据库迁移

所有建表和后续修改 SQL 统一保存在 `supabase/migrations`，使用 `YYYYMMDDHHMMSS_描述.sql` 命名。在云端 SQL Editor 中按文件名顺序执行，不需要 Supabase CLI。当前只有新项目的初始迁移，不包含旧版数据库升级脚本。

`public.migrations` 保存迁移 ID、文件名和成功执行时间；它与业务表 `public.profiles` 分开，不允许浏览器、App 或服务端认证客户端读写。SQL Editor 的数据库管理员负责执行迁移。每份迁移在事务内检查执行记录、执行变更并记录成功；重复执行跳过，失败全部回滚。一次性 `DO` 块不会创建存储函数或触发器。初始迁移要求数据库尚未创建旧版 profiles 表。

已执行的迁移文件不再修改。后续变更新建时间戳文件，沿用事务、迁移记录检查和写入方式，并在检查当前记录前验证前一份迁移已经成功；不要手动删除迁移记录来强制重跑。这里的 `public.migrations` 是手动执行历史，不是 Supabase CLI 的内部迁移记录表。

## 生产配置

1. 建立 Supabase 项目，新建项目执行 `supabase/migrations/20260928000100_create_profiles.sql`。昵称生成、校验和修改规则由 Next.js 服务层统一执行；数据库只保留表、唯一约束、RLS 和权限。认证用户没有 profiles 写权限，只有服务端管理客户端可以写入。不要给客户端 service-role/secret key。
2. 配置 Email Auth：允许注册、要求邮箱确认，OTP 长度 6、有效期 600 秒、发送间隔 60 秒。将 **Confirm signup** 和 **Magic Link** 两个模板设置为 `supabase/templates/code.html` 内容，使用 `.Token`，不发送登录链接。
3. 配置自己的 SMTP（发件域名、发件人、主机、端口、用户名、密码），确认供应商域名验证及 SPF/DKIM。SMTP 密码只放在 Supabase 配置中。
4. 启用 Supabase Auth CAPTCHA，选择 Turnstile 并配置其 secret；在 Turnstile 中允许正式域名。Vercel 配置 `NEXT_PUBLIC_TURNSTILE_SITE_KEY`，删除 `AUTH_CAPTCHA_DISABLED`。
5. 配置 Supabase Auth rate limits：邮件 30/小时、登录/注册 30/5 分钟、验证码校验 30/5 分钟、刷新 150/5 分钟作为初始值，再根据实际流量调整。Supabase 对 API 来源地址实施的限流可能聚合 Vercel 出口流量；邮件发送间隔和 CAPTCHA 仍由 Auth 强制执行。
6. 启用 passkey：RP display name 设为 `r2f2`，RP ID 为稳定正式域名（无协议、端口和路径），origins 为确切 HTTPS 网站地址。RP ID 变更会使旧 passkey 失效。开发项目的 RP ID 设为 `localhost`，origin 设为 `http://localhost:3000`；部署在 `next.r2f2.com` 的项目应核对其 RP ID 和 `https://next.r2f2.com` origin，切换到 `r2f2.com` 前也应核对正式站点 origin。品牌改名不要求修改现有 RP ID。开发、预览和生产应使用独立 Supabase 项目；Vercel 临时预览域名使用邮件登录，不加入生产 RP 配置。
7. Vercel Root Directory 设为 `src/server-web-next`，Node >=22.18，构建 `npm run build`。配置服务端 `SUPABASE_URL`、`SUPABASE_PUBLISHABLE_KEY`、`SUPABASE_SECRET_KEY`、`APP_ORIGIN`（无尾斜杠）和公开 Turnstile site key。Supabase Site URL 设为正式站点。正式环境 HTTPS 下 Cookie 自动启用 Secure。

Passkey 是 Supabase 的实验接口；`@supabase/supabase-js` 固定为 2.105.0，实验设置显式开启。登录使用 SDK；受认证的注册/管理 REST 请求封装于 `src/lib/passkeys.ts`，直接携带用户 access token，使 App 不需要为这些操作交出 refresh token。升级 SDK 前运行类型检查和真实认证测试。

## API 与会话

OpenAPI 文档：[/openapi.json](./public/openapi.json)。所有 JSON API 位于 `/api/v1/`；响应为 `{data: ...}` 或 `{error: {code,message,...}}`。

浏览器会话使用 `r2f2_access`、`r2f2_refresh` 两个 HttpOnly、SameSite=Lax Cookie。浏览器写请求必须携带与 `APP_ORIGIN` 完全一致的 Origin。浏览器认证响应不返回令牌；私有页面动态渲染，API 禁止缓存。访问令牌失效时页面或 API 客户端通过刷新端点轮换会话。会话 Cookie 保留 30 天；Auth 可以更早撤销或使会话失效。

App 所有请求设置 `X-Client-Type: app`，私有请求另携带 `Authorization: Bearer <accessToken>`。App 请求不使用 Cookie。验证码/passkey 验证成功返回 accessToken、refreshToken、expiresAt；刷新发送 `{refreshToken}`，不要求尚未过期的 access token。原生 WebAuthn 在 App 侧完成，options 中二进制字段使用 base64url，verify 提交 `challengeId` 和序列化 credential。将 refresh token 存储于设备安全存储；未来原生 App 还需要单独配置系统域名关联和允许的原生 origin。

退出撤销当前 refresh 会话并清除 Cookie；Supabase 已签发的 access JWT 在过期前仍可能有效，这是 Auth 的令牌语义。不要将退出理解为立即吊销所有已签发的 access JWT。

昵称：NFC、去首尾空白、2–24 个 Unicode 字符，仅 Han 字符、ASCII 英文字母/数字和下划线；英文大小写不影响唯一性，大小写调整仍算修改。首次修改立即可用；后续间隔 720 小时（30 天），以 Next.js 服务端 UTC 时间为准。修改操作在同一个 UPDATE 中匹配旧昵称、旧修改时间和允许修改的时间范围，防止并发请求覆盖。相同昵称不消耗次数；失败不更新冷却时间；旧昵称释放后可被其他用户使用。邮箱修改、旧账号迁移和公开个人资料不在当前范围内。

## 验证

```sh
npm run typecheck
npm run lint
npm test
npm run test:db
npm run build
npx playwright install chromium --no-shell
npm run test:e2e
```

`test:db` 默认自动启动独立 PostgreSQL 17，在临时数据库执行实际 schema，并让 Next.js 业务服务通过测试 repository 操作真实数据库，验证唯一性、RLS、昵称初始化、并发修改和冷却时间。同时验证迁移失败整体回滚、重复执行跳过且保留数据，以及迁移记录的访问权限。它不模拟真实 Supabase Auth 服务。也可设置 `TEST_DATABASE_URL` 连接有 CREATE DATABASE 权限的测试 PostgreSQL；测试创建并删除随机名称的临时数据库，不操作已有业务表。

资料写入使用服务端 secret key，会绕过 RLS；因此权限和规则必须由共享服务执行。浏览器和 App 不接收此密钥，不允许直接修改资料表。该设计将业务规则保留在 TypeScript 中，未来更换数据库时可复用服务层，只需更换 repository。

浏览器测试包含无需真实服务的页面与虚拟认证器交互测试；真实 Supabase 测试需要 `E2E_LIVE_AUTH=true`、`SUPABASE_URL`、`SUPABASE_PUBLISHABLE_KEY`、`TEST_SUPABASE_SECRET_KEY` 及对应 localhost passkey 配置。测试密钥用于创建/删除专用测试用户、获取验证码，以及作为测试网站的服务端资料写入密钥，不配置为 NEXT_PUBLIC 环境变量，不在生产项目运行测试。真实测试覆盖虚拟认证器注册、登录、删除、challenge 重放和跨账号操作；没有这些配置时明确跳过该组测试。

发布前还应真实发送邮件并在支持 passkey 的手机/桌面设备上检查注册、恢复访问和删除凭据。自动测试不能代替 SMTP 投递及设备兼容性联调。
