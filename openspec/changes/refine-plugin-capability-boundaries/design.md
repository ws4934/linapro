## Context

当前插件体系已经具备源码插件、动态插件、`pluginbridge`、`pluginhost`、`pluginservice`、`plugindb`和插件生命周期资源治理等基础，但职责边界仍然存在交叉：

- `pluginhost`是源码插件贡献入口，但当前也容易承载宿主能力目录语义，导致“插件贡献能力”和“插件消费宿主能力”混在一起。
- `pluginservice`最接近统一宿主能力目录，但动态插件仍通过`pluginbridge`host service wire 感知部分能力语义，造成源码插件和动态插件能力暴露路径不一致。
- `pluginbridge`应是动态插件 ABI、WASM transport 和协议适配层，但低层 codec、artifact、host service wire 和 guest SDK 曾经扩成较大公开面，容易被当作宿主业务能力 owner。
- `plugindb`既包含 guest 侧受限 DSL，也包含 host-side plan、DB wrapper 和审计实现；这些 host-side 实现不应成为插件开发者可依赖的公共 API。
- `pluginfs`和`sourceupgrade`更接近宿主内部资源扫描、路径治理和升级治理实现，不应作为插件公共能力直接暴露。
- `orgcap`、`tenantcap`等能力属于“宿主定义框架能力契约，插件提供具体实现”，但接口、空实现、全局 provider 注册和插件实现分散在不同组件中，理解和生命周期治理成本高。

本设计以没有历史兼容负担为前提，直接收敛到长期更清晰的模型。

## Goals / Non-Goals

**Goals:**

- 将插件能力消费统一到`pkg/pluginservice`，源码插件和动态插件使用同一能力目录和同一授权语义。
- 新增`pkg/frameworkcap`作为框架能力统一公开入口，使`org`、`tenant`等能力通过`frameworkcap.Org()`、`frameworkcap.Tenant()`等根组件方法获取，避免为每个能力新增公开子组件；provider registry、fallback、delegation 和激活 manager 位于`frameworkcap/internal`。
- 明确 provider adapter 是插件内部实现，不作为插件公开包暴露；官方插件默认使用`backend/internal/provider/<capability>adapter/`。
- 明确插件间调用必须通过`frameworkcap`消费服务或`pluginservice`能力目录，禁止插件直接依赖其他插件内部实现。
- 将`pluginhost`、`pluginbridge`、`plugindb`、`pluginfs`和`sourceupgrade`收敛到单一职责，并把非公开资源放入`internal`。
- 为既有插件依赖、启停降级、provider 激活、缓存一致性、数据权限和治理扫描建立明确规则。

**Non-Goals:**

- 不引入通用 DI 容器或全局 service locator。
- 不把`frameworkcap`设计成插件可以任意发布新业务 API 的市场；它只承载框架认可的稳定能力契约。
- 不允许动态插件直接实现包含 Go 运行时对象、`*gdb.Model`、`*ghttp.Request`或宿主 DAO 的 provider 接口。
- 不为普通插件开放跨插件直接调用、直接 import provider 包或直接访问其他插件数据库表的能力。
- 不在本变更中设计在线热替换 provider 的复杂多版本并行运行模型；provider 切换跟随插件安装、启用、禁用、升级和刷新治理。

## Decisions

### 1. 组件职责重新划分

目标职责如下：

```text
pkg/pluginhost
  source plugin contribution API only
  routes / hooks / cron / lifecycle / provider factory declaration

pkg/pluginservice
  unified plugin-facing host capability directory
  source adapters + dynamic host service handlers + guest clients

pkg/pluginbridge
  dynamic plugin ABI and transport only
  WASM envelope / host call / guest bridge / protocol facade

pkg/frameworkcap
  host-defined capability contracts with plugin-provided implementations
  root accessors: Org() / Tenant()
  root DTOs and stable service interfaces
  internal/org / internal/tenant / internal/registry / internal/fallback / internal/activation

pkg/plugindb
  dynamic guest-side restricted data DSL and root facade
  host execution details under internal or pluginservice/data

internal/service/plugin/...
  runtime / lifecycle / catalog / resourcefs / sourceupgrade / datahost
```

备选方案是继续让`pluginbridge`和`pluginservice`分别服务动态插件和源码插件。该方案短期改动少，但会让同一宿主能力拥有两套公开入口、两套授权语义和两套 SDK 叙事，因此不采用。

### 2. `pluginservice`成为统一能力目录

`pluginservice.Services`是所有插件消费宿主能力的唯一入口。源码插件通过 registrar 拿到该目录；动态插件通过`pluginservice/guest` client 发起调用，底层使用`pluginbridge`transport。

```text
source plugin
  -> pluginhost registrar context
  -> pluginservice.Services
  -> frameworkcap / config / data / cache / notify / auth / i18n

dynamic plugin
  -> pluginservice/guest client
  -> pluginbridge host service envelope
  -> pluginservice host handlers
  -> same runtime services
```

`pluginbridge`不再决定业务能力命名、授权资源语义或框架能力降级语义；它只传输`service`、`method`、payload 和结构化错误。

### 3. `frameworkcap`使用根组件聚合公开 API

每个框架能力至少包含公开消费面和内部治理面：

```go
// OrgService 由宿主和其他插件消费。
type OrgService interface {
    Available(ctx context.Context) bool
    // 面向消费方的稳定 DTO / batch / projection 方法。
}

// Org 获取组织能力消费实例。
func Org() OrgService

// Tenant 获取租户能力消费实例。
func Tenant() TenantService
```

`frameworkcap`根包是唯一公开导入入口。它直接维护插件和宿主稳定依赖所需的契约：capability ID、版本、DTO、`OrgService`、`TenantService`、provider factory 声明 facade 和必要错误类型。provider registry、fallback/delegation、provider 激活状态、冲突检测、缓存和 manager 实现必须放到`frameworkcap/internal`下；内部实现可以按`internal/org`、`internal/tenant`等子组件拆分，但这些拆分不成为用户可 import 的公共组件。

消费面不得暴露 provider 实例，也不得暴露实现插件的`internal/service`、DAO、Entity、缓存快照或`*gdb.Model`。

备选方案是让消费方直接拿 provider 接口。该方案省一层包装，但会把 provider 生命周期、实现细节和消费契约混在一起，无法为可选能力降级、fallback、缓存和审计提供统一治理，因此不采用。

### 4. Provider 通过 factory 注册，不能直接写全局注册表

提供方插件只在插件入口声明 provider factory：

```text
plugin backend/plugin.go
  -> pluginhost.RegisterSourcePlugin(...)
  -> frameworkcap.ProvideOrg(...) / frameworkcap.ProvideTenant(...)
  -> factory(ctx, ProviderEnv) returns Provider
```

`ProvideOrg(...)`、`ProvideTenant(...)`是`frameworkcap`根包中的窄 facade，用于收集 provider factory 声明；声明会进入`frameworkcap/internal`的 registry。`ProviderEnv`由宿主激活时提供，包含插件 ID、插件版本或 release 快照、`pluginservice.Services`、运行模式和必要的治理上下文。`frameworkcap/internal`中的 manager 根据插件安装、启用、依赖、版本、租户状态、升级和禁用事件激活或撤销 provider。

不采用`frameworkcap.RegisterProvider(provider)`这种直接全局注册模型，因为它无法表达插件生命周期，也难以在禁用、升级失败、同版本刷新或集群状态传播时保持一致。

### 5. Provider adapter 放在插件内部

官方插件采用以下结构：

```text
apps/lina-plugins/<plugin-id>/backend/
  plugin.go
  internal/
    service/
      <domain>/
    provider/
      <capability>adapter/
```

`backend/internal/service/`承载业务实现和领域编排；`backend/internal/provider/<capability>adapter/`只做薄适配，把内部 service 映射到对应 framework capability 的 provider-facing 契约。简单场景可以由`internal/service`中的私有类型直接实现 provider-facing 契约，但不得把 provider 实现放到插件公开包中。

这样其他插件无法 import provider adapter，只能通过`pluginservice.Services.Framework().<Capability>()`消费稳定能力。

### 6. 插件依赖复用既有`dependencies.plugins`，不新增能力依赖配置

framework capability 消费不新增 manifest 配置概念：不新增顶层`capabilities`配置块，也不在`dependencies`下新增`capabilities`子块。`dependencies`继续只表达 LinaPro 框架版本约束和插件间依赖约束。

当消费方插件必须保证某个 provider 插件存在并满足版本时，使用既有`dependencies.plugins`声明具体 provider 插件和版本范围：

```yaml
dependencies:
  plugins:
    - id: linapro-tenant-core
      version: ">=1.0.0"
```

`dependencies.plugins`只表达硬依赖，不再支持`required`、`install`或等价策略字段。声明即表示目标插件安装、启用、升级、禁用、卸载和 provider 发布切换时必须保护该依赖；依赖插件是否自动安装属于更高层安装流程或管理员操作，不由插件清单决定。

当消费方只是可选使用某项框架能力时，不声明硬 provider 插件依赖；运行时通过`frameworkcap.Org().Available(ctx)`、`frameworkcap.Tenant().Available(ctx)`或等价能力状态检查可用性，并按规范隐藏功能、返回零值或执行降级行为。

动态插件仍必须通过`hostServices`声明对应能力服务、方法和资源边界。`hostServices`表达动态 host service 调用授权，`dependencies.plugins`表达安装、启用和升级所需的具体 provider 插件依赖，两者职责不重叠。

### 7. 数据权限和性能进入 capability 契约

`org`、`tenant`这类能力经常参与用户、组织、租户、候选项、树形数据和数据范围计算。消费服务必须提供 DTO 化、批量化和投影化方法，例如批量解析用户、部门树、租户上下文、可见范围投影，而不是让插件逐项调用详情或直接修改底层查询对象。

低层查询注入能力如果确实存在，只能留给宿主内部或受信 provider 适配层，不作为普通插件消费契约。所有跨插件能力消费必须避免`N+1`、前端瀑布式调用和数据存在性泄漏。

### 8. 缓存与集群一致性绑定插件运行时修订

framework capability provider 激活状态属于关键运行时快照。插件安装、启用、禁用、卸载、升级、同版本刷新、依赖状态变化或租户可用性变化后，宿主必须刷新 provider 激活状态，并通过插件 runtime revision、事件广播、共享缓存或等价机制传播到集群节点。

如果`frameworkcap`内部缓存消费服务结果或 provider 状态，缓存键必须包含 capability ID、版本、provider 插件 ID、provider release 或 generation，以及必要租户作用域。缓存失效必须幂等、可重试、可观测。

### 9. 公开面治理作为实现门禁

本变更需要新增或扩展静态治理，至少覆盖：

- 非测试代码不得 import 其他插件的`backend/internal/**`。
- 非授权边界不得 import `pluginbridge/internal/**`、`plugindb/internal/**`或插件资源扫描内部包。
- 官方 provider adapter 不得位于`backend/provider/**`等公开目录。
- `pkg/pluginfs`和`pkg/sourceupgrade`如不承载稳定公共契约，生产使用应迁移到`internal`边界。
- `orgcap`、`tenantcap`旧公共路径迁移后不得继续作为新代码入口。

治理入口优先使用 Go 工具或`linactl`，避免平台专属 Shell 脚本。

## Risks / Trade-offs

- Provider 激活状态变成运行时关键状态，可能引入缓存一致性风险。缓解方式：绑定插件 runtime revision，集群模式通过事件或共享修订传播，并为能力可用性提供重建路径。
- 统一能力目录会触及多个组件，迁移面较大。缓解方式：按“新增目标边界、迁移官方调用、收口旧入口、增加治理扫描”的顺序实施。
- 动态插件无法直接实现 Go provider-facing 契约。缓解方式：动态插件提供能力时必须使用 DTO 化命令、查询和事件协议，由宿主 proxy adapter 转为`frameworkcap/internal`可激活的 provider 实现；本轮优先支持源码插件 provider。
- `org`、`tenant`等能力如果过度暴露底层查询能力，会形成数据权限绕过。缓解方式：普通消费面只暴露高层 DTO、批量和投影方法，低层查询注入保留在宿主内部或受信 adapter。
- 不新增能力依赖配置会让“依赖某个能力”回到“依赖某个 provider 插件”表达。缓解方式：仅在确实需要硬阻断、安装顺序或版本约束时声明`dependencies.plugins`；可选能力一律通过`frameworkcap`运行时可用性和降级表达，避免 manifest 概念膨胀。删除`required`和`install`会减少清单表达力，但也避免把插件清单变成安装编排 DSL；自动安装和批量安装应由管理流程显式处理。

## Migration Plan

1. 建立`pkg/frameworkcap`骨架，根包定义 capability ID、版本、DTO、消费 service、`Org()`、`Tenant()`和 provider factory facade；`frameworkcap/internal`定义 provider env、activation manager、fallback/delegation 和 registry。
2. 将`orgcap`、`tenantcap`公开契约迁移到`frameworkcap`根包，通过`frameworkcap.Org()`、`frameworkcap.Tenant()`、`frameworkcap.ProvideOrg(...)`和`frameworkcap.ProvideTenant(...)`暴露稳定入口，并把 fallback、provider 激活和 registry 实现放入`frameworkcap/internal`；删除或封存旧全局 provider 注册入口。
3. 将`pluginhost.HostServices`等能力目录语义迁移到`pluginservice.Services`，`pluginhost`只保留源码插件贡献 API 和 provider factory 声明。
4. 为动态插件提供`pluginservice/guest`能力 client，将动态 host service handler 接到同一`pluginservice.Services`运行期实例。
5. 调整官方 provider 插件，把 provider adapter 移到`backend/internal/provider/<capability>adapter/`，业务逻辑留在`backend/internal/service/`，`backend/plugin.go`只声明 factory。
6. 调整消费方插件和宿主模块，使其通过`pluginservice.Services.Framework().<Capability>()`、`frameworkcap.Org()`、`frameworkcap.Tenant()`或宿主内部注入的 framework capability 消费服务访问能力。
7. 收敛`pluginbridge`、`plugindb`、`pluginfs`、`sourceupgrade`公开面，把低层实现迁入`internal`或宿主内部职责包，并提供必要 facade。
8. 复用既有插件`dependencies.plugins`依赖校验，在 provider 插件启用、禁用、卸载、升级和发布切换时保护已声明的下游插件硬依赖；可选能力只通过`frameworkcap`可用性快照和降级策略治理。
9. 增加治理扫描、导入边界检查、Go 编译门禁、动态插件 WASM 构建 smoke、OpenSpec 严格校验和必要的单元测试。

回滚策略：如果实施中发现某个能力迁移风险过大，可以先保留旧路径的只读 facade，但 facade 只能委托到新`frameworkcap`或`pluginservice`实现，并必须在治理记录中标注为迁移期入口；不得继续让新代码依赖旧全局注册或旧公开实现包。

## Open Questions

- 动态插件作为 provider 的完整模型是否纳入首轮实现：建议首轮只要求动态插件能消费`frameworkcap`能力，动态插件提供 provider 通过后续独立设计处理。
