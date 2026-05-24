## ADDED Requirements

### Requirement: 框架能力必须分离公开契约和内部治理实现

系统 SHALL 通过`pkg/frameworkcap`集中维护由宿主定义契约且由插件提供实现的框架能力。`frameworkcap`根包 MUST 作为唯一公开导入入口，直接维护 capability ID、版本、DTO、`OrgService`、`TenantService`等消费 service 接口、provider factory 声明 facade 和必要错误类型；系统 MUST NOT 为组织、租户等具体能力新增`frameworkcap/org`、`frameworkcap/tenant`公开子组件。provider registry、fallback、delegation、provider 激活状态、冲突检测、缓存和 manager 实现 MUST 放入`frameworkcap/internal`下的职责明确子组件。接口契约不得泄漏提供方插件的`internal/service`、`DAO`、`DO`、`Entity`、缓存快照、运行时状态或私有配置。

#### Scenario: 组织能力集中维护

- **WHEN** 系统定义`framework.org.v1`能力
- **THEN** 该能力的 capability ID、DTO、`OrgService`接口和`ProvideOrg(...)`声明 facade 位于`pkg/frameworkcap`根包公开边界下
- **AND** fallback、delegation、provider registry 和激活治理位于`pkg/frameworkcap/internal`边界下
- **AND** 消费方不得依赖提供方插件的 provider adapter 或内部业务 service

#### Scenario: 租户能力通过根组件获取

- **WHEN** 系统定义`framework.tenant.v1`能力
- **THEN** 消费方通过`frameworkcap.Tenant()`获取租户能力实例
- **AND** 系统不得要求消费方 import `pkg/frameworkcap/tenant`
- **AND** 租户能力内部实现位于`pkg/frameworkcap/internal`边界下

#### Scenario: 能力契约不泄漏内部模型

- **WHEN** framework capability service 返回组织、租户、用户或可见范围信息
- **THEN** 返回值使用该能力契约自有的 DTO、投影或值对象
- **AND** 返回值不得包含宿主或插件内部`DAO`、`DO`、`Entity`、`*gdb.Model`或私有缓存结构

### Requirement: 插件必须通过 Provider Factory 声明框架能力实现

系统 SHALL 要求插件通过 provider factory 声明其对 framework capability 的实现。Provider factory MUST 通过`frameworkcap`根包的窄 facade 在插件入口或 registrar 阶段声明，例如`frameworkcap.ProvideOrg(...)`或`frameworkcap.ProvideTenant(...)`；provider 实例 MUST 由`frameworkcap/internal`中的 capability manager 根据插件安装、启用、版本、依赖、租户状态、升级和禁用状态激活；插件不得在路由注册、controller 构造或业务请求路径中直接写入全局 provider 注册表。

#### Scenario: 源码插件声明组织能力 Provider

- **WHEN** `linapro-org-core`提供`framework.org.v1`实现
- **THEN** 插件入口通过`frameworkcap.ProvideOrg(...)`声明一个组织能力 provider factory
- **AND** `frameworkcap/internal`中的 framework capability manager 在该插件满足安装、启用和依赖条件后创建并激活 provider
- **AND** 路由注册回调不得直接调用全局`RegisterProvider(provider)`完成激活

#### Scenario: 插件禁用后 Provider 被撤销

- **WHEN** 提供 framework capability 的插件被禁用、卸载或升级失败
- **THEN** `frameworkcap/internal`中的 framework capability manager 撤销该插件对应 provider 的激活状态
- **AND** 消费 service 的`Available(ctx)`或等价状态反映该能力不可用或降级状态

### Requirement: Provider 实现必须封装在提供方插件内部

系统 SHALL 将官方插件的 framework capability provider adapter 视为插件内部实现。Provider adapter MUST 默认位于`apps/lina-plugins/<plugin-id>/backend/internal/provider/<capability>adapter/`，真实业务实现 MUST 位于`backend/internal/service/`；`backend/plugin.go`只负责声明路由、生命周期和 provider factory，不得承载 provider 业务实现。

#### Scenario: Provider Adapter 不作为公开包暴露

- **WHEN** 一个官方源码插件实现`framework.tenant.v1`
- **THEN** 该插件的 tenant provider adapter 位于`backend/internal/provider/tenantadapter/`
- **AND** 其他插件无法通过 Go import 直接依赖该 adapter
- **AND** 其他插件只能通过 framework capability 消费 service 使用该能力

#### Scenario: 业务逻辑保留在插件 Service

- **WHEN** provider adapter 需要访问插件业务数据或领域逻辑
- **THEN** adapter 通过显式依赖调用同插件`backend/internal/service/`中的业务 service
- **AND** adapter 只做 framework capability 契约到内部业务 service 的薄适配

### Requirement: 插件消费框架能力必须通过消费 Service

系统 SHALL 要求宿主模块、源码插件和动态插件通过 framework capability 的消费 service 使用能力。消费方 MUST 不直接获取 provider 实例，也不得直接 import 提供方插件的 provider adapter、内部 service、DAO、Entity 或其他内部实现；消费方在需要硬阻断时 MUST 使用既有 provider 插件依赖，在可选使用能力时 MUST 支持运行时可用性降级语义。

#### Scenario: 源码插件消费组织能力

- **WHEN** 源码插件需要读取组织树或批量解析用户组织信息
- **THEN** 插件通过`pluginservice.Services.Framework().Org()`、`frameworkcap.Org()`或等价注入的消费 service 发起调用
- **AND** 插件不得 import `linapro-org-core/backend/internal/**`

#### Scenario: 动态插件消费组织能力

- **WHEN** 动态插件需要消费`framework.org.v1`
- **THEN** guest SDK 通过`pluginservice/guest`发起版本化 host service 调用
- **AND** 宿主将调用路由到同一个 framework capability 消费 service
- **AND** 调用必须满足动态插件`hostServices`授权；若消费方需要硬依赖具体 provider 插件，则由既有`dependencies.plugins`声明和生命周期校验表达

#### Scenario: 可选能力不可用时降级

- **WHEN** 插件可选使用`framework.org.v1`
- **AND** 当前环境没有可用组织能力 provider
- **THEN** 插件仍可启用
- **AND** 插件必须通过`Available(ctx)`或等价能力状态隐藏相关功能、返回零值或执行规范定义的降级行为

### Requirement: 框架能力必须治理 Provider 冲突和版本

系统 SHALL 为 framework capability 定义 provider 选择、冲突处理和版本兼容规则。单例能力 MUST 在多个 provider 同时满足条件时拒绝激活或按规范定义的确定性优先级选择；不兼容契约变更 MUST 使用新的 capability 版本，不得破坏既有`v1`消费契约。

#### Scenario: 单例能力重复 Provider 被拒绝

- **WHEN** 两个已启用插件同时声明同一个单例 capability 的 provider
- **THEN** framework capability manager 按该能力规范拒绝第二个 provider 或进入明确冲突状态
- **AND** 错误或状态包含 capability ID 和冲突插件 ID

#### Scenario: 不兼容契约使用新版本

- **WHEN** `framework.org.v1`无法兼容新增的响应结构、错误语义或授权边界
- **THEN** 系统定义`framework.org.v2`
- **AND** 已声明依赖`framework.org.v1`的插件继续按`v1`契约运行或降级

### Requirement: 框架能力状态必须保持缓存一致性

系统 SHALL 将 provider 激活状态和能力可用性快照视为关键运行时数据。插件安装、启用、禁用、卸载、升级、同版本刷新、既有插件依赖状态变化或租户状态变化后，系统 MUST 触发受影响 capability 的状态刷新；集群模式下不得只刷新当前节点本地内存。

#### Scenario: 插件启用后集群节点刷新能力状态

- **WHEN** 集群模式下启用一个提供`framework.tenant.v1`的插件
- **THEN** 主节点更新 provider 激活状态并发布插件运行时修订或等价事件
- **AND** 其他节点观察到事件后刷新本地能力可用性快照

#### Scenario: 只读能力检查不触发全局失效

- **WHEN** 插件或宿主只读取某个 capability 的`Available(ctx)`状态
- **THEN** 系统不得写入插件 registry 或清空无关插件、语言包、路由或前端 bundle 缓存
- **AND** 只读检查不得产生跨实例失效事件
