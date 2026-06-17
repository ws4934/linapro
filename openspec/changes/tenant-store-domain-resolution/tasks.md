## 1. 规则与前置

- [x] 1.1 读取命中规则`backend-go.md`、`api-contract.md`、`i18n.md`（`plugin.md`、`data-permission.md`、`database.md`、`architecture.md`已读），记录无遗漏
- [x] 1.2 确认`linapro-tenant-core`无插件本地`AGENTS.md`，按项目顶层规范与命中规则执行

## 2. 数据模型

- [x] 2.1 在`manifest/sql/002-tenant-store-domain-resolution.sql`新增域名表（`CREATE TABLE IF NOT EXISTS`，`tenant_id`、`domain`、`is_primary`、`is_verified`、`verification_token`、`status`与`created_at`/`updated_at`/`deleted_at`），`domain`唯一索引、`tenant_id`与解析索引，满足幂等
- [x] 2.2 在`backend/hack/config.yaml`覆盖域名表；生成 DAO/DO/Entity（codegen 阶段以幂等`002`SQL 直接建表于 dev 库后`make dao`；运行期建表走插件安装 SQL 流程）
- [x] 2.3 通过 Go 编译门禁（`GOWORK=off go build ./...`exit 0，`go vet`变更包 exit 0）

## 3. 解析能力

- [x] 3.1 `shared`新增`ResolverDomain`常量与`TableDomain`表名、`DomainStatus`枚举；新增`domainResolver`（规范化`Host`、查域名表`is_verified`且`status=active`、再校验租户`active`、命中返回`tenant_id`）
- [x] 3.2 `resolver.New()`注册`domainResolver`；默认链顺序在`subdomain`之后、`default`之前加入`domain`
- [x] 3.3 新增 host-only 链`StorefrontResolverChain()`（`[domain, subdomain]`，排除`default`），供 storefront 复用
- [x] 3.4 `ToResolverConfig`改为按 policy passthrough，移除硬编码强制空；默认`RootDomain`空保持 subdomain 禁用。storefront 子域名以显式域名行经`domainResolver`覆盖，manifest 配置文件级`RootDomain`插桩后续迭代再补

## 4. 域名管理

- [x] 4.1 `api/platform/v1`新增域名 DTO（list/create/delete/verify），`g.Meta`含`dc`/`eg`/`permission`，`createdAt`为 Unix 毫秒`*int64`
- [x] 4.2 域名`service`（`internal/service/domain`）：List/Create/Delete/SetVerified，typed 插入/更新结构，域名归一化+全局唯一校验，操作前可见性校验，bizerr 错误码集中于`domain_code.go`
- [x] 4.3 `make ctrl`生成`IPlatformV1`与控制器骨架并填充 4 个域名方法；`platform.go`加`domainSvc`与`toAPIDomain`；`plugin.go`装配`domainSvc`并绑定
- [x] 4.4 `plugin.yaml`新增 4 个域名隐藏权限点（list/add/remove/verify）；`manifest/i18n/{en-US,zh-CN}/apidoc/platform/domain.json`补 apidoc 翻译。可见菜单随前端页一并落地
- [x] 4.5 `frontend/pages`域名管理页（`domain-client.ts` + `domain-management.vue` + `domain-modal.vue`）+ 可见菜单（`plugin.yaml`域名菜单+按钮）+ `menu.json`与 lina-vben `pages.json`双语文案；镜像既有租户页模式，JSON 合法、复用键存在。前端 vite 构建、`vue-tsc`类型检查与 E2E 须在前端环境验证（本环境无 vben 工具链）

## 5. 测试

- [x] 5.1 `domainResolver`单测：已验证命中、未验证不匹配、停用租户不匹配、未知`Host`不匹配、`Host`大小写与端口规范化
- [x] 5.2 host-only 链单测：`StorefrontResolverChain()`排除`default`且含`domain`；未知`Host`经`domainResolver`不返回平台`0`
- [x] 5.3 既有名称测试补`domain`；`domainResolver`不回落平台由链结构与 resolve 行为共同保证
- [x] 5.4 域名`service`单测：归一化、全局唯一拒绝、无效输入、删除 not-found、验证标记、按租户筛选（连容器 DB，`ok`通过）

## 6. 验证与记录

- [x] 6.1 运行`openspec validate tenant-store-domain-resolution --strict`并通过
- [x] 6.2 SQL 幂等性：建表/索引均`IF NOT EXISTS`；数据分类：仅 DDL 无 Seed；软删除与时间维护由`deleted_at`/`created_at`/`updated_at`自动处理，未手写
- [x] 6.3 数据权限例外：域名映射为平台治理数据，管理按平台权限`system:tenant:domain:*`门控，不施加行级租户数据范围（与既有平台租户管理同构）；host-only 公开解析权威边界=仅`verified`+`active`、未匹配即拒绝、永不平台，单测覆盖
- [x] 6.4 DI 来源检查：`domainSvc`owner=`registerRoutes`启动装配，经`domainsvc.New(services.BizCtx())`创建并构造注入平台控制器，无运行期`New()`独立服务图、无共享缓存状态
- [x] 6.5 影响分析：`i18n`新增 4 权限点、域名 apidoc 双语、菜单与 lina-vben 页面文案双语；缓存一致性无影响；开发工具跨平台无影响；后端测试覆盖解析与 service 路径
- [x] 6.6 i18n 治理验证：`make i18n.check` exit 0、零 error、domain 文案/菜单/apidoc 双语完整、无新增模块级 `$t` 警告
- [x] 6.7 前端构建/类型检查：经 Volta 安装 `pnpm@10.30.3` 并 `pnpm install`；web-antd `vue-tsc --noEmit --skipLibCheck` exit 0（`src` + `pages.json` 干净）；`turbo build --filter=@lina/web-antd` 11/11 成功（应用 + i18n 构建干净）。插件域名页忠实镜像既有租户页；standalone `vue-tsc` 报错均为「包外文件模块解析够不到 web-antd `node_modules`」产物与租户页一致的 implicit-any，且插件页在真实流程由 esbuild 转译不做类型门禁，无真实缺陷
## 7. E2E — TC001 域名管理页

- [x] 7.1 创建 `apps/lina-plugins/linapro-tenant-core/hack/tests/pages/DomainManagementPage.ts`（POM，mock 工作台 shell + 域名 CRUD API，零真实后端依赖）
- [x] 7.2 创建 `apps/lina-plugins/linapro-tenant-core/hack/tests/e2e/tenant-domain/TC001-domain-management.ts`，覆盖：TC-1a 列表渲染 + i18n 列标题/按钮实际文案断言（非双语正则）；TC-1b 创建并出现在列表 + 断言 POST 载荷；TC-1c 重复域名唯一约束拒绝（HTTP 400）且不新增行；TC-1d 验证开关写入 `verified` 状态；TC-1e 删除后行消失
- [x] 7.3 `node hack/tests/scripts/validate-e2e.mjs` 通过（251 文件零错误，TC 命名/分组/结构合规）
- [ ] 7.4 浏览器运行（remaining gate）：须 `hack/tests` 依赖 + Playwright chromium + 在 `:9120` 提供本插件域名页的前端栈（vite dev / 后端服务）。本环境未完成该运行时编排。替代验证：后端解析 6 例 + service 5 例连库全绿、前端 `turbo build` 11/11、`vue-tsc` exit 0、`make i18n.check` exit 0、E2E 结构校验通过。残余风险：运行时浏览器交互未跑通，由 mocked-API 设计 + 镜像既有 `MultiTenantPage` 模式缓解
