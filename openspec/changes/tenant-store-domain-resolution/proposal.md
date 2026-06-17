## Why

storefront 按域名定店是商城地基，但`linapro-tenant-core`现有 resolver 链有三处缺口：无自定义域名（CNAME）解析（`subdomainResolver`只做`label == code`，覆盖不了`acme.com`这类整域名）、subdomain 解析被`resolverconfig`强制`RootDomain`空而禁用、且现有 resolver 链以`default`回落到平台`tenant_id = 0`，不能用于顾客匿名按域名访问店面。本变更在`tenant-core`内补齐域名解析能力，作为商城 P0 地基。

本变更范围闭环在`tenant-core`插件内，不改动`lina-core`核心；storefront 公开请求路径的接入由后续`storefront-host-surface`变更承担。

## What Changes

- 新增域名映射表`plugin_linapro_tenant_core_domain`（`tenant_id`↔域名，含`is_primary`、`is_verified`、`status`与审计字段），并经`make dao`生成 DAO/DO/Entity。
- 新增`domainResolver`：按`request.Host`全量匹配已验证且租户有效的域名解析到`tenant_id`，注册进现有 resolver 链。
- 开启 subdomain 解析：`resolverconfig`从插件配置读取`RootDomain`与保留子域名，不再硬编码为空。
- 提供 host-only 解析能力：通过`Config.Chain`仅运行域名与子域名 resolver、排除`default`，保证未匹配时不回落平台`tenant_id = 0`，供`storefront-host-surface`消费。
- 新增域名管理 API（平台作用域 CRUD 与验证标记）、权限点、菜单与`i18n`资源，并接入数据权限。
- 新增单元测试：域名匹配、host-only 链永不回落平台、未知域名不匹配、subdomain 启用与保留子域名拦截。

## Capabilities

### New Capabilities

- `tenant-domain-resolution`：租户自定义域名与子域名解析、host-only 公开解析的不回落平台保证，以及域名映射的管理与验证。

### Modified Capabilities

本变更为新增能力，不修改既有规范需求。既有`tenant-platform-access-control`与基线`commerce-platform-architecture`（storefront 隔离红线）作为遵循依据被引用。

## Impact

- 改动范围：`apps/lina-plugins/linapro-tenant-core`（fork 的 submodule）——`backend/internal/service/resolver`、`resolverconfig`、`shared`，新增 domain 的`service`/`dao`/`model`/`controller`/`api`，`manifest/sql`、`manifest/config`、`manifest/i18n`与`frontend`域名管理页。
- 数据权限：域名管理 CRUD 为平台作用域并注入数据权限；`host-only`解析是公开解析例外（无身份），其权威边界为「仅匹配`is_verified`且租户`active`的域名、未匹配即拒绝、永不返回平台」，在 design 与 specs 显式说明并测试覆盖。
- 依赖：无新增运行期依赖；新增 domain service 经构造函数注入，复用现有`bizctx`与 DB 访问。
- 不涉及：`lina-core`核心改动；store 业务属性（货币、套餐、主题绑定）由后续`commerce-store-profile`承担，不进本变更。
