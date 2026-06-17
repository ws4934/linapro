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
- [x] 3.4 `ToResolverConfig`改为按 policy passthrough（honor`RootDomain`/`ReservedSubdomains`），移除硬编码强制空；默认`RootDomain`空保持 subdomain 禁用。说明：storefront 子域名以显式域名行经`domainResolver`覆盖，manifest 配置文件级`RootDomain`插桩后续迭代再补

## 4. 域名管理

- [ ] 4.1 `backend/api`新增域名 DTO 与路由契约（`g.Meta`、权限标签、时间字段约定）
- [ ] 4.2 域名`service`：CRUD 与验证标记，注入数据权限，唯一域名校验，操作前目标可见性校验
- [ ] 4.3 `controller`与路由绑定（`_new.go`构造结构）
- [ ] 4.4 `plugin.yaml`新增域名管理菜单与权限点；`manifest/i18n/<locale>`补充文案
- [ ] 4.5 `frontend/pages`新增最小域名管理页（列表、创建、删除、验证）

## 5. 测试

- [x] 5.1 `domainResolver`单测：已验证命中、未验证不匹配、停用租户不匹配、未知`Host`不匹配、`Host`大小写与端口规范化
- [x] 5.2 host-only 链单测：`StorefrontResolverChain()`排除`default`且含`domain`；未知`Host`经`domainResolver`不返回平台`0`
- [x] 5.3 既有名称测试补`domain`；`domainResolver`不回落平台由链结构与 resolve 行为共同保证
- [ ] 5.4 域名`service`单测：唯一性约束、数据权限注入、操作前可见性校验（随 task 4）

## 6. 验证与记录

- [ ] 6.1 运行`openspec validate tenant-store-domain-resolution --strict`并通过
- [ ] 6.2 记录 SQL 幂等性、数据分类（DDL 与 Seed 分离）、软删除与自动时间维护验证
- [ ] 6.3 记录数据权限例外：host-only 公开解析的权威边界（仅`verified`+`active`、未匹配即拒绝、永不平台）与测试覆盖
- [ ] 6.4 记录 DI 来源检查：域名 service 的 owner、创建位置、传递路径与共享实例判断（随 task 4）
- [ ] 6.5 记录影响分析：`i18n`（新增菜单与文案）、缓存一致性（无）、开发工具跨平台（无）、测试策略已覆盖
