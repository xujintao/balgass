<!--
首次使用时，将 AGENTS.example.md 复制为同目录的 AGENTS.md：
cp AGENTS.example.md AGENTS.md
共享规则请更新到 AGENTS.example.md；本地定制保留在被 Git 忽略的 AGENTS.md 中。
-->

# 项目指引

## 项目目标

使用 Next.js、Vercel 和 Supabase 重写原 Python/Django 网站。当前目录 `server-web-next` 是新项目的工作目录。

## 原项目参考

需要追溯原 Python/Django 网站行为时，请查阅 Git 历史中的旧项目路由、视图、模型、模板和配置。以旧代码实际实现的行为作为功能参考。

区分功能复现与新增功能。原项目尚未完成的购物车、支付、发货等部分，不自动纳入本次重写需求。

## 架构原则

- Next.js 负责网站页面渲染、HTTP API 和服务端业务逻辑，Vercel 负责网站部署。
- Supabase 负责网站用户认证与网站数据存储。
- 现有游戏服务继续负责游戏账号和实时状态，新网站通过其接口接入。

## 浏览器与 App 支持

Next.js 同时支持浏览器页面和供未来 App 使用的 HTTP API。

- 使用 App Router 实现页面，使用 Route Handlers 提供 `/api/v1/...` JSON API。
- 页面与 API 共用业务服务层，统一处理权限和业务规则。Server Components 直接调用业务服务层，避免通过 HTTP 请求本站 API。
- 面向 App 的业务功能必须提供 HTTP API，不能仅通过 Server Actions 实现。
- 浏览器使用 Cookie 会话，App 使用 Bearer Token；服务端验证身份后形成统一的用户上下文。
- API 保持稳定的请求、响应、状态码和错误代码，并维护 OpenAPI 文档。
- API 认证失败时返回 JSON 错误及相应状态码，不重定向到登录页面。
- 用户私有数据不得进入共享缓存。

## 安全要求

- 在服务端校验身份、权限和输入，不依赖浏览器端校验来保护数据或操作。
- 为用户私有数据配置 Supabase 行级安全策略（RLS），确保用户只能访问获授权的数据。
- 服务端密钥仅保存在服务端，不得暴露给浏览器或提交到版本库。

## 开发与验证

使用项目实际配置的包管理器、开发命令和检查命令，不预设尚未确定的工具或命令。

完成修改后，运行与改动相关且项目已有的检查，并说明验证结果。如果检查无法运行，明确说明原因，不声称检查通过。

## 提交规范

提交信息使用 `server-web-next/<type>: <英文小写描述>` 格式。根据改动内容选择类型，沿用仓库已有的 `feat`、`fix`、`refactor`、`docs` 等类型。

例如：`server-web-next/feat: add player game account management`。
