## MODIFIED Requirements

### Requirement: 动态插件通过受限 ORM 风格 SDK 访问 data service

系统 SHALL 为动态插件提供`pkg/plugindb`受限 ORM 风格 SDK，作为`data service`的推荐 guest 侧访问入口；该 SDK 必须继续建立在结构化 hostService 协议和`pluginservice/data`语义之上，而不是直接向插件暴露完整`gdb.DB`、`gdb.Model`、宿主 DAO、typed plan 内部结构或 host-side DB wrapper。

#### Scenario: 插件通过 plugindb 发起单表查询

- **WHEN** 插件作者使用`plugindb.Open().Table(...).WhereEq(...).Page(...).All()`等链式 API 访问数据
- **THEN** guest SDK 将该请求转换为结构化、可验证的数据访问请求
- **AND** 宿主继续按当前 release 授权快照、字段白名单和数据范围执行治理
- **AND** 插件不得借此获得 raw SQL、JOIN、任意表达式拼接能力或 host-side query plan 内部类型

#### Scenario: plugindb 作为推荐路径而兼容层暂时保留

- **WHEN** 宿主开始引入`plugindb`guest SDK
- **THEN** 开发文档、样例和 demo 应优先使用`plugindb.Open()`作为数据访问主路径
- **AND** 旧的`pluginbridge.Data()`可作为兼容层短期保留
- **AND** 宿主不得要求插件作者直接拼装底层 hostService envelope

## ADDED Requirements

### Requirement: plugindb 的 Host-Side 实现不得作为公共 API

系统 SHALL 将`plugindb`中的 typed plan、plan codec、host DB wrapper、DoCommit 拦截、执行器、审计上下文和治理校验实现视为宿主内部细节。插件可导入的公共 API MUST 限于 guest 侧受限 DSL、公开枚举、公开 DTO 和必要 facade；host-side 实现 MUST 位于`plugindb/internal/**`、`pluginservice/data`内部边界或宿主`internal/service/plugin/internal/datahost`中。

#### Scenario: 插件不能导入 Host DB Wrapper

- **WHEN** 插件代码尝试 import `pkg/plugindb/internal/host`或旧的公开 host-side DB 包
- **THEN** Go internal 边界或治理扫描阻止该依赖
- **AND** 插件只能通过`plugindb`guest DSL 发起结构化 data service 请求

#### Scenario: Host 执行器通过 Facade 解码计划

- **WHEN** 宿主 data host 需要执行动态插件数据请求
- **THEN** 宿主通过受控 facade 或内部包解码并校验 typed plan
- **AND** 解码后的内部计划不得回传给插件或成为插件接口契约
