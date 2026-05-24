## ADDED Requirements

### Requirement: 源码插件升级实现不得作为公共 pkg 能力暴露

系统 SHALL 将源码插件升级、发现版本对比、发布快照同步和升级执行器视为插件运行时升级治理的内部实现。除稳定的插件管理 API、运行时状态 DTO 和必要治理契约外，旧`sourceupgrade`实现不得作为插件或外部组件可直接依赖的公共`pkg`能力；相关实现 MUST 收敛到职责明确的宿主`internal`组件。

#### Scenario: 源码插件升级通过运行时治理入口执行

- **WHEN** 管理员升级源码插件
- **THEN** 操作通过插件运行时升级治理 API 执行
- **AND** 业务代码不得直接 import 旧`pkg/sourceupgrade`执行升级逻辑

#### Scenario: 内部升级执行器不被插件依赖

- **WHEN** 源码插件需要声明升级回调、SQL 或治理资源
- **THEN** 插件通过`pluginhost`生命周期和插件 manifest 资源声明
- **AND** 插件不得依赖宿主内部 source upgrade scanner、executor 或 state reconciler
