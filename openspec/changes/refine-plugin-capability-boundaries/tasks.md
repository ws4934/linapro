## 1. 现状审计与目标边界确认

- [ ] 1.1 审计`pluginbridge`、`pluginhost`、`pluginservice`、`plugindb`、`pluginfs`、`sourceupgrade`、`orgcap`、`tenantcap`的导出类型、导入方和生产调用路径，形成迁移清单。
- [ ] 1.2 标记需要保留为稳定公共契约的入口、需要迁移到`internal`的实现入口、需要删除或封存的旧入口。
- [ ] 1.3 记录影响分析：`i18n`、缓存一致性、数据权限、开发工具跨平台、测试策略、SQL 和 API 契约影响。

## 2. Framework Capability 基础组件

- [ ] 2.1 新建`pkg/frameworkcap`基础结构，由根包定义 capability ID、版本、DTO、消费 service、`Org()`、`Tenant()`和 provider factory facade，`frameworkcap/internal`定义 provider env、provider state、fallback、registry 和激活实现。
- [ ] 2.2 实现`frameworkcap/internal`中的 capability manager，支持 provider factory 注册、按插件生命周期激活、撤销、冲突检测、fallback 和能力可用性查询。
- [ ] 2.3 将 provider 激活状态绑定插件 runtime revision 或等价事件，补齐单机和集群模式下的刷新、重建和故障降级策略。
- [ ] 2.4 为 framework capability manager 增加单元测试，覆盖 provider 激活、撤销、重复 provider 冲突、fallback、可选能力降级和既有插件硬依赖阻断。

## 3. 迁移 org 与 tenant 能力

- [ ] 3.1 将`pkg/orgcap`公开消费契约迁移到`pkg/frameworkcap`根包，并通过`frameworkcap.Org()`暴露组织能力实例；将 fallback、delegation、provider registry 和激活实现迁移到`pkg/frameworkcap/internal`。
- [ ] 3.2 将`pkg/tenantcap`公开消费契约迁移到`pkg/frameworkcap`根包，并通过`frameworkcap.Tenant()`暴露租户能力实例；将 fallback、delegation、provider registry 和激活实现迁移到`pkg/frameworkcap/internal`。
- [ ] 3.3 移除路由注册期间直接写全局 provider 的实现，改为插件入口声明 provider factory，由 framework capability manager 激活。
- [ ] 3.4 调整宿主认证、会话、数据权限、用户管理、插件 host service 等调用方，改为通过`frameworkcap.Org()`和`frameworkcap.Tenant()`消费能力。
- [ ] 3.5 增加或更新测试，覆盖 org/tenant provider 可用、不可用、插件禁用、插件升级和依赖不满足时的降级或阻断行为。

## 4. Provider Adapter 内部封装

- [ ] 4.1 将官方 provider 插件的 adapter 移入`backend/internal/provider/<capability>adapter/`，业务实现保留在`backend/internal/service/`。
- [ ] 4.2 调整`backend/plugin.go`，只声明 provider factory、路由、hook、cron 和生命周期，不承载 provider 业务实现。
- [ ] 4.3 检查官方插件中是否存在公开 provider 包、跨插件内部导入或直接依赖其他插件内部 service 的调用，并迁移到 framework capability 消费路径。
- [ ] 4.4 增加静态治理验证，拒绝官方插件新增公开 provider adapter 目录和跨插件`backend/internal/**`生产导入。

## 5. 统一 pluginservice 能力目录

- [ ] 5.1 将源码插件 host services 目录语义统一迁移到`pluginservice.Services`，并让`pluginhost` registrar 只传递该统一能力目录。
- [ ] 5.2 在`pluginservice`下补齐 framework capability 消费入口，例如`Services.Framework().Org()`和`Services.Framework().Tenant()`或等价明确子目录。
- [ ] 5.3 为动态插件实现`pluginservice/guest`能力 client，使 guest SDK 通过`pluginbridge`transport 调用同一`pluginservice`语义。
- [ ] 5.4 调整动态插件 host service registry，将 framework capability、config、manifest、data、cache、lock、notify 等 handler 委托到同一`pluginservice.Services`运行期实例。
- [ ] 5.5 增加源码插件和动态插件的同语义测试，确认同一能力在两类插件中的授权、错误、DTO 和降级行为一致。

## 6. 收敛插件相关组件公开面

- [ ] 6.1 收敛`pluginhost`为源码插件贡献 API，移除或迁移其中的宿主能力目录实现。
- [ ] 6.2 收敛`pluginbridge`为动态插件 ABI 和 transport facade，确保低层 codec、artifact、hostcall 和 hostservice 实现位于`pluginbridge/internal/**`。
- [ ] 6.3 收敛`plugindb`为 guest 侧受限 DSL 和必要 facade，确保 typed plan、host DB wrapper、执行器和审计上下文位于`plugindb/internal/**`或宿主内部边界。
- [ ] 6.4 将`pluginfs`中仅服务宿主扫描、路径治理、资源索引和 artifact 视图的实现迁移到职责明确的宿主`internal`组件。
- [ ] 6.5 将`sourceupgrade`中源码插件升级发现、对比和执行实现迁移到插件运行时升级治理内部组件，对外只保留稳定管理 API 和状态契约。

## 7. 既有插件依赖与生命周期治理

- [ ] 7.1 复用既有插件 manifest parser 和校验逻辑，禁止新增顶层`capabilities`或`dependencies.capabilities`等 framework capability 依赖配置块。
- [ ] 7.2 将`dependencies.plugins`依赖项收敛为`id`和可选`version`，移除`required`、`install`及等价软依赖、自动安装策略字段。
- [ ] 7.3 在插件安装、启用、禁用、卸载、升级、同版本刷新和发布切换路径中复用`dependencies.plugins`检查和下游硬依赖保护。
- [ ] 7.4 调整插件管理 API 和运行时状态投影，返回 provider 插件依赖检查结果、provider 冲突状态和能力不可用原因，但不返回自动安装计划、自动安装结果或软依赖提示。
- [ ] 7.5 若新增或修改 API DTO、错误码或文档源文本，按`api-contract`和`i18n`规则补齐文档标签、结构化错误和翻译资源。

## 8. 数据权限、性能与缓存治理

- [ ] 8.1 审查 framework capability 消费 service 的数据边界，确保普通插件消费面只暴露 DTO、批量接口和投影接口，不暴露`*gdb.Model`或内部查询对象。
- [ ] 8.2 为 org/tenant 等能力补齐批量方法或投影方法，避免列表、树形、候选项和聚合路径通过逐项调用产生`N+1`。
- [ ] 8.3 确认涉及组织、租户、用户、候选项和数据范围的 capability service 在读取阶段接入数据权限或明确受信边界。
- [ ] 8.4 验证 provider 状态、enabled snapshot、runtime revision 和 framework capability 缓存失效在单机和集群模式下具备一致性策略。

## 9. 治理扫描与验证

- [ ] 9.1 增加或扩展 Go 治理扫描，检查旧`pkg/orgcap`、`pkg/tenantcap`、公开 provider adapter、跨插件 internal 导入、旧`sourceupgrade`公共入口和低层 bridge/plugindb 公开实现导入。
- [ ] 9.2 运行覆盖变更包的 Go 编译门禁，至少包含`apps/lina-core`相关 pkg、插件运行时、宿主启动绑定包和受影响官方插件包。
- [ ] 9.3 运行动态插件 WASM 构建 smoke，确认 guest 侧不引入 host-side 数据库、provider 或宿主内部实现依赖。
- [ ] 9.4 运行或补充单元测试，覆盖 provider factory、可选能力降级、`dependencies.plugins`硬阻断、provider 冲突和跨节点失效策略。
- [ ] 9.5 执行`openspec validate refine-plugin-capability-boundaries --strict`和文档格式检查。
- [ ] 9.6 完成实现后执行`lina-review`，并按审查结论修复阻断问题。

## Feedback

- [x] **FB-1**: 明确`orgcap`、`tenantcap`迁移到`frameworkcap`后的公开契约与内部实现边界，避免把 registry、fallback 和激活实现继续作为公共组件暴露。
- [x] **FB-2**: 移除独立顶层`capabilities`配置设计，避免在插件清单中新增与`dependencies`并列的配置块。
- [x] **FB-3**: 将 framework capability 的用户使用面收敛到`frameworkcap`根组件，通过`frameworkcap.Org()`和`frameworkcap.Tenant()`获取能力实例，避免新增`frameworkcap/org`和`frameworkcap/tenant`公开子组件。
- [x] **FB-4**: 移除`dependencies.capabilities`设计，不新增任何 framework capability 依赖配置块；硬依赖复用既有`dependencies.plugins`和版本约束，可选能力通过`frameworkcap`运行时可用性降级。
- [x] **FB-5**: 删除`dependencies.plugins`中的`required`和`install`策略字段，插件依赖项只保留`id`和可选`version`；声明即硬依赖，自动安装和软依赖不由插件清单表达。

## Feedback Execution Record

- FB-1 根因：原提案把`orgcap`、`tenantcap`“迁移到`frameworkcap`”描述得过宽，容易把 provider registry、fallback、delegation 和激活 manager 继续作为公开实现面。已修正为`frameworkcap`根包只保留 capability ID、版本、DTO、消费 service、`Org()`、`Tenant()`、provider factory facade 和必要错误类型，内部实现统一进入`frameworkcap/internal`。
- FB-2 根因：原提案新增独立`capabilities`配置块，与现有插件`dependencies`能力重叠，增加 manifest 结构和维护复杂度。已先修正为不新增并列顶层`capabilities`配置；后续 FB-4 进一步移除了`dependencies.capabilities`子块。
- FB-3 根因：即使`frameworkcap/org`、`frameworkcap/tenant`只承载公开契约，也会让用户理解为多个能力子组件，增加使用和维护复杂度。已修正为`frameworkcap`根包统一暴露`Org()`、`Tenant()`和 provider factory facade，能力内部实现可在`frameworkcap/internal`下按职责拆分但不作为公共组件暴露。
- FB-4 根因：上一版虽然移除了顶层`capabilities`，但仍在`dependencies`下新增`dependencies.capabilities`，本质上继续引入第二套能力依赖模型，和既有插件依赖、版本约束及`frameworkcap`运行时可用性职责重叠。已修正为不新增任何 framework capability 依赖配置块：硬阻断、安装顺序和版本约束继续使用`dependencies.plugins`依赖 provider 插件；可选能力只通过`frameworkcap.Org()`、`frameworkcap.Tenant()`等消费 service 的`Available(ctx)`或等价状态降级；动态插件仍通过`hostServices`声明调用授权。
- FB-5 根因：旧`dependencies.plugins`中的`required`和`install`使插件清单承担软依赖、手动安装和自动安装编排语义，增加理解成本，并与“可选能力通过运行时可用性降级、安装编排由治理流程显式处理”的目标冲突。已修正为插件依赖项只保留`id`和可选`version`；声明即硬依赖；清单校验拒绝`required`、`install`及等价策略字段；插件安装接口不再返回清单驱动的自动安装计划、自动安装结果或软依赖提示。
- 影响分析：本轮仅修改 OpenSpec 文档和规范，不修改 Go 代码、HTTP API、SQL、前端 UI、插件清单实际文件或运行时文案；`i18n`资源无影响，数据权限无新增数据操作影响，缓存一致性无新增运行时缓存实现影响，开发工具跨平台无新增脚本或工具入口影响。
- 验证方式：本轮属于项目治理类反馈，使用`openspec validate refine-plugin-capability-boundaries --strict`和`git diff --check -- openspec/changes/refine-plugin-capability-boundaries`作为治理验证。
