## MODIFIED Requirements

### Requirement: 插件宿主服务适配器必须由宿主运行期统一构造
系统 SHALL 由宿主运行期统一构造并发布源码插件和动态插件 host service 适配器。适配器 MUST 复用启动期共享的宿主服务实例或共享后端，MUST 不在插件调用路径中自行构造孤立宿主服务图。源码插件和动态插件访问同一宿主能力时 MUST 共享`pkg/plugin/capability`与`capability.Services`能力集合语义，动态插件 host service handler 只作为 transport 与授权适配层。

#### Scenario: 源码插件使用宿主服务适配器
- **WHEN** 源码插件调用`pkg/plugin/capability/*`发布的宿主能力
- **THEN** 该能力适配器由宿主运行期构造并通过 registrar 传递给插件
- **AND** 适配器复用宿主共享的 auth、session、notify、config、i18n、pluginstate、orgcap、tenantcap 或其他依赖
- **AND** 插件生产路径不得无参创建该适配器

#### Scenario: 动态插件 host service 调用共享宿主能力
- **WHEN** 动态插件通过统一 host service 协议调用 cache、lock、notify、config、runtime、storage、data 或 capability 等宿主能力
- **THEN** host service handler 使用插件 runtime 注入的共享宿主服务或共享后端
- **AND** handler 不得在每次调用中创建独立 cache、lock、notify、config、plugin service 或 capability manager 实例
- **AND** data host service 的 guest SDK、typed plan 和宿主治理 facade 归属`pkg/plugin/capability/data`

#### Scenario: WASM host service 配置入口由启动期注入
- **WHEN** 宿主启动并初始化 WASM host service
- **THEN** 启动路径显式配置 cache、lock、notify、storage、config、runtime、orgcap、tenantcap 和`capability.Services`等 host service 的共享依赖
- **AND** 包级默认实例不得在生产启动后继续作为实际运行依赖

## ADDED Requirements

### Requirement: capability 必须作为源码插件和动态插件的统一能力集合

系统 SHALL 将`pkg/plugin/capability`定义为插件消费宿主能力的统一集合，并将根接口命名为`capability.Services`。源码插件 MUST 通过 registrar 或等价上下文获取 capability services；动态插件 MUST 通过 guest client 和 host service handler 进入同一组服务语义；两类插件不得分别使用不同组件暴露同一能力。

#### Scenario: 源码插件和动态插件读取同一能力

- **WHEN** 源码插件和动态插件分别消费插件作用域配置、宿主公开配置、manifest、数据服务或框架 capability
- **THEN** 二者共享同一 service 契约、授权边界、错误语义和降级策略
- **AND** 仅 transport 和运行时加载方式存在差异

#### Scenario: 新能力只注册到统一目录

- **WHEN** 系统新增一个插件可消费宿主能力
- **THEN** 该能力必须注册到`pkg/plugin/capability`根服务集合或其子包
- **AND** 动态插件的 host service handler 只把 bridge 请求映射到该统一服务集合

### Requirement: 动态 hostServices 与 capability services 必须语义分层

系统 SHALL 保留动态插件 manifest 中的`hostServices`作为授权声明和 bridge transport 调用面，同时将宿主能力语义归属到`pkg/plugin/capability`和`capability.Services`。`hostServices`不得被重新解释为 Go 公共能力集合名称，`capability`也不得绕过动态插件授权快照直接授予动态插件访问权。

#### Scenario: 动态插件声明 hostServices

- **WHEN** 动态插件在`plugin.yaml`中声明`hostServices`
- **THEN** 该声明表示动态插件申请调用的 service、method 和资源边界
- **AND** 宿主在安装、启用或升级阶段生成授权快照
- **AND** 该声明不改变`pkg/plugin/capability`中能力契约的 owner

#### Scenario: 动态插件通过 capability guest client 调用宿主能力

- **WHEN** 动态插件调用`pkg/plugin/capability/guest`中的能力 client
- **THEN** guest client 通过`pkg/plugin/pluginbridge/guest`raw transport 发起 host service 调用
- **AND** 宿主先校验`hostServices`授权快照
- **AND** 授权通过后再委托到同一`capability.Services`能力集合
- **AND** `RuntimeHostService`、`StorageHostService`、`ConfigHostService`、`DataHostService`等 guest 能力 client 接口、默认实例、WASI 实现和非 WASI stub 均归属`pkg/plugin/capability/guest`
- **AND** `pkg/plugin/pluginbridge/guest`只提供 raw host call transport、guest runtime 和 route binding，不拥有宿主能力 client 语义
- **AND** 能力 client 方法签名中的 host service DTO、cron 合约、日志等级和 codec 均直接使用`pkg/plugin/pluginbridge/protocol`，`capability/guest`不得重复导出这些协议别名

#### Scenario: 动态插件通过 data SDK 调用宿主 data 能力

- **WHEN** 动态插件调用`pkg/plugin/capability/data`中的 ORM-style data SDK
- **THEN** data SDK 通过`pkg/plugin/pluginbridge/guest`raw transport 和`pkg/plugin/pluginbridge/protocol`协议 DTO 发起 data host service 调用
- **AND** 宿主先校验`hostServices`中的 data 授权快照、资源表和方法范围
- **AND** 授权通过后再执行同一 data host service 治理路径
