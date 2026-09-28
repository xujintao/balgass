# Balgass 账号网站

Next.js App Router + Supabase Auth/PostgreSQL。邮箱为登录标识，默认 Passkey 登录，备用邮件验证码；昵称由数据库随机生成，唯一且每 30 天可修改一次。

## 本地启动

需要 Node.js >=22.18 和 npm。提交的 `package-lock.json` 用于复现依赖。

```sh
npm ci
cp .env.example .env.local
npm run dev
```

完整本地认证还需要 Docker 与支持 passkey 的近期 Supabase CLI：运行 `supabase start`、`supabase db reset`，将 CLI 输出的 URL 和 anon/publishable key 写入 `.env.local`。本地邮箱在 http://localhost:54324 查看。打开 http://localhost:3000，使用同一个 hostname；`127.0.0.1` 不是这里配置的 passkey RP ID。

`.env.example` 中的 CAPTCHA 关闭开关只适用于 localhost/loopback。正式环境需要 Turnstile；客户端只接收公开 site key。没有认证配置时 API 返回 `503 CONFIGURATION_REQUIRED`。

## 生产配置

1. 建立 Supabase 项目，执行 `supabase/migrations/202609280001_profiles.sql`。数据库函数和 RLS 共同保护昵称；认证用户没有 profiles 写权限，修改通过 `change_nickname` RPC 执行。不要给客户端 service-role/secret key。
2. 配置 Email Auth：允许注册、要求邮箱确认，OTP 长度 6、有效期 600 秒、发送间隔 60 秒。将 **Confirm signup** 和 **Magic Link** 两个模板设置为 `supabase/templates/code.html` 内容，使用 `.Token`，不发送登录链接。
3. 配置自己的 SMTP（发件域名、发件人、主机、端口、用户名、密码），确认供应商域名验证及 SPF/DKIM。SMTP 密码只放在 Supabase 配置中。
4. 启用 Supabase Auth CAPTCHA，选择 Turnstile 并配置其 secret；在 Turnstile 中允许正式域名。Vercel 配置 `NEXT_PUBLIC_TURNSTILE_SITE_KEY`，删除 `AUTH_CAPTCHA_DISABLED`。
5. 配置 Supabase Auth rate limits：邮件 30/小时、登录/注册 30/5 分钟、验证码校验 30/5 分钟、刷新 150/5 分钟作为初始值，再根据实际流量调整。Supabase 对 API 来源地址实施的限流可能聚合 Vercel 出口流量；邮件发送间隔和 CAPTCHA 仍由 Auth 强制执行。
6. 启用 passkey：RP display name `Balgass`，RP ID 为稳定正式域名（无协议、端口和路径），origins 为确切 HTTPS 网站地址。RP ID 变更会使旧 passkey 失效。开发、预览和生产应使用独立 Supabase 项目；Vercel 临时预览域名使用邮件登录，不加入生产 RP 配置。
7. Vercel Root Directory 设为 `src/server-web-next`，Node >=22.18，构建 `npm run build`。配置服务端 `SUPABASE_URL`、`SUPABASE_PUBLISHABLE_KEY`、`APP_ORIGIN`（无尾斜杠）和公开 Turnstile site key。Supabase Site URL 设为正式站点。正式环境 HTTPS 下 Cookie 自动启用 Secure。

Passkey 是 Supabase 的实验接口；`@supabase/supabase-js` 固定为 2.105.0，实验设置显式开启。登录使用 SDK；受认证的注册/管理 REST 请求封装于 `src/lib/passkeys.ts`，直接携带用户 access token，使 App 不需要为这些操作交出 refresh token。升级 SDK 前运行类型检查和真实认证测试。

## API 与会话

OpenAPI 文档：[/openapi.json](./public/openapi.json)。所有 JSON API 位于 `/api/v1/`；响应为 `{data: ...}` 或 `{error: {code,message,...}}`。

浏览器会话使用 HttpOnly、SameSite=Lax Cookie。浏览器写请求必须携带与 `APP_ORIGIN` 完全一致的 Origin。浏览器认证响应不返回令牌；私有页面动态渲染，API 禁止缓存。访问令牌失效时页面或 API 客户端通过刷新端点轮换会话。会话 Cookie 保留 30 天；Auth 可以更早撤销或使会话失效。

App 所有请求设置 `X-Client-Type: app`，私有请求另携带 `Authorization: Bearer <accessToken>`。App 请求不使用 Cookie。验证码/passkey 验证成功返回 accessToken、refreshToken、expiresAt；刷新发送 `{refreshToken}`，不要求尚未过期的 access token。原生 WebAuthn 在 App 侧完成，options 中二进制字段使用 base64url，verify 提交 `challengeId` 和序列化 credential。将 refresh token 存储于设备安全存储；未来原生 App 还需要单独配置系统域名关联和允许的原生 origin。

退出撤销当前 refresh 会话并清除 Cookie；Supabase 已签发的 access JWT 在过期前仍可能有效，这是 Auth 的令牌语义。不要将退出理解为立即吊销所有已签发的 access JWT。

昵称：NFC、去首尾空白、2–24 个 Unicode 字符，仅 Han 字符、ASCII 英文字母/数字和下划线；英文大小写不影响唯一性，大小写调整仍算修改。首次修改立即可用；后续间隔 720 小时（30 天），以数据库时间为准。相同昵称不消耗次数；失败不更新冷却时间；旧昵称释放后可被其他用户使用。邮箱修改、旧账号迁移和公开个人资料不在当前范围内。

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

`test:db` 默认自动启动独立 PostgreSQL 17，在临时数据库中执行实际迁移，模拟 Supabase 的 auth schema 与 JWT 身份，检查 RLS、唯一性、并发、时间边界和 Unicode 一致性。它不模拟真实 Auth 服务。也可通过 `TEST_DATABASE_URL` 连接有 CREATE DATABASE 权限的测试 PostgreSQL；测试创建并删除随机名称的临时数据库，不操作已有业务表。

浏览器测试包含无需真实服务的页面与虚拟认证器交互测试；真实 Supabase 测试需要 `E2E_LIVE_AUTH=true`、`SUPABASE_URL`、`SUPABASE_PUBLISHABLE_KEY`、`TEST_SUPABASE_SECRET_KEY` 及对应 localhost passkey 配置。测试密钥只用于创建/删除专用测试用户及获取验证码，不配置为 NEXT_PUBLIC 环境变量，不在生产项目运行测试。真实测试覆盖虚拟认证器注册、登录、删除、challenge 重放和跨账号操作；没有这些配置时明确跳过该组测试。

发布前还应真实发送邮件并在支持 passkey 的手机/桌面设备上检查注册、恢复访问和删除凭据。自动测试不能代替 SMTP 投递及设备兼容性联调。
