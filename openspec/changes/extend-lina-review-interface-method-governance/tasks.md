## 1. 规范与审查规则

- [x] 1.1 更新 `AGENTS.md` 服务层接口规范，要求后端接口方法定义唯一且语义清晰。
- [x] 1.2 更新 `lina-review` 技能，新增后端接口方法定义审查规则。
- [x] 1.3 增加 `backend-conformance` 增量规范，覆盖重复方法、近义方法、歧义方法和兼容期说明。

## 2. 验证与审查

- [x] 2.1 运行 `openspec validate extend-lina-review-interface-method-governance --strict`。
- [x] 2.2 运行静态检索，确认 `AGENTS.md` 与 `lina-review` 均包含接口方法定义治理规则。
- [x] 2.3 记录 i18n、缓存一致性、数据权限、开发工具跨平台和测试影响评估，并执行 `lina-review` 审查。

## 3. 执行记录

- i18n 影响评估：本变更只调整治理规范和审查技能说明，不新增、修改或删除运行时用户可见文案、菜单、路由、API 文档源文本、插件清单或语言包资源。
- 缓存一致性影响评估：本变更不新增或修改生产缓存、缓存键、失效路径、跨实例协调或运行时状态。
- 数据权限影响评估：本变更不新增或修改数据操作接口、服务数据访问路径、插件 host service 数据访问或权限边界。
- 开发工具跨平台影响评估：本变更不新增或修改开发工具、脚本、CI 命令实现或默认测试入口。
- 测试策略：项目治理类反馈不新增单元测试或 E2E；使用 OpenSpec 严格校验、静态检索和审查结论作为验证证据。
- Review：已按 `lina-review` 口径完成审查。审查范围限定为 `AGENTS.md`、`.agents/skills/lina-review/SKILL.md` 和 `openspec/changes/extend-lina-review-interface-method-governance`；未发现阻塞问题。该变更不涉及 Go 生产代码、REST API、后端数据权限、运行时缓存、前端 UI 或 i18n 资源，Go 编译门禁和 E2E 不适用。

## 4. Feedback

- [x] **FB-1**: 补充暴露给源码插件的 `IsEnabledAuthoritative` 接口方法详细注释，并增加中文用途说明。
- [x] **FB-2**: 移除插件状态契约中与 `PluginStateService` 重复的 `EnablementReader` 接口定义。
- [x] **FB-3**: 将 `pluginstate` adapter 启用状态方法注释调整为英文在上、中文在下的双语格式。

## 5. Feedback Execution Record

- FB-1 影响分析：修改源码插件可见的插件状态契约、host service reader 边界、pluginstate adapter 注释和核心插件服务接口注释；不修改方法签名、运行时逻辑、REST API、数据库访问或用户可见前端行为。
- FB-1 i18n 影响评估：本次仅新增 Go 文档注释中的接口用途说明，不新增、修改或删除运行时语言包、插件 manifest i18n、apidoc i18n 或前端 UI 文案。
- FB-1 缓存一致性影响评估：本次不修改缓存键、快照刷新、失效路径、跨实例同步或运行时启用状态读取逻辑；注释明确该方法绕过进程内快照读取持久化治理状态。
- FB-1 数据权限影响评估：本次不新增或修改数据操作接口、查询条件、写操作、下载、聚合统计或插件数据访问路径。
- FB-1 开发工具跨平台影响评估：本次不新增或修改开发工具、脚本、CI 命令实现或默认测试入口。
- FB-1 验证：已运行 `cd apps/lina-core && go test ./pkg/pluginservice/contract ./pkg/pluginservice/pluginstate ./internal/service/pluginhostservices ./internal/service/plugin -count=1`，通过；已运行 `openspec validate extend-lina-review-interface-method-governance --strict`，通过；已运行静态检索确认 `IsEnabledAuthoritative` 公开契约包含详细说明和中文用途说明。
- FB-1 Review：已按 `lina-review` 口径完成审查。审查范围限定为本次修改的插件状态接口注释、adapter 注释、host service reader 注释、核心插件服务接口注释和本变更任务记录；未发现阻塞问题。该反馈不修改运行时行为、REST API、数据权限、缓存实现、前端 UI 或 i18n 资源；Go 编译门禁、OpenSpec 校验和静态检索均已通过。
- FB-2 影响分析：移除 `pkg/pluginservice/contract` 中与 `PluginStateService` 方法集合重复的 `EnablementReader`，并将 `pkg/pluginservice/pluginstate` adapter 与 `internal/service/pluginhostservices.New` 的插件状态入参统一为 `contract.PluginStateService`；不改变运行时启用状态判断逻辑、REST API、数据库访问或前端行为。
- FB-2 i18n 影响评估：本次仅收敛 Go 接口定义和入参类型，不新增、修改或删除运行时语言包、插件 manifest i18n、apidoc i18n 或前端 UI 文案。
- FB-2 缓存一致性影响评估：本次不修改插件启用状态缓存、快照刷新、权威读取、失效路径、跨实例同步或缓存键；仅移除重复接口定义，保留原有 adapter 的空值保护和 pluginID 归一化行为。
- FB-2 数据权限影响评估：本次不新增或修改数据操作接口、查询条件、详情读取、写操作、下载、聚合统计或插件数据访问路径。
- FB-2 开发工具跨平台影响评估：本次不新增或修改开发工具、脚本、CI 命令实现或默认测试入口。
- FB-2 验证：已运行 `cd apps/lina-core && go test ./pkg/pluginservice/contract ./pkg/pluginservice/pluginstate ./internal/service/pluginhostservices ./internal/cmd -count=1`，通过；已运行 `cd apps/lina-core && go test ./internal/service/plugin -run '^$' -count=1`，通过编译烟测；已运行 `openspec validate extend-lina-review-interface-method-governance --strict`，通过；已运行静态检索确认 `apps/lina-core` 目标范围内无 `contract.EnablementReader` 或 `PluginStateReader` 残留。完整 `go test ./internal/service/plugin -count=1` 因当前测试数据库缺少 `plugin_linapro_tenant_core_user_membership` 表失败，失败点与本次接口收敛无关，已用该包 `-run '^$'` 编译烟测补充覆盖。
- FB-2 Review：已按 `lina-review` 口径完成审查。审查范围限定为插件状态公共契约、pluginstate adapter、host service directory 构造入参和本变更任务记录；未发现阻塞问题。该反馈不修改运行时行为、REST API、数据权限、缓存实现、前端 UI 或 i18n 资源；Go 变更已通过相关包测试、启动绑定包测试和插件服务包编译烟测。
- FB-3 影响分析：仅调整 `apps/lina-core/pkg/pluginservice/pluginstate/pluginstate_enablement.go` 中 `IsEnabled` 与 `IsEnabledAuthoritative` 的方法注释顺序和英文说明，不修改方法签名、运行时逻辑、REST API、数据库访问或前端行为。
- FB-3 i18n 影响评估：本次只修改 Go 文档注释，不新增、修改或删除运行时语言包、插件 manifest i18n、apidoc i18n 或前端 UI 文案。
- FB-3 缓存一致性影响评估：本次不修改插件启用状态缓存、权威读取、快照刷新、失效路径、跨实例同步或缓存键；注释继续区分普通 snapshot 查询与权威状态查询语义。
- FB-3 数据权限影响评估：本次不新增或修改数据操作接口、查询条件、详情读取、写操作、下载、聚合统计或插件数据访问路径。
- FB-3 开发工具跨平台影响评估：本次不新增或修改开发工具、脚本、CI 命令实现或默认测试入口。
- FB-3 验证：已运行 `cd apps/lina-core && go test ./pkg/pluginservice/pluginstate -count=1`，通过；已运行 `openspec validate extend-lina-review-interface-method-governance --strict`，通过；已运行 `git diff --check -- apps/lina-core/pkg/pluginservice/pluginstate/pluginstate_enablement.go openspec/changes/extend-lina-review-interface-method-governance/tasks.md`，通过。
- FB-3 Review：已按 `lina-review` 口径完成审查。审查范围限定为 `pluginstate_enablement.go` 方法注释和本反馈任务记录；未发现阻塞问题。该反馈不修改运行时行为、REST API、数据权限、缓存实现、前端 UI、i18n 资源或开发工具脚本；Go 编译门禁、OpenSpec 校验和 diff 空白检查均已通过。
