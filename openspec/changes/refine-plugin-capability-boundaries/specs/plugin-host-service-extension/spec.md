## MODIFIED Requirements

### Requirement: 动态插件通过版本化宿主服务协议获取全部宿主能力

系统 SHALL 在保留`lina_env.host_call`入口的前提下，为动态插件提供版本化的宿主服务调用协议。动态插件访问宿主能力时 MUST 通过统一宿主服务通道进入`pluginservice`能力目录或其受控适配器，而不是继续线性增加新的专用 opcode，也不得让`pluginbridge`成为与源码插件平行的宿主能力语义 owner。

#### Scenario: Guest 调用结构化宿主服务

- **WHEN** guest SDK 发起一次宿主服务调用
- **THEN** 宿主通过统一请求 envelope 解析`service`、`method`、资源标识（如`storage.resources.paths`、URL 模式、`table`、pluginservice capability ID 或 manifest 资源路径）和请求载荷
- **AND** 宿主服务注册表定位对应的服务处理器
- **AND** 服务处理器委托到`pluginservice`能力目录、`orgcap.Service`、`tenantcap.Service`或其他受控宿主适配器
- **AND** 宿主以统一响应 envelope 返回业务结果或结构化错误

#### Scenario: 未知 service 或 method 被拒绝

- **WHEN** 插件调用一个宿主不支持的`service`或`method`
- **THEN** 宿主返回显式的“不支持”错误
- **AND** 宿主不得进入任何实际宿主服务逻辑

#### Scenario: 已有最小 Host Call 被统一重构

- **WHEN** 宿主收敛当前动态插件能力模型
- **THEN** 已实现的日志、状态和数据访问能力也通过统一宿主服务协议对 guest 暴露
- **AND** 宿主不得继续维护面向插件的平行公开协议

#### Scenario: 动态插件消费 Pluginservice Capability

- **WHEN** 动态插件通过 guest SDK 调用`framework.org.v1`
- **THEN** host service handler 校验动态插件的`hostServices`授权
- **AND** 调用进入`pluginservice.Services.Org()`对应的消费 service
- **AND** 如该动态插件需要硬依赖具体 provider 插件，则由既有`dependencies.plugins`在生命周期路径中校验
- **AND** `pluginbridge`仅承担 transport 和 payload 编解码职责

### Requirement: 插件宿主服务适配器必须由宿主运行期统一构造
系统 SHALL 由宿主运行期统一构造并发布源码插件和动态插件 host service 适配器。适配器 MUST 复用启动期共享的宿主服务实例或共享后端，MUST 不在插件调用路径中自行构造孤立宿主服务图。源码插件和动态插件访问同一宿主能力时 MUST 共享`pluginservice`能力目录语义，动态插件 host service handler 只作为 transport 适配层。

#### Scenario: 源码插件使用宿主服务适配器
- **WHEN** 源码插件调用`pkg/pluginservice/*`发布的宿主能力
- **THEN** 该能力适配器由宿主运行期构造并通过 registrar 传递给插件
- **AND** 适配器复用宿主共享的 auth、session、notify、config、i18n、pluginstate、orgcap、tenantcap 或其他依赖
- **AND** 插件生产路径不得无参创建该适配器

#### Scenario: 动态插件 host service 调用共享宿主能力
- **WHEN** 动态插件通过统一 host service 协议调用 cache、lock、notify、config、runtime、storage、data 或 pluginservice capability 等宿主能力
- **THEN** host service handler 使用插件 runtime 注入的共享宿主服务或共享后端
- **AND** handler 不得在每次调用中创建独立 cache、lock、notify、config、plugin service 或 capability manager 实例

#### Scenario: WASM host service 配置入口由启动期注入
- **WHEN** 宿主启动并初始化 WASM host service
- **THEN** 启动路径显式配置 cache、lock、notify、storage、config、runtime、orgcap、tenantcap 和 pluginservice 等 host service 的共享依赖
- **AND** 包级默认实例不得在生产启动后继续作为实际运行依赖

## ADDED Requirements

### Requirement: pluginservice 必须作为源码插件和动态插件的统一能力目录

系统 SHALL 将`pkg/pluginservice`定义为插件消费宿主能力的统一目录。源码插件 MUST 通过 registrar 或等价上下文获取`pluginservice.Services`；动态插件 MUST 通过 guest client 和 host service handler 进入同一组服务语义；两类插件不得分别使用不同组件暴露同一能力。

#### Scenario: 源码插件和动态插件读取同一能力

- **WHEN** 源码插件和动态插件分别消费插件作用域配置、宿主公开配置、manifest、数据服务或 pluginservice capability
- **THEN** 二者共享同一 service 契约、授权边界、错误语义和降级策略
- **AND** 仅 transport 和运行时加载方式存在差异

#### Scenario: 新能力只注册到统一目录

- **WHEN** 系统新增一个插件可消费宿主能力
- **THEN** 该能力必须注册到`pluginservice.Services`或其子目录
- **AND** 动态插件的 host service handler 只把 bridge 请求映射到该统一目录
