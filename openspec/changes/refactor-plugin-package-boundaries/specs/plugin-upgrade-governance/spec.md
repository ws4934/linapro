## ADDED Requirements

### Requirement: 源码插件升级治理不得通过公共`pkg/sourceupgrade`暴露

系统 SHALL 将源码插件升级发现、版本对比、升级执行、失败状态和发布切换视为宿主插件运行时内部治理能力。`apps/lina-core/pkg/sourceupgrade`公共入口 MUST 被删除；宿主内部调用方 MUST 通过`internal/service/plugin`服务接口或其内部 sourceupgrade 组件访问该能力。

#### Scenario: 宿主内部查询源码插件升级状态

- **WHEN** 宿主插件管理服务需要查询源码插件有效版本和发现版本差异
- **THEN** 调用方通过`internal/service/plugin`服务接口或其内部 sourceupgrade 组件查询
- **AND** 不得 import `lina-core/pkg/sourceupgrade`

#### Scenario: 源码插件升级执行

- **WHEN** 管理员显式升级一个源码插件
- **THEN** 插件管理 API 委托到宿主插件运行时内部 sourceupgrade 实现
- **AND** 升级流程继续遵守插件升级治理中的依赖检查、生命周期回调、SQL 迁移、治理资源同步和缓存失效要求

#### Scenario: 插件开发者声明升级资源

- **WHEN** 源码插件需要提供升级 SQL、生命周期回调或 manifest 资源
- **THEN** 插件通过`pluginhost`生命周期契约和插件资源目录声明
- **AND** 插件不得依赖公共`pkg/sourceupgrade`SDK
