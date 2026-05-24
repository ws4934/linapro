## ADDED Requirements

### Requirement: 插件生命周期变化必须刷新 Framework Capability Provider 状态

系统 SHALL 在插件安装、启用、禁用、卸载、升级、同版本刷新和发布切换成功后，重新计算受影响的 framework capability provider 激活状态。若插件提供 provider，则其 provider 激活、撤销或替换 MUST 与插件有效 release、运行时状态和依赖校验结果一致；集群模式下 MUST 通过插件 runtime revision、事件广播、共享缓存或等价机制传播。

#### Scenario: Provider 插件升级成功后切换 Provider

- **WHEN** 提供`framework.org.v1`的插件升级成功并切换到新有效 release
- **THEN** framework capability manager 使用新 release 对应的 provider factory 重新创建或刷新 provider
- **AND** 旧 provider 不再作为 active provider 处理新调用
- **AND** 集群其他节点收到运行时修订后刷新本地 provider 状态

#### Scenario: Provider 插件禁用后能力降级

- **WHEN** 提供`framework.tenant.v1`的插件被禁用
- **THEN** framework capability manager 撤销该 provider 激活状态
- **AND** 消费 service 返回不可用状态、fallback 行为或规范定义的降级结果
- **AND** 通过`dependencies.plugins`硬依赖该 provider 插件的下游插件在后续启用、升级或健康检查中被标记为依赖不满足

### Requirement: 插件升级必须校验下游 Provider 插件依赖

插件升级 SHALL 校验升级后的 provider 插件状态不会破坏其他已安装插件通过既有`dependencies.plugins`声明的硬依赖。如果升级、禁用或发布切换会导致下游插件硬 provider 插件依赖不满足，系统 MUST 拒绝该操作或进入规范明确的阻断状态；framework capability 的可选消费仍通过运行时可用性降级表达，不引入独立 capability 依赖阻断模型。

#### Scenario: Provider 升级后不满足下游插件依赖版本

- **WHEN** 已安装插件`consumer`在`dependencies.plugins`中硬依赖`linapro-org-core`版本范围`>=1.0.0 <2.0.0`
- **AND** 管理员尝试将`linapro-org-core`升级为不满足该范围的新版本
- **THEN** 升级请求失败或要求先处理下游依赖
- **AND** 错误包含下游插件 ID、provider 插件 ID 和所需版本范围

#### Scenario: 禁用唯一 Provider 时保护下游硬依赖

- **WHEN** 插件`consumer`已启用且通过`dependencies.plugins`硬依赖唯一 tenant provider 插件
- **AND** 管理员尝试禁用唯一 tenant provider 插件
- **THEN** 禁用请求失败
- **AND** 错误包含依赖该 provider 插件的下游插件列表
