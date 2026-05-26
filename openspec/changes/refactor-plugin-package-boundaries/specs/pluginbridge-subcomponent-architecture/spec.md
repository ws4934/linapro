## MODIFIED Requirements

### Requirement:pluginbridge 必须按职责提供公开子组件

系统 SHALL 将`apps/lina-core/pkg/plugin/pluginbridge`组织为职责明确的公开子组件包。子组件至少覆盖公开协议出口、bridge 合约、bridge 编解码、WASM 产物辅助、host call 协议、host service 协议、guest runtime、route binding 和 raw host call transport。根目录不得继续承载大量跨职责实现文件；根包只允许保留包说明，不得重新导出协议或 guest helper。

#### Scenario:开发者按职责定位 bridge 能力
- **WHEN** 开发者需要查看动态插件 bridge 合约、编解码、WASM 产物解析、host call、host service、guest runtime、route binding 或 raw host call transport
- **THEN** 对应源码位于语义明确的`pkg/plugin/pluginbridge/<subcomponent>/`子组件目录
- **AND** 根包目录下的生产源码文件数量保持为 1 个包说明文件，公开协议入口位于`pkg/plugin/pluginbridge/protocol`

#### Scenario:子组件名称表达稳定职责
- **WHEN** 系统完成 pluginbridge 子组件化
- **THEN** 子组件包名必须使用清晰职责名称
- **AND** 不得使用`common`、`util`、`helper`等兜底包名承载跨领域逻辑

### Requirement:公开协议出口必须唯一

系统 SHALL 以`lina-core/pkg/plugin/pluginbridge/protocol`作为 dynamic plugin bridge 的唯一公开协议出口。公开协议 DTO、ABI 常量、WASM section 名称、host call payload、host service payload 和协议 codec 必须通过`protocol`包访问；`pluginbridge`根包不得作为 facade 重新导出这些类型、常量或函数。协议实现只能存在于一个权威子组件中。

#### Scenario:目标协议 import 路径可以编译
- **WHEN** 宿主内部代码、动态插件样例或用户插件需要访问 bridge envelope、ABI 常量、WASM section、host call、host service 或 codec
- **THEN** 代码 import `lina-core/pkg/plugin/pluginbridge/protocol`
- **AND** 返回行为与迁移前目标协议语义保持一致

#### Scenario:根包不重复公开协议逻辑
- **WHEN** 开发者查看`pkg/plugin/pluginbridge`根包
- **THEN** 根包只提供包说明
- **AND** 根包不得维护 type alias、const alias、codec wrapper、protobuf wire 编解码、WASM section 遍历或 host service payload 编解码入口

### Requirement:子组件依赖方向必须防止循环依赖

系统 SHALL 明确定义`pluginbridge`子组件的依赖方向。底层合约和协议子组件不得依赖根包、guest runtime 或 route helper；`protocol`可以依赖底层 contract、artifact、codec、hostcall 和 hostservice 子组件并作为公开聚合出口。任何子组件下沉的`internal`实现包必须服务于明确父组件，不得成为跨组件兜底依赖。

#### Scenario:子组件构建无 import cycle
- **WHEN** 执行`go test ./pkg/plugin/pluginbridge/...`
- **THEN** 所有子组件包必须通过编译
- **AND** 不得出现 import cycle

#### Scenario:底层包不依赖根包
- **WHEN** 检查`contract`、`codec`、`artifact`、`hostcall`、`hostservice`子组件 import
- **THEN** 这些子组件不得 import `lina-core/pkg/plugin/pluginbridge`
- **AND** 只能依赖职责更底层或同层允许的子组件

### Requirement:宿主内部调用必须优先使用精确子组件

系统 SHALL 将项目可控的宿主内部调用逐步迁移到精确子组件 import。动态插件 guest 代码使用`pluginbridge/guest`获取 runtime、route dispatcher、request/response helper 和 raw host call transport，使用`pluginbridge/protocol`获取协议 DTO、常量和 codec；宿主 runtime、WASM host function、artifact 解析、i18n/apidoc 资源加载和 data host 应使用能表达职责边界的子组件包。

#### Scenario:宿主 runtime 使用精确子组件
- **WHEN** 宿主运行时解析动态插件产物或执行 Wasm bridge 请求
- **THEN** 代码优先 import `plugin/pluginbridge/protocol`或职责更窄的内部子组件
- **AND** 不再因为只需要单一协议能力而 import 根包

#### Scenario:插件侧 bridge guest 路径可用
- **WHEN** 动态插件 guest 代码调用`NewGuestRuntime`、`BindJSON`或其他 bridge runtime/route helper
- **THEN** 系统通过`lina-core/pkg/plugin/pluginbridge/guest`提供入口
- **AND** 这些入口不通过`pluginbridge`根包中转

#### Scenario:插件侧能力 client 路径归属 capability
- **WHEN** 动态插件 guest 代码调用`Runtime()`、`Storage()`、`Data()`、`Config()`、`Cron()`等宿主能力 client
- **THEN** 新代码应使用`lina-core/pkg/plugin/capability/guest`
- **AND** `pluginbridge`根包不得作为 facade 转发到`capability/guest`
- **AND** `pluginbridge/guest`不得重新定义这些能力 client 接口、默认实例、WASI 实现、非 WASI stub 或协议 DTO/常量别名
- **AND** 能力 client 方法签名中需要的协议 DTO、常量和 codec 由`pluginbridge/protocol`提供，`capability/guest`不得重新导出为自身别名

### Requirement:子组件化不得改变 bridge 协议行为

系统 SHALL 保证子组件化和路径迁移是结构重构，不改变动态插件 bridge 协议行为。ABI 常量、WASM custom section 名称、protobuf 字段编号、host call 状态码、host service service/method 字符串、payload 编解码结果和 guest helper 行为必须保持不变。

#### Scenario:bridge envelope 编解码保持不变
- **WHEN** 使用重构后的 API 编码并解码`BridgeRequestEnvelopeV1`或`BridgeResponseEnvelopeV1`
- **THEN** round trip 结果与重构前等价
- **AND** 现有协议测试必须继续通过

#### Scenario:host service payload codec 保持不变
- **WHEN** 使用重构后的 API 编码并解码 runtime、storage、network、data、cache、lock、config、notify 或 cron host service payload
- **THEN** round trip 结果与重构前等价
- **AND** 字段编号和默认值语义不得变化

#### Scenario:protocol 和底层子组件结果一致
- **WHEN** 同一协议调用同时通过`pluginbridge/protocol`和底层目标子组件执行
- **THEN** 两者返回相同结果或等价错误
- **AND** 测试必须覆盖至少 bridge envelope、WASM section 和 host service payload 三类代表性入口

### Requirement:子组件化必须有自动化验证

系统 SHALL 为`pluginbridge`子组件化和路径迁移提供自动化验证。验证必须覆盖公开`protocol`出口、子组件编译、宿主内部调用、动态插件样例和 Wasm guest 构建。

#### Scenario:pluginbridge 子组件测试通过
- **WHEN** 执行`go test ./pkg/plugin/pluginbridge/...`
- **THEN** `protocol`出口和所有子组件测试必须通过

#### Scenario:宿主插件运行时测试通过
- **WHEN** 执行插件运行时、WASM host function 和 data host 相关 Go 测试
- **THEN** 测试必须通过
- **AND** 不得出现因 import 迁移导致的协议行为回归

#### Scenario:动态插件样例可构建
- **WHEN** 对动态插件样例执行普通 Go 测试和`GOOS=wasip1 GOARCH=wasm`构建
- **THEN** 样例必须通过编译
- **AND** guest 侧 bridge runtime helper 与`capability/guest`能力 client 调用必须可用
