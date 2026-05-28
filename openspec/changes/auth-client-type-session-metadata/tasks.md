## 1. 规范和影响确认

- [x] 1.1 读取并记录命中的 OpenSpec、文档、架构、API、后端 Go、数据库、缓存一致性、数据权限、插件、前端 UI、测试和 i18n 规则。
- [x] 1.2 确认`clientType`仅表示用户会话客户端，枚举值只包含`web`、`mobile`、`desktop`、`cli`。
- [x] 1.3 记录破坏式变更范围：登录请求必须显式传入合法`clientType`，后端认证内核不再隐式补`web`。

## 2. 认证模型实现

- [x] 2.1 在认证服务中新增`ClientType`命名类型、常量、解析校验函数和结构化错误码。
- [x] 2.2 将`ClientType`写入`LoginInput`、`preTokenRecord`、JWT claims、租户令牌签发、租户切换、刷新和 impersonation 令牌路径。
- [x] 2.3 将`Logout`改为接收`LogoutInput`并从当前业务上下文传入`ClientType`，移除退出 Hook 中的硬编码`web`。
- [x] 2.4 更新认证中间件和`bizctx`，把 claims 中的`ClientType`注入请求业务上下文。

## 3. 会话存储和插件投影

- [x] 3.1 更新宿主`004-online-session.sql`建表源，使`sys_online_session.client_type`成为必填投影字段。
- [x] 3.2 执行 DAO 生成流程并更新`sys_online_session`相关 DAO/DO/Entity 生成文件。
- [x] 3.3 更新`session.Session`、DBStore、CoordinationStore Redis hot state 和在线会话列表投影，保证`ClientType`随会话写入和读取。
- [x] 3.4 更新插件 session contract、hostservices adapter 和`linapro-monitor-online`在线用户列表 DTO/投影，返回`clientType`。

## 4. API、前端和 i18n

- [x] 4.1 更新`LoginReq`，新增必填`clientType`字段和 OpenAPI 文档说明。
- [x] 4.2 更新默认 Web 前端登录 API 适配层，显式传入`clientType=web`。
- [x] 4.3 同步宿主`zh-CN` apidoc 翻译资源，确认无新增前端用户可见文案。

## 5. 测试和验证

- [x] 5.1 新增或更新认证单元测试，覆盖合法/非法`clientType`、登录失败 Hook、二阶段租户选择、刷新继承、租户切换继承和退出 Hook。
- [x] 5.2 新增或更新会话存储与 hostservices 投影测试，覆盖 DB/Redis hot state/session contract 中的`ClientType`。
- [x] 5.3 运行 Go 验证：`cd apps/lina-core && go test ./internal/service/auth ./internal/service/session ./internal/service/plugin/internal/hostservices -count=1`。
- [x] 5.4 涉及 API/路由绑定后运行`cd apps/lina-core && go test ./internal/cmd -count=1`。
- [x] 5.5 运行前端类型检查、OpenSpec 严格校验、静态检索和格式检查。

## 6. 执行记录

- [x] 6.1 记录 i18n、缓存一致性、数据权限、插件边界、前端 UI、开发工具跨平台和测试影响分析。
- [x] 6.2 记录 SQL 幂等性、数据分类、索引判断、DAO 生成和编译门禁结果。
- [x] 6.3 完成实现和验证后执行`lina-review`审查。

## 执行记录

### 规范和范围确认

- 已按`AGENTS.md`读取本次命中的规则：OpenSpec、文档、架构、API、后端 Go、数据库、缓存一致性、数据权限、插件、前端 UI、测试和`i18n`。Go 后端实现同步使用`goframe-v2`技能；实现任务使用`openspec-apply-change`流程。
- 根因确认：认证 Hook 的`clientType`当前只在登录失败、登录成功、退出成功 payload 组装处硬编码为`web`，且插件分发层对空值继续补`web`；该字段未进入登录输入、`pre_token`、JWT claims、在线会话、Redis hot state 或`bizCtx`，导致`Logout`无法从当前会话事实源读取客户端类型。
- 设计范围确认：`clientType`只表示用户会话客户端，枚举值限定为`web`、`mobile`、`desktop`、`cli`；`service`和`plugin`属于未来独立的主体类型、授权方式或 actor 语义，不进入`ClientType`。
- 破坏式变更确认：登录请求必须显式提交合法`clientType`；默认 Web 工作台由前端 API 适配层显式传入`web`，后端认证内核、插件 Hook 分发层和退出逻辑均不再隐式补默认值。

### 实现记录

- 认证模型：在`pkg/authtoken`新增共享`ClientType`命名类型、`web/mobile/desktop/cli`常量和基础校验函数；`auth`服务复用该类型并通过`CodeAuthClientTypeInvalid`返回结构化错误。`plugin`、`service`和空值均被拒绝。
- 生命周期传递：`ClientType`已进入`LoginInput`、`preTokenRecord`、JWT claims、`IssueTenantToken`、`ReissueTenantToken`、`Refresh`和 impersonation 签发路径；impersonation 从当前`bizCtx.ClientType`继承用户会话来源，缺少认证上下文时拒绝签发。
- 退出路径：`Logout`改为接收`LogoutInput`，控制器从当前`bizCtx`传入`ClientType`，退出 Hook 使用该事实源，不再硬编码`web`。
- 请求上下文：认证中间件从 JWT claims 注入`bizctx.SetUser(..., clientType)`，`model.Context`保留`ClientType`供控制器、动态插件路由和 host service 适配使用；动态插件路由解析 JWT 时也通过`pkg/authtoken.ParseClientType`拒绝空值、`plugin`和`service`。
- 会话投影：`session.Session`、DBStore、Redis session hot state payload、`sys_online_session` DAO/DO/Entity 均新增`ClientType/client_type`字段；读写、分页列表和 scoped 列表均随既有投影返回，不新增查询。
- 插件契约：`pkg/plugin/capability/contract.Session`、hostservices adapter 和`linapro-monitor-online`在线用户列表 DTO/控制器投影均返回`clientType`；插件目录无本地`AGENTS.md`，`plugin.yaml`启用`i18n`，已维护插件自身`zh-CN` apidoc 翻译。
- API 与前端：`LoginReq.clientType`为必填字段；默认 Web 前端登录 API 适配层显式提交`clientType=web`，没有新增前端页面、交互或运行时可见文案。

### 影响分析

- `i18n`：宿主新增错误键、校验键和`LoginReq.clientType`接口文档翻译；同步维护`manifest/i18n/zh-CN/apidoc/core-api-auth.json`与打包副本。英文 apidoc 按规则保持源文本和空占位策略。在线用户插件启用`i18n`，已补充插件自身`zh-CN` apidoc 翻译；无新增前端运行时文案。
- 缓存一致性：Redis session hot state 仅扩展 payload 字段，权威来源仍为登录时认证输入和 JWT/session 投影；key 维度、TTL、失效路径、跨实例共享 KV 和故障降级策略不变。`pre_token`继续通过共享 KV 单次消费，额外保存`ClientType`。
- 数据权限：没有新增数据读取或写入边界；在线用户列表仍走`sys_online_session`投影、tenantcap 和 datascope 过滤，新增字段随同一行投影返回，不通过逐项补查暴露额外存在性。
- 插件边界：插件只接收宿主发布的 session contract 投影字段，不依赖宿主 DAO/DO/Entity 或内部缓存结构；动态插件路由校验 JWT claims 的`clientType`后再透传到用户上下文。
- 前端 UI：仅修改 Web 登录 API 适配层请求体，未改页面、表单、按钮、路由、布局或用户可见交互，因此无需新增 E2E；通过类型检查覆盖调用契约。
- 开发工具跨平台：没有修改 Makefile、脚本、CI 或`linactl`；仅使用既有`make dao`生成流程。命中开发工具规则的影响为无新增跨平台入口。
- DI 来源：本次没有新增运行期服务依赖、构造函数参数或独立服务图；只扩展既有接口参数、DTO、claims 和投影字段。缓存敏感服务仍复用启动期传入的`sessionStore`、`kvCacheSvc`、`roleSvc`和`tenantSvc`。
- 性能：`clientType`随登录、会话写入、会话读取、Redis hot state 和在线用户列表既有查询路径读写，不新增数据库查询、远程调用或前端瀑布式请求；当前没有按`client_type`过滤或排序需求。

### SQL、生成和验证记录

- SQL：更新`004-online-session.sql`建表源，直接在`sys_online_session`中声明`"client_type" VARCHAR(32) NOT NULL`和列注释；按破坏式设计不保留独立追加迁移、不设置默认值、不回填历史数据。`CREATE TABLE IF NOT EXISTS`和列注释保持初始化脚本可重复执行。
- 数据分类：本次 SQL 仅包含 DDL 和注释，无 Seed DML、Mock 数据、自增主键写入或软删除语义变化。
- 索引判断：本次没有新增按`client_type`筛选、排序、聚合或关联查询路径，因此不新增低价值索引；未来出现在线用户按客户端类型筛选时再按实际查询增加组合索引。
- DAO 生成：首次`make init confirm=init`因代码已引用未生成的`do.SysOnlineSession.ClientType`失败；本地开发库存在旧行且本机无`psql`，因此使用一次性本地 schema 修复辅助流程让开发库具备该列后执行`cd apps/lina-core && make dao`成功。临时本地修复未提交，提交的迁移仍保持破坏式无默认/无回填。
- 生成结果：`sys_online_session`相关 DAO/DO/Entity 已通过`make dao`更新；`sys_kv_cache.ValueBytes`在同次生成中回到实际`BYTEA`对应的`[]byte`类型，既有 kvcache 代码按该类型使用。
- 验证通过：
  - `cd apps/lina-core && go test ./internal/service/auth ./internal/service/session ./internal/service/plugin/internal/hostservices ./internal/service/apidoc -count=1`
  - `cd apps/lina-core && go test ./pkg/authtoken ./internal/service/plugin/internal/runtime ./internal/service/kvcache -count=1`
  - `cd apps/lina-core && go test ./internal/cmd -count=1`
  - 使用临时`go.work`包含`apps/lina-core`和`apps/lina-plugins/linapro-monitor-online`后执行`GOWORK=<tmp>/go.work go test ./... -count=1`
  - `cd apps/lina-vben && pnpm -F @lina/web-antd typecheck`
  - `openspec validate auth-client-type-session-metadata --strict`
  - 静态检索生产路径确认不存在`ClientType: "web"`、`input.ClientType = "web"`或`clientType = "web"`后端兜底硬编码；`plugin/service`仅出现在拒绝用例中。
  - `git diff --check`覆盖本次相关文件，未发现空白格式问题。

### Lina 审查记录

- 变更：`auth-client-type-session-metadata`。
- 范围：全部变更。审查范围来源为`git status --short`、`git ls-files --others --exclude-standard`、`git -C apps/lina-plugins status --short`和本变更任务上下文；审查文件数 53。当前工作区存在大量`consolidate-plugin-service-boundaries`等无关既有变更，已排除在本次审查结论之外。
- 已读取规则文件：`AGENTS.md`、`.agents/rules/documentation.md`、`.agents/rules/openspec.md`、`.agents/rules/architecture.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`。使用技能：`lina-review`、`goframe-v2`、`openspec-apply-change`。
- 审查发现：初审发现动态插件路由的 JWT 解析路径读取`clientType`后未校验，已修复为复用`pkg/authtoken.ParseClientType`，并补充拒绝空值、`plugin`、`service`的单元测试。复查未发现阻塞问题。
- 规则域结论：OpenSpec、文档、架构、API、后端 Go、数据库、缓存一致性、数据权限、插件、前端 UI、测试、`i18n`和开发工具跨平台均通过。无新增运行期依赖；无新增前端可见 UI；无 E2E 必要性，原因是默认 Web 前端仅改变 API 适配层请求体且已有类型检查覆盖。
- 额外生成影响：`make dao`同次更新`sys_kv_cache.ValueBytes`为`[]byte`，已通过`cd apps/lina-core && go test ./internal/service/kvcache -count=1`覆盖。
- 验证证据：所有“SQL、生成和验证记录”中的命令均已在当前工作区重跑通过。`openspec validate auth-client-type-session-metadata --strict`通过；静态扫描确认生产路径无后端`web`兜底硬编码，且`service/plugin`未进入`ClientType`枚举。
- 摘要：严重 0，警告 0。剩余风险：提交的破坏式建表源不兼容已有非空旧表数据，符合“不考虑兼容性”的用户要求和设计记录。

## Feedback

- [x] **FB-1**: 将`013-auth-client-type-session-metadata.sql`整理回原始在线会话 SQL 数据文件
- [x] **FB-2**: 修复在线会话`client_type`必填后 CI mock 数据初始化失败
- [x] **FB-3**: 修复 CI 冒烟登录请求缺少必填`clientType`
- [x] **FB-4**: 修复 E2E 直接登录请求缺少必填`clientType`

### Feedback 执行记录

#### FB-1：整理在线会话 SQL 文件

- 根因确认：`sys_online_session.client_type`已经同步写入`004-online-session.sql`建表源，但独立的`013-auth-client-type-session-metadata.sql`仍残留，导致同一破坏式 schema 变更同时存在原始建表源和追加迁移两种表达，不符合用户要求的“修改到原来的 SQL 数据文件中”。
- 修复内容：删除`apps/lina-core/manifest/sql/013-auth-client-type-session-metadata.sql`，保留`apps/lina-core/manifest/sql/004-online-session.sql`中的`client_type`字段和列注释；同步更新`proposal.md`、`design.md`和本任务记录中关于 SQL 文件形态的描述。
- `i18n`影响：无运行时用户可见文案、API 文档源文本、错误消息、插件清单、语言包或翻译缓存变更。
- 缓存一致性影响：无缓存 key、payload、失效、刷新、跨实例同步或故障降级策略变更。
- 数据权限影响：无新增或修改数据读写接口、过滤条件、租户边界、存在性暴露或插件宿主数据访问路径。
- 开发工具跨平台影响：未修改 Makefile、脚本、CI、`linactl`或代码生成入口；本次仅整理 SQL 和 OpenSpec 文档。
- 测试策略：纯项目治理类反馈，不改变运行时行为，不新增单元测试或 E2E；使用 OpenSpec 严格校验、SQL 静态检索、文件存在性检查和格式检查验证。
- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`。
- 验证记录：`rg -n "client_type|auth-client-type-session-metadata" apps/lina-core/manifest/sql`确认`client_type`仅在`004-online-session.sql`中声明且无`013`文件引用；`ls apps/lina-core/manifest/sql`确认`013-auth-client-type-session-metadata.sql`不存在。

#### FB-2：修复 CI mock 数据初始化失败

- 根因确认：GitHub Actions run `26513008846`在 Host-only build smoke 阶段执行 mock 数据时失败，PostgreSQL 报错`null value in column "client_type" of relation "sys_online_session" violates not-null constraint`。`sys_online_session.client_type`已在源码 SQL 中改为`NOT NULL`，但宿主 mock 数据、打包 mock 数据、打包建表 SQL 和部分测试 fixture 仍按旧列清单写入在线会话，导致初始化或测试插入时缺少必填字段。
- 修复内容：为宿主`manifest/sql/mock-data/004-online-sessions.sql`和`linapro-monitor-online`插件 mock 数据补齐`client_type`列和值；Web/桌面浏览器演示会话使用`web`，iOS 演示会话使用`mobile`。同步修复`cluster_multiprocess_test.go`和`linapro-tenant-core` E2E support 中直接插入`sys_online_session`的 fixture，避免后续测试路径继续触发同一 NOT NULL 约束失败。运行`make pack.assets`确认忽略的`internal/packed/manifest`可由源码 manifest 重新生成。
- `i18n`影响：无运行时用户可见文案、API 文档源文本、错误消息、插件清单、语言包或翻译缓存变更；仅补齐 mock/fixture 数据字段。
- 缓存一致性影响：无缓存 key、payload、失效、刷新、跨实例同步或故障降级策略变更；仅影响初始化演示数据和测试 fixture 写入。
- 数据权限影响：不新增接口或数据可见性边界；`linapro-tenant-core` fixture 保持原有 tenant A/B 隔离断言，只是补齐同一行的必填`client_type`字段。
- 开发工具跨平台影响：未修改 Makefile、脚本、CI 或`linactl`实现；使用既有跨平台`make pack.assets`、`make init`和`make mock`入口验证。
- 测试策略：这是可执行初始化数据修复，使用真实 PostgreSQL 隔离容器跑 CI 同类`make init`和`make mock`路径复现验证；另用 Go 包编译、E2E TypeScript 编译、E2E 结构校验、静态检索、OpenSpec 严格校验和格式检查覆盖相关回归。
- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/backend-go.md`、`.agents/rules/plugin.md`；Go 后端测试 fixture 修复同步使用`goframe-v2`技能。
- 验证记录：`make pack.assets`通过；临时 PostgreSQL 容器`linapro-ci-fb2-postgres`上以隔离`GF_GCFG_PATH`执行`make init confirm=init rebuild=true`通过，随后`make mock confirm=mock`通过，确认`004-online-sessions.sql`不再触发`client_type`非空约束失败；`cd apps/lina-core && go test ./internal/service/cluster -run TestClusterTwoHostProcessesSharePostgreSQL -count=1`通过包编译并按缺少`LINA_TEST_PGSQL_LINK`/`LINA_TEST_REDIS_ADDR`跳过真实多进程测试；`pnpm -C hack/tests exec tsc --noEmit`通过；`pnpm -C hack/tests test:validate`通过，校验 238 个 E2E 文件；`openspec validate auth-client-type-session-metadata --strict`通过；`git diff --check -- apps/lina-core/internal/service/cluster/cluster_multiprocess_test.go apps/lina-core/manifest/sql/mock-data/004-online-sessions.sql openspec/changes/auth-client-type-session-metadata/tasks.md`和`git -C apps/lina-plugins diff --check -- linapro-monitor-online/manifest/sql/mock-data/001-linapro-monitor-online-mock-data.sql linapro-tenant-core/hack/tests/support/linapro-tenant-core-scenarios.ts`通过；静态检索`rg -n -U -P 'INSERT INTO sys_online_session\s*\((?:(?!\)).)*\)' apps/lina-core apps/lina-plugins hack --glob '!**/node_modules/**' --glob '!**/dist/**' --glob '!**/.git/**'`确认一行式 fixture 插入均包含`client_type`，并通过上下文检索确认宿主与插件 mock SQL 的每段在线会话插入均包含`"client_type"`。

##### FB-2 Lina 审查记录

- 审查范围：反馈级，文件包括`apps/lina-core/manifest/sql/mock-data/004-online-sessions.sql`、`apps/lina-core/internal/service/cluster/cluster_multiprocess_test.go`、`apps/lina-plugins/linapro-monitor-online/manifest/sql/mock-data/001-linapro-monitor-online-mock-data.sql`、`apps/lina-plugins/linapro-tenant-core/hack/tests/support/linapro-tenant-core-scenarios.ts`和本`tasks.md`记录；`internal/packed/manifest`为`.gitignore`忽略的生成资源，已通过`make pack.assets`验证由源码 manifest 重新生成，不纳入提交范围。
- 已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/backend-go.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`；技能：`lina-feedback`、`lina-review`、`goframe-v2`、`lina-e2e`。
- 规则域结论：OpenSpec 反馈先记录根因再修复，任务已标记完成且有验证证据；SQL mock 数据与测试 fixture 均补齐`client_type`，不新增自增主键写入或非幂等 DML；Go 后端只改测试 fixture，包编译门禁通过；插件目录未发现本地`AGENTS.md`，按顶层插件/测试规则执行；E2E support TypeScript 编译和结构校验通过，无新增 E2E 文件编号影响；无运行时 UI/API 文案、缓存、数据权限或开发工具实现变更。
- 验证证据：复查`openspec validate auth-client-type-session-metadata --strict`通过；复查`git diff --check`覆盖宿主和插件本次文件通过；审查 diff 未发现旧式`sys_online_session`插入缺少`client_type`的残留。严重问题 0，警告 0。

#### FB-3：修复 CI 冒烟登录请求缺少`clientType`

- 根因确认：GitHub Actions run `26515440723`中 Host-only build smoke 和 Redis cluster smoke 均已成功启动后端，但冒烟脚本仍按旧登录契约请求`/api/v1/auth/login`，请求体只有`username`和`password`。`LoginReq.clientType`已经按本变更改为必填，因此服务正确返回`validation.auth.login.clientType.required`，导致 smoke 失败。
- 修复内容：在`.github/actions/host-only-artifact-smoke/action.yml`和`hack/tests/scripts/run-redis-cluster-smoke.sh`的登录请求体中显式传入`clientType:"web"`，让 CI smoke 与默认 Web 工作台登录契约保持一致。
- `i18n`影响：无运行时用户可见文案、API 文档源文本、错误消息、插件清单、语言包或翻译缓存变更。
- 缓存一致性影响：无缓存 key、payload、失效、刷新、跨实例同步或故障降级策略变更；Redis smoke 只是通过正确登录参数验证既有集群 session store 写入路径。
- 数据权限影响：无新增接口、数据操作、租户边界或存在性暴露；仅修复 CI 登录请求参数。
- 开发工具跨平台影响：修改 GitHub composite action 和 Bash smoke 脚本；它们属于 CI 专属 Linux 入口，不作为默认跨平台开发入口。Windows 命令 smoke 不受影响。
- 测试策略：这是 CI 治理和运行时接口联动修复，不新增业务单元测试；使用静态检索、格式检查、OpenSpec 严格校验，并运行相关 Go smoke/编译门禁验证同一 Action 中其他失败项。
- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/api-contract.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`。
- 验证记录：静态检索确认`.github/actions/host-only-artifact-smoke/action.yml`和`hack/tests/scripts/run-redis-cluster-smoke.sh`的登录 payload 均包含`clientType:"web"`；`openspec validate auth-client-type-session-metadata --strict`通过；`git diff --check`通过。关联 CI 中 plugin-full 失败项已在`consolidate-plugin-service-boundaries`的`FB-3`中继续修复和验证。

##### FB-3 Lina 审查记录

- 审查范围：反馈级，文件包括`.github/actions/host-only-artifact-smoke/action.yml`、`hack/tests/scripts/run-redis-cluster-smoke.sh`和本`tasks.md`记录；`git ls-files --others --exclude-standard`无未跟踪文件。
- 已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/api-contract.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`；技能：`lina-feedback`、`lina-review`。
- 规则域结论：OpenSpec 反馈先记录根因再修复且任务完成；CI composite action 和 Bash smoke 属于 CI/Linux 专属入口，已记录跨平台边界；登录请求与必填`clientType`接口契约一致；无运行时 UI、API 文档源文本、语言包、缓存、数据权限或生产后端行为变更。
- 验证证据：`openspec validate auth-client-type-session-metadata --strict`通过；静态检索确认两个 smoke 登录 payload 均包含`clientType:"web"`；`git diff --check`通过。严重问题 0，警告 0。

#### FB-4：修复 E2E 直接登录请求缺少`clientType`

- 根因确认：GitHub Actions run `26518762946`中 Go 单测、前端单测、命令 smoke 和构建 smoke均通过，失败集中在 E2E 矩阵；失败日志显示多个用例在`hack/tests/support/api/job.ts:157`断言登录响应`payload.code`时收到业务码`51`。`POST /api/v1/auth/login`已按本变更收紧为必须显式提交`clientType`，前端登录适配层和 CI smoke 已补齐，但宿主 E2E helper、宿主插件框架 E2E、源码插件 E2E 和插件测试 support 中仍有多处直接调用`auth/login`只传`username/password`，导致 E2E 前置 API 登录按新契约被拒绝。
- 修复内容：为所有 E2E 直接登录请求补齐`clientType:"web"`，覆盖宿主共享 helper、宿主插件框架用例、配置/菜单/用户下拉用例，以及`linapro-content-notice`、`linapro-demo-dynamic`、`linapro-demo-source`、`linapro-monitor-online`、`linapro-ops-demo-guard`、`linapro-tenant-core`插件测试路径。同步增强`hack/tests/scripts/validate-e2e.mjs`，将宿主测试资产和插件`hack/tests/{e2e,pages,support}`中的直接`auth/login`请求纳入治理校验，后续缺少`clientType`会在`pnpm -C hack/tests test:validate`中失败。
- `i18n`影响：无运行时用户可见文案、API 文档源文本、错误消息、插件清单、语言包或翻译缓存变更；仅修复测试请求体和测试治理校验。
- 缓存一致性影响：无缓存 key、payload、失效、刷新、跨实例同步或故障降级策略变更；E2E 登录仍使用正常 Web 用户会话路径。
- 数据权限影响：无新增接口、数据操作、租户边界或存在性暴露；租户插件测试保持原有租户隔离断言，只补齐同一登录请求的必填客户端类型。
- 插件边界影响：插件目录均未发现本地`AGENTS.md`；变更仅发生在插件自有`hack/tests`测试资产中，不改插件生产后端、清单、host service、WASM bridge 或插件能力契约。
- 前端 UI 影响：不修改页面、路由、组件、表单、按钮、布局或用户可观察交互；仅修复 Playwright APIRequestContext/fetch 直连后端的测试前置登录。
- 开发工具跨平台影响：修改 Node E2E 治理脚本`validate-e2e.mjs`，该脚本由`pnpm -C hack/tests test:validate`调用，使用 Node 标准库和现有治理工具函数，适用于 macOS/Linux/Windows 的 Node 运行环境；无新增 Bash/PowerShell/Makefile/`linactl`入口。
- 测试策略：这是 E2E 测试和治理脚本修复，不新增业务 E2E 用例编号；使用 TypeScript 编译、E2E 结构校验、静态扫描、OpenSpec 严格校验和格式检查验证。未运行完整 Playwright 矩阵，原因是该矩阵依赖完整服务启动和插件分片，成本高且本次缺陷可由登录请求静态门禁直接覆盖。
- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/api-contract.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`；技能：`lina-feedback`、`lina-e2e`。
- 验证记录：`pnpm -C hack/tests exec tsc --noEmit`通过；`pnpm -C hack/tests test:validate`通过，校验 238 个 E2E 文件；静态扫描确认 37 个包含`auth/login`的文件中，测试直接请求窗口均包含`clientType`；`openspec validate auth-client-type-session-metadata --strict`通过；`git diff --check`覆盖宿主 E2E、插件 E2E/support、治理脚本和本任务记录通过。`gh run view 26518762946 --repo linaproai/linapro --log-failed`已确认旧提交失败点为 E2E 登录响应业务码`51`；E2E artifact zip 下载曾因网络读超时失败，但不影响日志证据。

##### FB-4 Lina 审查记录

- 审查范围：反馈级，文件包括宿主`hack/tests`下 14 个 E2E/helper/治理脚本文件、插件仓库`linapro-content-notice`、`linapro-demo-dynamic`、`linapro-demo-source`、`linapro-monitor-online`、`linapro-ops-demo-guard`、`linapro-tenant-core`下 11 个测试文件，以及本`tasks.md`记录。
- 已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/api-contract.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`；技能：`lina-feedback`、`lina-review`、`lina-e2e`。
- 规则域结论：OpenSpec 反馈已先记录根因再修复并补充验证记录；E2E 直接登录请求与必填`clientType`接口契约一致；插件目录本地规范检查通过且仅改插件自有测试资产；治理脚本使用 Node 标准库和既有执行入口，跨平台边界清晰；无生产后端、数据库、缓存、数据权限、运行时 UI 或`i18n`资源变更。
- 验证证据：复查`pnpm -C hack/tests exec tsc --noEmit`通过；复查`pnpm -C hack/tests test:validate`通过；复查静态扫描确认测试直接`auth/login`请求均包含`clientType`；复查`openspec validate auth-client-type-session-metadata --strict`通过；复查`git diff --check`通过。严重问题 0，警告 0。剩余风险：未本地运行完整 Playwright E2E 矩阵，需由 GitHub Actions 或后续手动全量 E2E 覆盖服务启动后的端到端页面行为。
