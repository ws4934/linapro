## 1. 现状审计与目标边界确认

- [x] 1.1 审计`pluginbridge`、`pluginhost`、`pluginservice`、`plugindb`、`pluginfs`、`sourceupgrade`、`orgcap`、`tenantcap`的导出类型、导入方和生产调用路径，形成迁移清单。
- [x] 1.2 标记需要保留为稳定公共契约的入口、需要迁移到`internal`的实现入口、需要删除或封存的旧入口。
- [x] 1.3 记录影响分析：`i18n`、缓存一致性、数据权限、开发工具跨平台、测试策略、SQL 和 API 契约影响。

## 2. Pluginservice Capability 基础组件

- [x] 2.1 新建`pkg/pluginservice/orgcap`与`pkg/pluginservice/tenantcap`独立能力组件，由各组件定义 capability ID、版本、DTO、消费`Service`、fallback/delegation 和 provider factory facade；`pkg/pluginservice/internal/capabilityregistry`只承载共享 provider factory、使用时状态判断、懒加载实例和冲突治理实现。
- [x] 2.2 实现`pkg/pluginservice/internal/capabilityregistry`中的 capability manager，支持 provider factory 注册、使用时插件启用判断、懒加载 provider、冲突检测、fallback 和能力可用性查询。
- [x] 2.3 将 provider 可用性绑定插件 enabled snapshot 与 runtime cache revision，补齐单机和集群模式下的状态刷新、重建和故障降级策略。
- [x] 2.4 为`capabilityregistry`manager 及`orgcap`、`tenantcap`组件增加单元测试，覆盖 provider 声明、插件启用/禁用、重复 provider 冲突、fallback、可选能力降级和既有插件硬依赖阻断。

## 3. 迁移 org 与 tenant 能力

- [x] 3.1 将旧`pkg/orgcap`公开消费契约迁移到`pkg/pluginservice/orgcap`独立组件，并通过`orgcap.New(...)`或注入的`orgcap.Service`暴露组织能力实例；删除旧`pkg/orgcap`兼容包，将共享 provider registry 和状态判断实现迁移到`pkg/pluginservice/internal/capabilityregistry`。
- [x] 3.2 将旧`pkg/tenantcap`公开消费契约迁移到`pkg/pluginservice/tenantcap`独立组件，并通过`tenantcap.New(...)`或注入的`tenantcap.Service`暴露租户能力实例；删除旧`pkg/tenantcap`兼容包，将共享 provider registry 和状态判断实现迁移到`pkg/pluginservice/internal/capabilityregistry`。
- [x] 3.3 移除路由注册期间直接写全局 provider 的实现，改为插件入口声明 provider factory，由 pluginservice capability manager 在使用时判断可用性。
- [x] 3.4 调整宿主认证、会话、数据权限、用户管理、插件 host service 等调用方，改为通过`orgcap.New(...)`、`tenantcap.New(...)`或显式注入的`orgcap.Service`、`tenantcap.Service`消费能力。
- [x] 3.5 增加或更新测试，覆盖 org/tenant provider 可用、不可用、插件禁用、插件升级和依赖不满足时的降级或阻断行为。

## 4. Provider Adapter 内部封装

- [x] 4.1 将官方 provider 插件的 adapter 移入`backend/internal/provider/<capability>adapter/`，业务实现保留在`backend/internal/service/`。
- [x] 4.2 调整`backend/plugin.go`，只声明 provider factory、路由、hook、cron 和生命周期，不承载 provider 业务实现。
- [x] 4.3 检查官方插件中是否存在公开 provider 包、跨插件内部导入或直接依赖其他插件内部 service 的调用，并迁移到 pluginservice capability 消费路径。
- [x] 4.4 使用静态检索、Go 编译门禁和审查记录确认生产代码未导入公开 provider adapter 或跨插件`backend/internal/**`实现。

## 5. 统一 pluginservice 能力目录

- [x] 5.1 将源码插件 host services 目录语义统一迁移到`pluginservice.Services`，并让`pluginhost` registrar 只传递该统一能力目录。
- [x] 5.2 在`pluginservice.Services`目录下直接暴露`Org()`和`Tenant()`能力入口，分别返回`pkg/pluginservice/orgcap.Service`与`pkg/pluginservice/tenantcap.Service`，不再暴露`Framework()`聚合入口。
- [x] 5.3 为动态插件实现`pluginservice/guest`能力 client，使 guest SDK 通过`pluginbridge`transport 调用同一`pluginservice`语义。
- [x] 5.4 调整动态插件 host service registry，将 pluginservice capability、config、manifest、data、cache、lock、notify 等 handler 委托到同一`pluginservice.Services`运行期实例。
- [x] 5.5 增加源码插件和动态插件的同语义测试，确认同一能力在两类插件中的授权、错误、DTO 和降级行为一致。

## 6. 收敛插件相关组件公开面

- [x] 6.1 收敛`pluginhost`为源码插件贡献 API，移除或迁移其中的宿主能力目录实现。
- [x] 6.2 收敛`pluginbridge`为动态插件 ABI 和 transport facade，确保低层 codec、artifact、hostcall 和 hostservice 实现位于`pluginbridge/internal/**`。
- [x] 6.3 收敛`plugindb`为 guest 侧受限 DSL 和必要 facade，确保 typed plan、host DB wrapper、执行器和审计上下文位于`plugindb/internal/**`或宿主内部边界。
- [x] 6.4 将`pluginfs`中仅服务宿主扫描、路径治理、资源索引和 artifact 视图的实现迁移到职责明确的宿主`internal`组件。
- [x] 6.5 将`sourceupgrade`中源码插件升级发现、对比和执行实现迁移到插件运行时升级治理内部组件，对外只保留稳定管理 API 和状态契约。

## 7. 既有插件依赖与生命周期治理

- [x] 7.1 复用既有插件 manifest parser 和校验逻辑，禁止新增顶层`capabilities`或`dependencies.capabilities`等 pluginservice capability 依赖配置块。
- [x] 7.2 将`dependencies.plugins`依赖项收敛为`id`和可选`version`，移除`required`、`install`及等价软依赖、自动安装策略字段。
- [x] 7.3 在插件安装、启用、禁用、卸载、升级、同版本刷新和发布切换路径中复用`dependencies.plugins`检查和下游硬依赖保护。
- [x] 7.4 调整插件管理 API 和运行时状态投影，返回 provider 插件依赖检查结果、provider 冲突状态和能力不可用原因，但不返回自动安装计划、自动安装结果或软依赖提示。
- [x] 7.5 若新增或修改 API DTO、错误码或文档源文本，按`api-contract`和`i18n`规则补齐文档标签、结构化错误和翻译资源。

## 8. 数据权限、性能与缓存治理

- [x] 8.1 审查 pluginservice capability 消费 service 的数据边界，确保普通插件消费面只暴露 DTO、批量接口和投影接口，不暴露`*gdb.Model`或内部查询对象。
- [x] 8.2 为 org/tenant 等能力补齐批量方法或投影方法，避免列表、树形、候选项和聚合路径通过逐项调用产生`N+1`。
- [x] 8.3 确认涉及组织、租户、用户、候选项和数据范围的 capability service 在读取阶段接入数据权限或明确受信边界。
- [x] 8.4 验证 provider 状态、enabled snapshot、runtime revision 和 pluginservice capability 缓存失效在单机和集群模式下具备一致性策略。

## 9. 治理验证

- [x] 9.1 通过静态检索、Go 编译门禁和审查记录检查非目标能力契约导入、公开 provider adapter 导入、跨插件 internal 导入、旧`sourceupgrade`公共入口和低层 bridge/plugindb 公开实现导入。
- [x] 9.2 运行覆盖变更包的 Go 编译门禁，至少包含`apps/lina-core`相关 pkg、插件运行时、宿主启动绑定包和受影响官方插件包。
- [x] 9.3 运行动态插件 WASM 构建 smoke，确认 guest 侧不引入 host-side 数据库、provider 或宿主内部实现依赖。
- [x] 9.4 运行或补充单元测试，覆盖 provider factory、可选能力降级、`dependencies.plugins`硬阻断、provider 冲突和跨节点失效策略。
- [x] 9.5 执行`openspec validate refine-plugin-capability-boundaries --strict`和文档格式检查。
- [x] 9.6 完成实现后执行`lina-review`，并按审查结论修复阻断问题。

## Feedback

- [x] **FB-1**: 明确`orgcap`、`tenantcap`迁移到`pluginservice`后的公开契约与内部实现边界，避免把 registry、fallback 和激活实现继续作为公共组件暴露。
- [x] **FB-2**: 移除独立顶层`capabilities`配置设计，避免在插件清单中新增与`dependencies`并列的配置块。
- [x] **FB-3**: 收敛 pluginservice capability 的使用面，先避免 registry、fallback 和激活实现外露；后续 FB-7 进一步调整为`pkg/pluginservice/orgcap`与`pkg/pluginservice/tenantcap`独立组件，不再采用`frameworkcap`或根包聚合。
- [x] **FB-4**: 移除`dependencies.capabilities`设计，不新增任何 pluginservice capability 依赖配置块；硬依赖复用既有`dependencies.plugins`和版本约束，可选能力通过`pluginservice`运行时可用性降级。
- [x] **FB-5**: 删除`dependencies.plugins`中的`required`和`install`策略字段，插件依赖项只保留`id`和可选`version`；声明即硬依赖，自动安装和软依赖不由插件清单表达。
- [x] **FB-6**: 将宿主插件资源文件系统实现收敛到`service/plugin`内部子组件，避免新增顶层`internal`通用组件。
- [x] **FB-7**: 删除`pkg/frameworkcap`聚合包以及`internal/service/tenantcap`、`internal/service/orgcap`双重适配层，同时删除旧`pkg/tenantcap`、`pkg/orgcap`兼容包，将租户和组织能力分别收敛到`pkg/pluginservice/tenantcap`与`pkg/pluginservice/orgcap`独立组件。
- [x] **FB-8**: 删除`pluginservice.ResetCapabilitiesForTest`及其内部依赖的公开 reset 测试辅助方法，避免测试便利函数污染生产导出面。
- [x] **FB-9**: 将公共组件主文件治理收紧为与 service 主文件一致的契约入口要求，禁止公共组件主文件承载具体实现逻辑。
- [x] **FB-10**: 将`pluginbridge/guest`中仅包含`InvokeHostService`转发函数的 WASI host call 文件合并回底层 transport 文件，减少无独立职责的文件拆分。
- [x] **FB-11**: 移除原导入边界扫描中目录级旧能力路径扫描，避免在新项目中沉淀历史路径防回归方法。
- [x] **FB-12**: 修正`pkg/pluginservice/guest/guest.go`公共组件主文件承载具体实现逻辑和接口注释不足的问题。
- [x] **FB-13**: 删除独立`linactl govern.boundary`指令及其导入边界扫描实现，避免维护单独的仓库导入边界扫描入口。
- [x] **FB-14**: 删除 framework capability provider 生命周期刷新链路，改为 provider 注册声明加使用时平台级插件状态判断，避免 provider active 状态与插件 enabled 状态形成双状态同步。
- [x] **FB-15**: 移除`wasmbuilder`中针对`capabilities`旧清单字段的兼容迁移判断，当前构建器只维护`hostServices`等目标模型字段。

## Feedback Execution Record

- FB-1 根因：原提案把`orgcap`、`tenantcap`“迁移到`pluginservice`”描述得过宽，容易把 provider registry、fallback、delegation 和激活 manager 继续作为公开实现面。已修正为由`pkg/pluginservice/orgcap`、`pkg/pluginservice/tenantcap`分别维护 capability ID、版本、DTO、消费`Service`、fallback/delegation 和 provider factory facade；共享 registry 与激活治理实现统一进入`pkg/pluginservice/internal/capabilityregistry`。
- FB-2 根因：原提案新增独立`capabilities`配置块，与现有插件`dependencies`能力重叠，增加 manifest 结构和维护复杂度。已先修正为不新增并列顶层`capabilities`配置；后续 FB-4 进一步移除了`dependencies.capabilities`子块。
- FB-3 根因：当时担心`pluginservice/orgcap`、`pluginservice/tenantcap`公开子组件会增加使用面，但后续 FB-7 明确指出`frameworkcap`和根包聚合会继续放大耦合。最终设计废弃`frameworkcap`聚合包和旧兼容 facade，改为保留`pluginservice.Services.Org()`、`pluginservice.Services.Tenant()`目录入口，并将能力契约分别归属到`pkg/pluginservice/orgcap`、`pkg/pluginservice/tenantcap`独立组件。
- FB-4 根因：上一版虽然移除了顶层`capabilities`，但仍在`dependencies`下新增`dependencies.capabilities`，本质上继续引入第二套能力依赖模型，和既有插件依赖、版本约束及`pluginservice`运行时可用性职责重叠。已修正为不新增任何 pluginservice capability 依赖配置块：硬阻断、安装顺序和版本约束继续使用`dependencies.plugins`依赖 provider 插件；可选能力只通过`orgcap.Service`、`tenantcap.Service`的`Available(ctx)`或等价状态降级；动态插件仍通过`hostServices`声明调用授权。
- FB-5 根因：旧`dependencies.plugins`中的`required`和`install`使插件清单承担软依赖、手动安装和自动安装编排语义，增加理解成本，并与“可选能力通过运行时可用性降级、安装编排由治理流程显式处理”的目标冲突。已修正为插件依赖项只保留`id`和可选`version`；声明即硬依赖；清单校验拒绝`required`、`install`及等价策略字段；插件安装接口不再返回清单驱动的自动安装计划、自动安装结果或软依赖提示。
- FB-6 根因：`pluginfs`宿主扫描实现迁移后落在`apps/lina-core/internal/pluginresourcefs`，形成新的顶层宿主内部通用组件；实际调用方集中在`internal/service/plugin/internal/{catalog,frontend,runtime,wasm}`，只有`i18n`为读取`plugin.yaml`常量引入该包。已将资源文件系统实现移动到`internal/service/plugin/internal/resourcefs`，作为`service/plugin`内部子组件维护；`i18n`改为使用本文件私有`sourcePluginManifestPath`常量，不再依赖插件服务内部实现。
- FB-7 根因：上一版虽然新增了`pkg/pluginservice/orgcap`、`pkg/pluginservice/tenantcap`，但仍保留或描述了`pkg/frameworkcap`聚合包、旧`pkg/orgcap`/`pkg/tenantcap`兼容包以及宿主`internal/service/orgcap`、`internal/service/tenantcap`双重适配层，导致组织和租户能力存在多套入口，调用方仍可能绕过独立能力组件。已删除上述旧目录和旧 facade，`pluginservice.Services`直接暴露`Org()`、`Tenant()`，共享 provider registry 只保留在`pkg/pluginservice/internal/capabilityregistry`，生产调用方统一改为导入`pkg/pluginservice/orgcap`与`pkg/pluginservice/tenantcap`。
- FB-8 根因：`pkg/pluginservice.ResetCapabilitiesForTest`作为生产包导出符号只服务测试聚合清理，并且其内部依赖`orgcap.ResetForTest`、`tenantcap.ResetForTest`和 registry `Reset()`会清空 provider factory 与 provider state，误用后需要重新执行注册流程才能恢复能力，违反生产导出面最小化要求。已删除这些公开 reset 辅助方法，测试改为使用唯一 provider plugin ID、生命周期撤销或显式注入的 tenant capability fake 保持隔离。
- FB-9 根因：`.agents/rules/backend-go.md`对 service 主文件已经使用“必须只保留契约入口和构造”的强制表述，但对`pkg/<component>`等公共组件公开包源文件仍使用“应尽可能只保留”的建议口径，导致公共组件主文件能否承载具体实现逻辑在审查时缺少硬性依据。已将公共组件主文件规则收紧为强制契约入口要求，和 service 主文件保持同一治理标准。
- FB-10 根因：`guest_hostcall_invoke_wasip1.go`只包含`InvokeHostService`对`invokeHostService`的一行公共转发，与底层 WASI host call transport 强耦合，没有形成独立职责，反而增加 guest SDK 调用路径的文件跳转成本。已将`InvokeHostService`合并到`guest_hostcall_wasip1.go`中，紧邻私有`invokeHostService`实现，并删除独立薄包装文件。
- FB-11 根因：原导入边界扫描的目录检查同时承载旧能力目录防回归和 provider adapter 目录治理，方法名看似通用，实际依赖固定历史路径；在本项目无兼容负担的前提下，目录级旧路径扫描会把迁移期防回归逻辑固化为长期工具能力。已删除目录级旧路径检查及测试，当时仅保留生产 Go import 边界检查；provider adapter 治理改为拒绝生产代码 import`backend/provider/**`，不再因目录存在本身失败。
- FB-12 根因：`apps/lina-core/pkg/pluginservice/guest/guest.go`作为`guest`公共组件主文件，除包说明、公开接口、默认目录构造外，还承载了目录委托方法、组织/租户能力 host-service 调用、响应解码和错误包装等具体实现逻辑；同时`Services`、`OrgService`和`TenantService`方法注释未说明输入上下文语义、零值返回和错误语义。已将目录委托迁移到`guest_directory.go`，将框架能力状态读取和解码迁移到`guest_capability.go`，`guest.go`仅保留组件契约、接口、轻量实现类型和构造入口，并补充接口方法注释。
- FB-13 根因：`linactl govern.boundary`作为独立仓库导入边界扫描入口，需要长期维护一套与架构规则并行的静态扫描规则；在用户明确要求不维护该新增指令后，继续保留注册项、命令薄入口和导入边界扫描实现会增加开发工具命令面与验证记录复杂度。将删除该指令和对应实现，并把本变更验证记录收敛为 Go 编译门禁、OpenSpec 校验、文件存在性检查和格式检查。
- FB-14 根因：`RefreshCapabilityProvidersForPlugin`把插件 enabled 状态刷新、provider active 状态和各基础能力组件激活调用集中耦合在主框架流程里；后续每新增一种基础能力都需要修改该集中刷新函数，并形成 provider active 与插件 enabled 的双状态同步风险。已删除 provider 生命周期刷新链路，改为 provider 插件只声明 factory，`orgcap`、`tenantcap`在使用时通过公开`PluginStateService.IsProviderEnabled(ctx, pluginID)`读取平台级插件启用快照并懒加载 provider；`ProviderServices(pluginID)`只作为 runtime-owned provider 构造环境，不把启用判断封装进主框架私有 service。
- FB-15 根因：`wasmbuilder`通过 YAML AST 主动识别顶层`capabilities`和`dependencies.capabilities`并输出迁移到`hostServices`的提示，这会把迁移期旧字段知识固化到长期构建器。项目无兼容负担时，目标模型应只定义、解析和校验当前`pluginManifest`、`dependencies.plugins{id,version}`与`hostServices`契约，不维护旧清单字段的兼容迁移判断。
- FB-1 至 FB-5 影响分析：仅修改 OpenSpec 文档和规范，不修改 Go 代码、HTTP API、SQL、前端 UI、插件清单实际文件或运行时文案；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增运行时缓存实现影响，开发工具跨平台无新增脚本或工具入口影响。
- FB-1 至 FB-5 验证方式：属于项目治理类反馈，使用`openspec validate refine-plugin-capability-boundaries --strict`和`git diff --check -- openspec/changes/refine-plugin-capability-boundaries`作为治理验证。
- FB-6 影响分析：本轮修改 Go 包归属和导入路径，不改变插件资源扫描、路径校验、artifact 解析、host storage 路径规范化或运行时 i18n 资源加载行为；无 HTTP API、SQL、数据库、前端 UI、插件清单、语言包资源或用户可见文案变更。`i18n`影响仅限去除对插件资源文件系统组件的代码依赖，运行时资源和缓存无变化；数据权限无新增数据操作影响；缓存一致性无新增缓存写入或失效影响；开发工具跨平台无新增脚本或平台专属入口影响；测试策略采用受影响 Go 包编译门禁、OpenSpec 严格校验和格式检查。
- FB-6 验证方式：已运行`cd apps/lina-core && go test ./internal/service/plugin/internal/resourcefs ./internal/service/plugin/internal/catalog ./internal/service/plugin/internal/frontend ./internal/service/plugin/internal/runtime ./internal/service/plugin/internal/wasm ./internal/service/i18n -count=1`、`openspec validate refine-plugin-capability-boundaries --strict`、`git diff --check -- apps/lina-core/internal/service/plugin/internal/resourcefs/resourcefs.go apps/lina-core/internal/service/plugin/internal/catalog/embedded.go apps/lina-core/internal/service/plugin/internal/catalog/manifest.go apps/lina-core/internal/service/plugin/internal/catalog/manifest_validate.go apps/lina-core/internal/service/plugin/internal/frontend/frontend_runtime.go apps/lina-core/internal/service/plugin/internal/runtime/artifact.go apps/lina-core/internal/service/plugin/internal/wasm/hostfn_service_storage.go apps/lina-core/internal/service/i18n/i18n_impl.go openspec/changes/refine-plugin-capability-boundaries/tasks.md`，并通过`test ! -d apps/lina-core/internal/pluginresourcefs`确认旧顶层内部组件目录已删除。
- FB-7 影响分析：本轮修改 Go 包归属、构造装配和测试替身，不新增 HTTP API、SQL、前端 UI、插件清单字段、语言包资源或用户可见文案；`i18n`资源无影响；数据权限无新增数据操作面，既有数据权限调用改为通过`orgcap.Service`、`tenantcap.Service`独立组件注入；缓存一致性影响限于 provider 可用性继续依赖既有 enabled snapshot 和 runtime cache revision，不新增缓存键或失效机制；开发工具跨平台无新增脚本或平台专属入口。
- FB-7 验证方式：使用 Go 编译门禁覆盖`pkg/pluginservice`、受影响宿主 service、插件 host service 和官方 provider adapter；使用文件存在性检查确认`apps/lina-core/internal/service/tenantcap`、`apps/lina-core/internal/service/orgcap`、`apps/lina-core/pkg/frameworkcap`、`apps/lina-core/pkg/tenantcap`、`apps/lina-core/pkg/orgcap`已删除；使用`openspec validate refine-plugin-capability-boundaries --strict`和`git diff --check`验证治理文档。
- FB-8 影响分析：本轮缩小`pkg/pluginservice`、`orgcap`、`tenantcap`和`capabilityregistry`生产导出面，不新增 HTTP API、SQL、前端 UI、插件清单字段、语言包资源或用户可见文案；`i18n`资源无影响；数据权限无新增数据操作面；缓存一致性影响为避免测试清空全局 provider factory，不改变生产 provider factory 注册、使用时启用判断、runtime cache revision 或 enabled snapshot 刷新策略；开发工具跨平台无新增脚本或工具入口影响；测试策略采用受影响 Go 包编译门禁、静态检索和 OpenSpec 严格校验。
- FB-8 验证方式：已运行`cd apps/lina-core && go test ./pkg/pluginservice ./pkg/pluginservice/orgcap ./pkg/pluginservice/tenantcap ./pkg/pluginservice/internal/capabilityregistry ./internal/service/auth ./internal/service/notify ./internal/service/menu ./internal/service/role ./internal/service/user ./internal/service/plugin ./internal/service/plugin/internal/wasm -count=1`、静态检索确认`ResetCapabilitiesForTest`/`ResetForTest`/registry `Reset()`无残留、`openspec validate refine-plugin-capability-boundaries --strict`和`git diff --check -- apps/lina-core/pkg/pluginservice/capability_manager.go apps/lina-core/pkg/pluginservice/orgcap/orgcap_manager.go apps/lina-core/pkg/pluginservice/tenantcap/tenantcap_manager.go apps/lina-core/pkg/pluginservice/internal/capabilityregistry/manager.go apps/lina-core/internal/service/auth/auth_tenant_flow_test.go apps/lina-core/internal/service/notify/notify_send_tenant_test.go apps/lina-core/internal/service/menu/menu_platform_guard_test.go apps/lina-core/internal/service/role/role_tenant_boundary_test.go apps/lina-core/internal/service/user/user_tenant_membership_test.go apps/lina-core/internal/service/plugin/plugin_platform_guard_test.go apps/lina-core/internal/service/plugin/plugin_capability_revision_test.go apps/lina-core/internal/service/plugin/plugin_startup_consistency_test.go apps/lina-core/internal/service/plugin/internal/wasm/hostfn_service_framework_test.go openspec/changes/refine-plugin-capability-boundaries/tasks.md`。
- FB-9 影响分析：本轮只修改后端 Go 治理规则和 OpenSpec 反馈记录，不修改生产 Go 代码、HTTP API、SQL、前端 UI、插件清单、语言包资源、运行时文案、缓存实现、数据访问路径或开发工具入口；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增运行时缓存或失效影响，开发工具跨平台无新增脚本或工具影响；测试策略采用 OpenSpec 严格校验、静态检索和格式检查。
- FB-9 验证方式：已运行`openspec validate refine-plugin-capability-boundaries --strict`、`git diff --check -- .agents/rules/backend-go.md openspec/changes/refine-plugin-capability-boundaries/tasks.md`，并通过`rg`确认`.agents/rules/backend-go.md`不再保留公共组件主文件“应尽可能只保留”的建议口径，已改为强制契约入口要求。
- FB-10 影响分析：本轮只调整`pluginbridge/guest` WASI host call 文件组织，不改变`InvokeHostService`函数签名、host service envelope、opcode、错误处理、HTTP API、SQL、数据库访问、前端 UI、插件清单、语言包资源、运行时用户可见文案、缓存状态或开发工具入口；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增缓存或失效影响，开发工具跨平台无新增脚本或平台专属入口影响；测试策略采用变更包 Go 测试、WASI 交叉编译烟测、格式检查和 OpenSpec 严格校验。
- FB-10 验证方式：已运行`cd apps/lina-core && go test ./pkg/pluginbridge/guest -count=1`、`cd apps/lina-core && GOOS=wasip1 GOARCH=wasm go test -c ./pkg/pluginbridge/guest -o /tmp/pluginbridge_guest_wasip1.test`和`git diff --check -- apps/lina-core/pkg/pluginbridge/guest/guest_hostcall_wasip1.go openspec/changes/refine-plugin-capability-boundaries/tasks.md`；直接执行`GOOS=wasip1 GOARCH=wasm go test ./pkg/pluginbridge/guest`会因本机无法执行 wasm 测试二进制返回`exec format error`，因此使用`go test -c`覆盖 WASI 编译门禁。
- FB-11 影响分析：本轮修改 OpenSpec 治理描述，不修改生产后端服务、HTTP API、SQL、前端 UI、插件清单、语言包资源或运行时用户可见文案；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增缓存或失效影响。开发工具跨平台影响限于当时 Go 实现的导入边界扫描规则收敛，不新增 Shell、PowerShell 或平台专属默认入口；测试策略采用 Go 工具单元测试、全量`linactl`测试、OpenSpec 严格校验和格式检查。
- FB-11 验证方式：已运行`cd hack/tools/linactl && go test ./... -count=1`、`openspec validate refine-plugin-capability-boundaries --strict`和`git diff --check -- openspec/changes/refine-plugin-capability-boundaries/tasks.md openspec/changes/refine-plugin-capability-boundaries/specs/plugin-capability-boundary-governance/spec.md`，全部通过。
- FB-12 影响分析：本轮只调整`pkg/pluginservice/guest`文件组织和接口注释，不改变`Services`、`OrgService`、`TenantService`公开方法签名、host service method、payload、错误返回路径、HTTP API、SQL、数据库访问、前端 UI、插件清单、语言包资源或运行时用户可见文案；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增缓存或失效影响，开发工具跨平台无新增脚本或工具入口影响。已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`。测试策略采用变更包 Go 测试、OpenSpec 严格校验和格式检查。
- FB-12 验证方式：已运行`cd apps/lina-core && go test ./pkg/pluginservice/guest -count=1`、`openspec validate refine-plugin-capability-boundaries --strict`和`git diff --check -- apps/lina-core/pkg/pluginservice/guest/guest.go apps/lina-core/pkg/pluginservice/guest/guest_directory.go apps/lina-core/pkg/pluginservice/guest/guest_capability.go apps/lina-core/pkg/pluginservice/guest/guest_test.go openspec/changes/refine-plugin-capability-boundaries/tasks.md`，全部通过。
- FB-13 影响分析：本轮删除`hack/tools/linactl`中的独立导入边界扫描命令注册、命令薄入口和内部扫描实现，并调整 OpenSpec 规范和任务记录，不修改生产后端服务、HTTP API、SQL、前端 UI、插件清单、语言包资源或运行时用户可见文案；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增缓存或失效影响。开发工具跨平台影响为减少一个 Go 实现的`linactl`命令入口，不新增 Shell、PowerShell 或平台专属默认入口；测试策略采用`linactl`全量 Go 测试、OpenSpec 严格校验、残留引用检索和格式检查。已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/backend-go.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`。
- FB-13 验证方式：已运行`cd hack/tools/linactl && go test ./... -count=1`、`openspec validate refine-plugin-capability-boundaries --strict`、`git diff --check -- hack/tools/linactl/command.go hack/tools/linactl/command_govern.boundary.go hack/tools/linactl/internal/boundaryscan openspec/changes/refine-plugin-capability-boundaries/tasks.md openspec/changes/refine-plugin-capability-boundaries/specs/plugin-capability-boundary-governance/spec.md`和`rg -n "govern\\.boundary|boundaryscan|Boundary governance" openspec/changes/refine-plugin-capability-boundaries hack/tools/linactl -g '!**/node_modules/**'`；残留检索只保留 FB-13 任务和根因说明中的删除对象引用。
- FB-14 影响分析：本轮修改 Go 后端公共能力契约、插件服务启用状态读取、`orgcap`/`tenantcap`provider registry、动态插件 framework host service、data host 组织数据范围接入和相关单元测试；不新增 HTTP API、SQL、前端 UI、插件清单字段、语言包资源或用户可见文案。`i18n`资源无影响；数据权限没有新增数据操作面，data host 的部门数据范围继续通过注入的`orgcap.Service`在查询阶段约束；缓存一致性影响为复用插件 enabled snapshot 和 runtime cache revision 作为 provider 可用性权威来源，删除第二套 provider active 状态和主动刷新链路；开发工具跨平台无新增脚本、命令或平台专属入口影响。已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/data-permission.md`。
- FB-14 验证方式：已运行`cd apps/lina-core && go test ./internal/service/plugin/internal/datahost ./internal/service/plugin/internal/wasm ./internal/service/plugin/internal/integration ./internal/service/pluginhostservices ./internal/service/auth ./internal/service/notify ./internal/service/menu ./internal/service/role ./internal/service/user ./internal/service/plugin ./internal/cmd ./pkg/pluginservice/... -count=1`、`openspec validate refine-plugin-capability-boundaries --strict`、`git diff --check`，并通过`rg -n "RefreshCapabilityProvidersForPlugin|RefreshFrameworkCapability|RefreshProviderForPlugin|ProviderEnv\\{[^}]*PluginVersion|ProviderEnv\\{[^}]*RuntimeRevision|ProviderStatus\\{[^}]*RuntimeRevision" apps/lina-core apps/lina-plugins --glob '*.go'`确认旧 provider 刷新链路和 provider runtime revision 状态无残留。
- FB-15 影响分析：本轮修改`hack/tools/linactl/internal/wasmbuilder`构建器解析逻辑、对应单元测试和 OpenSpec 规范记录，不修改生产后端服务、HTTP API、SQL、数据库访问、前端 UI、运行时语言包或插件清单实际文件；`i18n`资源无影响；数据权限无新增数据操作影响；缓存一致性无新增缓存或失效影响。开发工具跨平台影响为减少 Go 构建器中的 YAML 旧字段兼容判断，不新增 Shell、PowerShell 或平台专属默认入口；测试策略采用`linactl` Go 测试、OpenSpec 严格校验、残留引用检索和格式检查。已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/plugin.md`、`.agents/rules/testing.md`、`.agents/rules/backend-go.md`、`.agents/rules/i18n.md`、`.agents/rules/architecture.md`。
- FB-15 验证方式：已运行`cd hack/tools/linactl && go test ./... -count=1`、`openspec validate refine-plugin-capability-boundaries --strict`、`git diff --check -- hack/tools/linactl/internal/wasmbuilder/wasmbuilder_manifest.go hack/tools/linactl/internal/wasmbuilder/wasmbuilder_embed.go hack/tools/linactl/internal/wasmbuilder/wasmbuilder_test.go openspec/changes/refine-plugin-capability-boundaries/specs/plugin-dependency-management/spec.md openspec/changes/refine-plugin-capability-boundaries/tasks.md`，并通过`rg -n 'rejectUnsupportedManifestDependencyFields|rejectUnsupportedDependencyFields|rejectUnsupportedPluginDependencyFields|deprecated capabilities|hostServices migration|capabilities is not supported|top-level capabilities|deprecated dependency|WasmSectionBackendCapabilities' hack/tools/linactl/internal/wasmbuilder openspec/changes/refine-plugin-capability-boundaries/specs/plugin-dependency-management/spec.md openspec/changes/refine-plugin-capability-boundaries/tasks.md`确认`wasmbuilder`中不再保留旧`capabilities`兼容迁移判断。

## Implementation Execution Record

- 现状审计与边界归属：`pluginhost`保留为源码插件贡献 API；`pluginservice.Services`成为源码插件和动态插件统一能力目录；`pluginbridge`保留动态插件 ABI、transport、guest 和协议 facade，低层 codec、artifact、hostcall、hostservice 实现位于`pluginbridge/internal/**`；`plugindb`保留 guest 侧受限 DSL，host DB、typed plan 和执行器位于`plugindb/internal/**`；`pluginfs`宿主扫描实现迁移到`internal/service/plugin/internal/resourcefs`；`sourceupgrade`公开包保留稳定 delegate facade，源码升级执行器迁移到`internal/service/plugin/internal/sourceupgrade`。
- `pluginservice`实现：新增`pkg/pluginservice/orgcap`与`pkg/pluginservice/tenantcap`独立能力组件，分别承载能力 ID、DTO、消费`Service`、fallback/delegation、`New(...)`、`Provide(...)`和 provider factory facade；`pkg/pluginservice/internal/capabilityregistry`承载共享 provider factory 注册、使用时启用判断、懒加载 provider 实例和冲突状态；旧`pkg/frameworkcap`、`pkg/orgcap`、`pkg/tenantcap`和宿主`internal/service/orgcap`、`internal/service/tenantcap`已删除。
- 官方插件封装：`linapro-org-core`、`linapro-tenant-core`的 provider adapter 已迁移到`backend/internal/provider/<capability>adapter/`；`backend/plugin.go`只声明 provider factory、路由和生命周期钩子；生产路径没有跨插件`backend/internal/**`导入。
- 统一宿主能力目录：新增`pkg/pluginservice/services.go`和`pkg/pluginservice/guest`；源码插件通过`pluginhost`registrar 接收同一个`pluginservice.Services`目录，并通过`Services.Org()`、`Services.Tenant()`访问独立能力组件；动态插件通过`pluginbridge/guest.InvokeHostService`进入同一 host service 语义；当前仅保留动态插件协议层已有`framework`服务名承载`org.status`、`tenant.status`只读状态查询，后续若要拆分动态协议名称应另立变更处理。
- 插件依赖治理：manifest parser、runtime artifact parser 和`wasmbuilder`仅维护当前依赖模型，不定义、读取或消费顶层`capabilities`、`dependencies.capabilities`等能力依赖字段；`dependencies.plugins`只保留`id`和可选`version`，声明即硬依赖；插件安装、启用、禁用、卸载、升级、同版本刷新和发布切换路径复用硬依赖检查与下游保护。
- API 契约影响：插件依赖检查和管理投影删除自动安装计划、手动安装策略和软依赖提示字段；保留结构化依赖检查、provider 状态和能力不可用原因；无新增时间字段，未引入前端逐项补查契约。
- 数据权限与性能影响：普通`orgcap.Service`、`tenantcap.Service`和`pluginservice.Services.Org()`、`pluginservice.Services.Tenant()`只暴露 DTO、状态、批量和投影方法；`*gdb.Model`与`*ghttp.Request`只保留在 provider-facing 或宿主受信迁移边界，不作为普通插件消费面；新增`ListUserTenantProjections`等批量投影，避免用户、组织、租户列表路径逐项调用。
- 缓存一致性影响：provider 可用性权威来源为插件 enabled snapshot 与 runtime cache revision；插件安装、启用、禁用、卸载、升级、同版本刷新和依赖状态变化后刷新 enabled snapshot 并发布 runtime revision，能力 registry 不再维护第二套 provider active 状态；只读依赖检查、`Available(ctx)`和 provider 使用时判断不写 registry、不发布失效事件、不清空无关`i18n`或前端缓存。
- `i18n`影响：未新增运行时 UI 文案、菜单文案、语言包或插件`manifest/i18n`资源；涉及 API 文档源文本和错误码的调整沿用已有结构化错误与文档标签，未新增需要补译的运行时语言资源；未启用多语言的插件不新增占位翻译文件。
- SQL 与数据库影响：未新增或修改宿主 SQL、插件安装 SQL、mock 数据、DAO、DO、Entity 或索引；无需执行`make init`或`make dao`；软删除和数据库时间维护无影响。
- 前端 UI 与 E2E 影响：未新增或修改前端页面、路由、表格、表单、按钮、菜单显示或用户可观察 UI 工作流；本轮不新增 E2E，采用 Go 单元测试、编译门禁、WASM 构建 smoke、OpenSpec 严格校验和审查记录覆盖。
- 开发工具跨平台影响：已删除独立导入边界扫描指令和实现；不新增 Shell、PowerShell 或平台专属默认入口；验证使用`cd hack/tools/linactl && go test ./... -count=1`。
- 验证记录：已运行`openspec validate refine-plugin-capability-boundaries --strict`、`git diff --check`、`cd hack/tools/linactl && go test ./... -count=1`、`cd apps/lina-core && go test ./pkg/pluginservice/... -count=1`、`cd apps/lina-core && go test ./internal/service/auth ./internal/service/user ./internal/service/role ./internal/service/menu ./internal/service/notify ./internal/service/plugin ./internal/service/plugin/internal/wasm ./internal/service/plugin/internal/integration ./internal/service/plugin/internal/datahost ./internal/service/pluginhostservices -run '^$' -count=1`、`cd apps/lina-core && go test ./internal/cmd -run '^$' -count=1`、`cd apps/lina-core && go test ./internal/service/plugin/internal/wasm -run 'TestHandleHostServiceInvokeFrameworkStatus|TestConfigureFrameworkHostServiceRejectsNil' -count=1`、`cd apps/lina-core && go test ./internal/service/plugin -run 'TestSourceProviderAvailabilityFollowsEnabledSnapshot|TestPluginGovernanceMethodsRejectTenantContext' -count=1`、`cd apps/lina-core && go test ./internal/service/auth -run 'TestLoginSelectTenantSwitchTenantLogoutFlow|TestLoginRejectsTenantUserWithoutActiveTenant|TestRefreshPreservesSessionOnProviderInfraError|TestRefreshRejectsAfterTenantMembershipRemoved' -count=1`、`cd apps/lina-core && go test ./internal/service/plugin/internal/datahost ./internal/service/plugin/internal/wasm ./internal/service/plugin/internal/integration ./internal/service/pluginhostservices ./internal/service/auth ./internal/service/notify ./internal/service/menu ./internal/service/role ./internal/service/user ./internal/service/plugin ./internal/cmd ./pkg/pluginservice/... -count=1`、`GOWORK=/Users/john/Workspace/github/linaproai/linapro/temp/go.work.plugins-test go test lina-plugin-linapro-tenant-core/backend/internal/service/membership lina-plugin-linapro-tenant-core/backend/internal/service/provider lina-plugin-linapro-tenant-core/backend/internal/provider/tenantadapter lina-plugin-linapro-org-core/backend/internal/provider/orgcapadapter lina-plugin-linapro-monitor-server/backend -run 'TestDoesNotExist' -count=1`。
- 静态检索记录：`rg`确认`apps/lina-plugins`的`plugin.yaml`不再包含`required:`、`install:`或`capabilities:`；旧`lina-core/pkg/frameworkcap`、`lina-core/pkg/orgcap`、`lina-core/pkg/tenantcap`、`lina-core/internal/service/orgcap`、`lina-core/internal/service/tenantcap`只剩规范、任务记录或测试 fixture；`dependencies.plugins`旧策略字段只出现在负向测试或无关 install/uninstall 语义中。
- 审查记录：已按`lina-review`读取`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/data-permission.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`、`.agents/rules/dev-tooling.md`、`.agents/rules/api-contract.md`和`.agents/rules/database.md`，审查范围来自`git status --short`、`git ls-files --others --exclude-standard`和当前 OpenSpec 变更；未发现阻塞问题。审查中移除未使用的 direct provider shim，并将残留旧聚合语义注释和错误文本收敛为`pluginservice capability`语义。
- 测试补充说明：曾并发执行`go run ./hack/tools/linactl wasm p=linapro-demo-dynamic out=temp/output`和`cd hack/tools/linactl && go test ./... -count=1`，导致同一动态插件源码目录的临时生成文件竞争；已串行重跑`cd hack/tools/linactl && go test ./... -count=1`通过。
