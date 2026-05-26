## MODIFIED Requirements

### Requirement:宿主和插件必须通过稳定能力接缝解耦

系统 SHALL 通过`pluginservice`、能力接口、事件 Hook、路由注册器和 Cron 注册器等稳定接缝完成宿主与插件的协作，而非在宿主业务代码中散落插件特定的占位逻辑、大量`if pluginEnabled`分支或对插件内部实现的直接依赖。

#### Scenario:宿主调用可选组织能力
- **当** 用户管理、认证或其他宿主核心模块需要访问部门、岗位、组织树或组织数据范围等可选能力时
- **则** 宿主通过显式注入的`orgcap.Service`或由`pluginservice.Services.Org()`发布的组织能力入口访问这些能力
- **且** `linapro-org-core`的插件状态判断和功能分支不直接散落在宿主实现中
- **且** 宿主仅持有该能力的接口、DTO、消费 service 和空实现，不直接查询或维护`linapro-org-core`的物理表

#### Scenario:宿主调用可选租户能力
- **当** 认证、会话、数据权限或插件 host service 需要访问租户上下文时
- **则** 宿主通过显式注入的`tenantcap.Service`或由`pluginservice.Services.Tenant()`发布的租户能力入口访问该能力
- **且** 租户 provider 的具体实现由满足生命周期条件的插件提供
- **且** 宿主不得在业务路径中直接 import 租户插件的内部 service 或 DAO

#### Scenario:宿主扩展插件日志或监控能力
- **当** 非核心能力拆分为源码插件时
- **则** 宿主仅保留稳定的事件、治理接口和注册入口
- **且** 不在宿主控制器、服务或中间件中为个别插件保留大量功能占位逻辑

## ADDED Requirements

### Requirement: 宿主定义插件实现的框架能力必须归属独立 pluginservice 能力组件

系统 SHALL 将“宿主定义接口、插件提供实现”的框架能力按能力领域归属到`pkg/pluginservice/orgcap`和`pkg/pluginservice/tenantcap`等独立组件。每个能力组件必须维护自身公开契约、DTO、消费`Service`、fallback、delegation 和 provider factory facade；共享 provider registry 与激活治理实现必须位于`pkg/pluginservice/internal/capabilityregistry`，不得再保留`pkg/frameworkcap`聚合包、旧`pkg/orgcap`/`pkg/tenantcap`兼容包或宿主`internal/service/orgcap`、`internal/service/tenantcap`双重适配层。

#### Scenario: orgcap 迁移为 pluginservice 独立组件

- **WHEN** 系统维护组织能力契约
- **THEN** 组织能力公开消费契约由`pkg/pluginservice/orgcap`暴露
- **AND** 消费方通过注入的`orgcap.Service`或`pluginservice.Services.Org()`获取组织能力实例
- **AND** 旧`pkg/orgcap`不得作为新代码入口
- **AND** fallback、delegation、provider factory facade 位于`pkg/pluginservice/orgcap`，共享 provider registry 和激活实现位于`pkg/pluginservice/internal/capabilityregistry`

#### Scenario: tenantcap 迁移为 pluginservice 独立组件

- **WHEN** 系统维护租户能力契约
- **THEN** 租户能力公开消费契约由`pkg/pluginservice/tenantcap`暴露
- **AND** 消费方通过注入的`tenantcap.Service`或`pluginservice.Services.Tenant()`获取租户能力实例
- **AND** 旧`pkg/tenantcap`不得作为新代码入口
- **AND** provider 实现由插件通过`tenantcap.Provide(...)`声明并由`pkg/pluginservice/internal/capabilityregistry`中的生命周期治理激活
