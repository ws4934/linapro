## MODIFIED Requirements

### Requirement:宿主和插件必须通过稳定能力接缝解耦

系统 SHALL 通过`frameworkcap`、`pluginservice`、能力接口、事件 Hook、路由注册器和 Cron 注册器等稳定接缝完成宿主与插件的协作，而非在宿主业务代码中散落插件特定的占位逻辑、大量`if pluginEnabled`分支或对插件内部实现的直接依赖。

#### Scenario:宿主调用可选组织能力
- **当** 用户管理、认证或其他宿主核心模块需要访问部门、岗位、组织树或组织数据范围等可选能力时
- **则** 宿主通过`frameworkcap.Org()`消费 service 或由`pluginservice`发布的 framework capability 入口访问这些能力
- **且** `linapro-org-core`的插件状态判断和功能分支不直接散落在宿主实现中
- **且** 宿主仅持有该能力的接口、DTO、消费 service 和空实现，不直接查询或维护`linapro-org-core`的物理表

#### Scenario:宿主调用可选租户能力
- **当** 认证、会话、数据权限或插件 host service 需要访问租户上下文时
- **则** 宿主通过`frameworkcap.Tenant()`消费 service 访问该能力
- **且** 租户 provider 的具体实现由满足生命周期条件的插件提供
- **且** 宿主不得在业务路径中直接 import 租户插件的内部 service 或 DAO

#### Scenario:宿主扩展插件日志或监控能力
- **当** 非核心能力拆分为源码插件时
- **则** 宿主仅保留稳定的事件、治理接口和注册入口
- **且** 不在宿主控制器、服务或中间件中为个别插件保留大量功能占位逻辑

## ADDED Requirements

### Requirement: 宿主定义插件实现的框架能力必须归属 frameworkcap

系统 SHALL 将“宿主定义接口、插件提供实现”的框架能力归属到`pkg/frameworkcap`。这类能力的公开契约必须由`frameworkcap`根包统一暴露，例如`frameworkcap.Org()`和`frameworkcap.Tenant()`；fallback、delegation、provider registry 和激活治理实现必须位于`frameworkcap/internal`，不得把接口留在一个公共包、fallback 留在宿主内部 service、provider 注册散落在插件路由注册路径，也不得为每个能力新增公开子组件增加使用复杂度。

#### Scenario: orgcap 迁移为 frameworkcap

- **WHEN** 系统维护组织能力契约
- **THEN** 组织能力公开消费契约由`pkg/frameworkcap`根包暴露
- **AND** 消费方通过`frameworkcap.Org()`获取组织能力实例
- **AND** 旧`pkg/orgcap`不得作为新代码入口
- **AND** fallback、delegation、provider registry 和激活实现位于`pkg/frameworkcap/internal`

#### Scenario: tenantcap 迁移为 frameworkcap

- **WHEN** 系统维护租户能力契约
- **THEN** 租户能力公开消费契约由`pkg/frameworkcap`根包暴露
- **AND** 消费方通过`frameworkcap.Tenant()`获取租户能力实例
- **AND** 旧`pkg/tenantcap`不得作为新代码入口
- **AND** provider 实现由插件通过`frameworkcap`根包的 provider factory facade 声明并由`frameworkcap/internal`中的生命周期治理激活
