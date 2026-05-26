## Why

当前插件相关组件的职责边界不够收敛：`pluginbridge`、`pluginhost`、`pluginservice`、`plugindb`、`pluginfs`、`sourceupgrade`以及`orgcap`、`tenantcap`等框架能力分散承担协议、能力暴露、实现注册、资源扫描和治理职责，导致源码插件与动态插件的能力使用路径不一致，也让不应公开的实现细节长期留在公共包中。

本变更通过统一插件能力目录、框架能力注册治理和内部实现封装，明确插件“提供实现”和“消费能力”的标准路径，降低长期维护复杂度。

## What Changes

- **BREAKING**：将插件能力消费统一收敛到`pkg/pluginservice`，源码插件和动态插件都通过同一组版本化能力契约访问宿主能力；`pluginbridge`只保留动态插件 ABI、WASM transport、guest bridge 和协议 facade，不再作为业务宿主能力语义的 owner。
- **BREAKING**：删除`pkg/frameworkcap`聚合包、旧`pkg/orgcap`/`pkg/tenantcap`兼容包以及`internal/service/orgcap`/`internal/service/tenantcap`双重适配层；组织和租户能力分别由`pkg/pluginservice/orgcap`、`pkg/pluginservice/tenantcap`独立组件维护，消费方通过`pluginservice.Services.Org()`、`pluginservice.Services.Tenant()`或显式注入的`orgcap.Service`、`tenantcap.Service`访问能力。
- **BREAKING**：插件提供的 provider 实现不得作为插件公开包暴露；官方源码插件的 provider adapter 默认放在`backend/internal/provider/<capability>adapter/`，真实业务实现继续放在`backend/internal/service/`。
- **BREAKING**：插件之间不得直接 import 其他插件的 provider、service、dao、entity 或任何`internal`实现；消费方插件只能通过`pluginservice.Services.Org()`或`pluginservice.Services.Tenant()`或动态插件等价 guest client 消费稳定能力。
- 新增 pluginservice capability provider factory 注册模型，由插件入口通过`orgcap.Provide(...)`或`tenantcap.Provide(...)`声明 provider factory；消费 service 使用 provider 时通过`PluginStateService.IsProviderEnabled(ctx, pluginID)`读取平台级插件可用性并懒加载 provider，禁止 provider 在路由注册期间直接写入全局注册表。
- 不新增顶层`capabilities`配置块，也不在`dependencies`下新增 pluginservice capability 依赖子块；硬依赖继续使用既有`dependencies.plugins`声明具体 provider 插件和版本约束，依赖项只保留`id`与`version`，不再支持`required`、`install`等软依赖或自动安装策略字段；可选能力通过`pluginservice`运行时可用性检查和降级表达。
- 收敛`pluginhost`职责为源码插件贡献入口，仅负责路由、hook、cron、生命周期和 provider factory 声明，不再拥有宿主能力目录实现。
- 收敛`plugindb`职责为动态插件 guest 侧受限数据 DSL 和根包 facade；typed plan、host DB、执行器和审计上下文等宿主实现细节必须留在`internal`或`pluginservice/data`内部边界。
- 收敛`pluginfs`职责为宿主内部插件资源、路径和 artifact 视图治理；除非存在稳定公共契约，扫描器、路径校验、资源索引和 cache 实现迁移到`internal`职责包。
- 收敛`sourceupgrade`职责为源码插件升级治理的内部实现；对外能力通过插件运行时升级治理和公开管理接口表达，不再把开发期源码升级实现作为公共 pkg 能力。
- 增加导入边界、公开面和目录结构治理验证，阻止低层 bridge、host service、provider adapter、resource scanner、plugin DB host executor 等实现包重新成为公共 API。

## Capabilities

### New Capabilities

- `framework-capability-registry`: 定义由`pkg/pluginservice/orgcap`与`pkg/pluginservice/tenantcap`承载的框架能力公开契约、provider factory facade、provider 懒加载、消费服务、运行时可用性和插件间能力消费治理。
- `plugin-capability-boundary-governance`: 定义`pluginhost`、`pluginservice`、`pluginbridge`、`plugindb`、`pluginfs`、`sourceupgrade`等插件相关组件的职责边界、公开面和内部实现封装规则。

### Modified Capabilities

- `core-host-boundary-governance`: 明确宿主通过`pluginservice`暴露可选框架能力，不直接依赖插件内部实现或业务存储。
- `plugin-host-service-extension`: 将源码插件与动态插件的宿主能力消费统一到`pluginservice`能力目录，并要求动态 host service 协议作为该目录的 transport 适配。
- `pluginbridge-subcomponent-architecture`: 将`pluginbridge`定位从“公开低层子组件集合”收敛为动态插件 ABI 与 transport facade，低层 codec、artifact、hostcall 和 hostservice 实现默认进入`internal`。
- `plugin-data-service`: 明确`plugindb`只暴露 guest 侧受限 DSL 和必要 facade，host-side plan、DB wrapper、执行器与审计实现不得作为公共 API。
- `plugin-dependency-management`: 明确 pluginservice capability 消费不新增依赖配置块；硬依赖复用既有`dependencies.plugins`的`id`和`version`约束，可选能力通过运行时可用性降级，并移除`required`、`install`等依赖策略字段。
- `plugin-upgrade-governance`: 要求插件升级、启用、禁用和刷新时重新校验 pluginservice capability provider、既有下游插件依赖和能力可用性，并刷新插件 enabled snapshot 与运行时派生状态。
- `source-upgrade-governance`: 明确旧`sourceupgrade`实现不再作为公共组件边界，源码插件升级治理归属插件运行时升级内部实现。

## Impact

- 影响宿主公共契约：`apps/lina-core/pkg/pluginservice`、新增`pkg/pluginservice/orgcap`与`pkg/pluginservice/tenantcap`公开契约入口、删除旧`pkg/frameworkcap`、`pkg/orgcap`、`pkg/tenantcap`入口，并继续收敛`pkg/pluginhost`、`pkg/pluginbridge`、`pkg/plugindb`公开面。
- 影响宿主内部实现：插件运行时、源码插件 registrar、动态插件 host service registry、插件资源扫描、插件数据库 host executor、源码插件升级治理和 provider 可用性判断。
- 影响官方插件结构：提供框架能力的插件需要把 provider adapter 放入`backend/internal/provider/<capability>adapter/`，业务实现保持在`backend/internal/service/`，插件入口只声明 provider factory。
- 影响动态插件 manifest：动态插件需要通过`hostServices`声明访问`pluginservice`消费服务的授权边界；如需要硬依赖某个 provider 插件，继续使用既有`dependencies.plugins`和版本约束，不能绕过`pluginservice`或直接依赖其他插件实现。
- 数据权限影响：能力消费服务涉及用户、组织、租户、数据范围、候选项或批量信息时，必须在契约中定义数据权限边界，并禁止向插件泄漏底层`DAO`、`DO`、`Entity`、`*gdb.Model`或内部查询状态。
- 缓存一致性影响：插件 enabled snapshot、runtime revision 和 provider 可用性判断属于关键运行时状态，集群模式下必须通过既有插件运行时修订、事件或共享缓存保持一致。
- `i18n`影响：本提案不新增运行时 UI 文案；后续实现如新增错误码、接口文档源文本或插件清单文案，必须按宿主或插件`i18n`启用状态维护英文 fallback 和翻译资源。
- 开发工具影响：可能需要新增或扩展治理扫描，默认使用 Go 工具或`linactl`入口，避免新增平台专属脚本。
