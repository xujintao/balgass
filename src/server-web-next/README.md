# r2f2 网站

Next.js App Router + Supabase Auth/PostgreSQL。网站提供账号管理、游戏实时地图、道具目录和待处理订单。邮箱为登录标识，默认 Passkey 登录，备用邮件验证码；昵称由 Next.js 服务端随机生成，唯一且每 30 天可修改一次。

## 本地启动

需要 Node.js >=22.18 和 npm。提交的 `package-lock.json` 用于复现依赖。

```sh
npm ci
cp .env.example .env.local
npm run dev
```

本地 Next.js 直接连接云端开发 Supabase，不需要 Supabase CLI 或 Docker。通过 Dashboard SQL Editor 按文件名顺序执行 `supabase/migrations/` 中的 SQL，在 `.env.local` 配置云端 URL、publishable key 和**服务端 secret key**。浏览器仍访问 http://localhost:3000。

首次邮箱验证成功后，Next.js 创建 `profiles` 并随机生成昵称；以后邮件登录、passkey 登录及读取本人资料时也会补齐缺失记录。已有昵称和修改时间保持不变，不在每次登录时重置。生成规则集中在 `src/lib/profile-service.ts` 的 `randomNickname()`，当前为 `player_` 加 16 位随机十六进制字符。自选昵称允许各种语言的字母及组合标记、数字 0–9 和下划线，长度为 2–24 个可见字符。

本地和 Vercel 共用启用 CAPTCHA 的云端 Supabase 项目，两处都需要 Turnstile。将对应小组件的公开 site key 配置为 `NEXT_PUBLIC_TURNSTILE_SITE_KEY`，并在 Cloudflare 中允许 `localhost` 和正式网站域名 `r2f2.com`；Supabase Auth 配置同一小组件的 secret。没有 token 时认证请求返回 `400 CAPTCHA_REQUIRED`。没有认证配置时 API 返回 `503 CONFIGURATION_REQUIRED`。

## 数据库迁移

所有建表和后续修改 SQL 统一保存在 `supabase/migrations`，使用 `YYYYMMDDHHMMSS_描述.sql` 命名。在云端 SQL Editor 中按文件名顺序执行，不需要 Supabase CLI。这些迁移不包含旧版数据库升级脚本。

`public.migrations` 保存迁移 ID、文件名和成功执行时间；它与业务表 `public.profiles` 分开，不允许浏览器、App 或服务端认证客户端读写。SQL Editor 的数据库管理员负责执行迁移。每份迁移在事务内检查执行记录、执行变更并记录成功；重复执行跳过，失败全部回滚。一次性 `DO` 块不会创建存储函数或触发器。初始迁移要求数据库尚未创建旧版 profiles 表。

已执行的迁移文件不再修改。后续变更新建时间戳文件，沿用事务、迁移记录检查和写入方式，并在检查当前记录前验证前一份迁移已经成功；不要手动删除迁移记录来强制重跑。这里的 `public.migrations` 是手动执行历史，不是 Supabase CLI 的内部迁移记录表。

## 生产配置

1. 建立 Supabase 项目，按顺序执行 `supabase/migrations/` 中的 SQL。昵称和订单的输入、归属规则由 Next.js 服务层校验；数据库保留约束、RLS 和权限。认证用户只能读取自己的资料与订单；只有服务端管理客户端可以写入。不要给客户端 service-role/secret key。
2. 配置 Email Auth：允许注册、要求邮箱确认，OTP 长度 6、有效期 600 秒、发送间隔 60 秒。将 **Confirm signup** 和 **Magic Link** 两个模板设置为 `supabase/templates/code.html` 内容，使用 `.Token`，不发送登录链接。
3. 配置自己的 SMTP（发件域名、发件人、主机、端口、用户名、密码），确认供应商域名验证及 SPF/DKIM。SMTP 密码只放在 Supabase 配置中。
4. 启用 Supabase Auth CAPTCHA，选择 Turnstile 并配置其 secret；在 Turnstile 中允许 `localhost` 和 `r2f2.com`。本地及 Vercel 均配置对应的 `NEXT_PUBLIC_TURNSTILE_SITE_KEY`。
5. 配置 Supabase Auth rate limits：邮件 30/小时、登录/注册 30/5 分钟、验证码校验 30/5 分钟、刷新 150/5 分钟作为初始值，再根据实际流量调整。Supabase 对 API 来源地址实施的限流可能聚合 Vercel 出口流量；邮件发送间隔和 CAPTCHA 仍由 Auth 强制执行。
6. 启用 passkey：RP display name 设为 `r2f2`，生产项目核对 RP ID 为 `r2f2.com`（无协议、端口和路径）、origin 为 `https://r2f2.com`。RP ID 变更会使旧 passkey 失效。开发项目的 RP ID 设为 `localhost`，origin 设为 `http://localhost:3000`。开发、预览和生产应使用独立 Supabase 项目；Vercel 临时预览域名使用邮件登录，不加入生产 RP 配置。
7. Vercel Root Directory 设为 `src/server-web-next`，Node >=22.18，构建 `npm run build`。配置服务端 `SUPABASE_URL`、`SUPABASE_PUBLISHABLE_KEY`、`SUPABASE_SECRET_KEY`、`APP_ORIGIN=https://r2f2.com`（无尾斜杠）和公开 Turnstile site key。Supabase Site URL 设为 `https://r2f2.com`。正式环境 HTTPS 下 Cookie 自动启用 Secure。

Passkey 是 Supabase 的实验接口；`@supabase/supabase-js` 固定为 2.105.0，实验设置显式开启。登录使用 SDK；受认证的注册/管理 REST 请求封装于 `src/lib/passkeys.ts`，直接携带用户 access token，使 App 不需要为这些操作交出 refresh token。升级 SDK 前运行类型检查和真实认证测试。

## API 与会话

游戏账号列表与创建由 Next.js 服务层通过游戏服务的 `/api/command` 执行。开发环境可将 `GAME_API_URL` 指向本机 `http://localhost:8080/api/command`；生产环境须使用 `https://game.r2f2.com/api/command`。两端配置相同、至少 32 字符的随机 `GAME_API_TOKEN`，例如用 `openssl rand -hex 32` 生成。此密钥授权 `/api/command` 的全部命令，包括删除账号及 bot 管理；Next.js 玩家功能目前只调用 `GetAccountList` 和 `CreateAccount`。密钥仅存于服务端环境，不提交到版本库。旧 `/api/accounts` 和 `/api/bots` 接口已移除。

部署时由 `game.r2f2.com` 的 HTTPS 入口代理 `/api/command` 与公开的 `/api/game` WebSocket；现有 VPS 防火墙规则继续阻止公网直连 8080。先部署并配置游戏服务，再配置网站环境变量；旧 Django 账号接口不参与新网站调用。

公开的 `/game` 与 `/items` 页面使用三个独立地址：`GAME_API_URL` 仅供 Next.js 服务端执行命令，`GAME_WEBSOCKET_URL` 传给浏览器连接实时地图，`GAME_CONFIG_URL` 仅供 Next.js 服务端读取配置目录。本地在 `.env.local` 设置 `GAME_CONFIG_URL=file:///home/pi/balgass/config/server-game-common/IGCData/`，直接读取 XML；生产设置 `GAME_CONFIG_URL=https://game.r2f2.com/config/`、`GAME_WEBSOCKET_URL=wss://game.r2f2.com/api/game`。三个变量均需显式配置，目录 URL 必须以 `/` 结尾。地图与道具配置的 HTTPS 解析结果缓存 60 秒；本地文件不缓存。

Caddy 容器将 VPS 的 `IGCData` 目录只读挂载到 `/srv/game-config`，例如使用 `-v ~/balgass/config/server-game-common/IGCData:/srv/game-config:ro`。在现有 `game.r2f2.com` 站点中增加以下路由；`/config/` 只放行列出的文件，不要对整个目录启用 `file_server`：

```caddyfile
@gameConfigFiles {
    path /config/IGC_MapList.xml
    path /config/Skills/IGC_SkillList.xml
    path /config/Items/IGC_ItemList.xml
    path /config/Items/IGC_ItemSetType.xml
    path /config/Items/IGC_ItemSetOption.xml
}
handle @gameConfigFiles {
    basic_auth {
        nextjs <caddy hash-password 生成的哈希>
    }
    root * /srv/game-config
    uri strip_prefix /config
    file_server
}
handle /api/game {
    reverse_proxy server-game:8080
}
handle /api/command {
    reverse_proxy server-game:8080
}
```

将示例中的 `server-game:8080` 换为 Caddy 容器实际可访问的游戏服务地址。现有 `docker/caddy-fail2ban/Caddyfile` 已预留四条注释的道具/技能路径；商城上线前应启用这四条精确路径并重新加载 Caddy，其他 `/config/` 路径仍返回 404。在 Vercel 设置仅服务端使用的 `GAME_CONFIG_BASIC_AUTH_USER=nextjs` 和 `GAME_CONFIG_BASIC_AUTH_PASSWORD`（哈希前的原密码）。这组只读凭据不要复用有命令权限的 `GAME_API_TOKEN`，也不要加 `NEXT_PUBLIC_`。浏览器不能携带这组凭据，且当前游戏服务的 `/api/game` 是公开实时地图接口，因此它不能套用配置文件的 Basic Auth。

`/items` 和 `GET /api/v1/items?kind=sword` 公开展示道具；`/orders`、`GET /api/v1/orders` 与 `POST /api/v1/orders` 使用网站登录身份。所有单件道具可选等级及追加，只有配置允许卓越属性的道具可选卓越属性；套装仅展示。提交只创建待处理订单，不付款也不自动发货。道具、技能、套装与订单名称按语言显示：英文使用 XML 的 `Name`，中文使用 `annotation`（技能为 `anotation`），西班牙语暂用英文。浏览器 API 使用 `r2f2-locale` Cookie，App API 使用 `Accept-Language`，缺省为英文。首次部署订单功能前执行 `20261008000100_create_item_orders.sql`，否则订单 API 返回暂不可用。

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
