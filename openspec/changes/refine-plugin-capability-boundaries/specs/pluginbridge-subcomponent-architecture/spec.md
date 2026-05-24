## MODIFIED Requirements

### Requirement:pluginbridge 必须按职责提供公开子组件

系统 SHALL 将`apps/lina-core/pkg/pluginbridge`组织为动态插件 ABI、WASM transport 和协议 facade 组件。根包必须保留稳定 facade、包说明和必要的兼容入口；`contract`和`guest`等确需被插件作者直接依赖的契约或 guest SDK 可以保持公开。bridge 编解码、WASM 产物解析、host call dispatcher、host service payload wire 等低层实现 MUST 默认放入`pluginbridge/internal/<subcomponent>/`，不得作为插件业务开发公共 API 暴露。

#### Scenario:开发者按职责定位 bridge 能力
- **当** 开发者需要查看动态插件 bridge 合约或 guest SDK
- **则** 稳定入口位于`pkg/pluginbridge`根 facade、`pkg/pluginbridge/contract`或`pkg/pluginbridge/guest`
- **且** 低层 codec、artifact、host call 和 host service wire 实现位于`pkg/pluginbridge/internal/**`

#### Scenario:子组件名称表达稳定职责
- **当** 系统完成 pluginbridge 职责收敛
- **则** 公开子组件和内部子组件包名必须使用清晰职责名称
- **且** 不得使用`common`、`util`、`helper`等兜底包名承载跨领域逻辑

### Requirement:宿主内部调用必须优先使用精确子组件

系统 SHALL 将项目可控的宿主内部调用迁移到能表达职责边界的入口。插件侧 guest 代码可继续使用`pluginbridge`根包 facade 或`pluginbridge/guest`；宿主 runtime、WASM host function、artifact 解析、i18n/apidoc 资源加载和 data host 等内部实现应通过根 facade、同一`internal`授权边界内的精确内部子组件或更上层`pluginservice`适配调用，不得把低层实现子包重新公开为公共 API。

#### Scenario:宿主 runtime 使用受控内部入口
- **当** 宿主运行时解析动态插件产物或执行 Wasm bridge 请求
- **则** 代码使用`pluginbridge`根 facade 或`pluginbridge/internal/artifact`、`pluginbridge/internal/codec`等授权内部入口
- **且** 不为了单一协议能力重新公开`pluginbridge/artifact`、`pluginbridge/codec`、`pluginbridge/hostcall`或`pluginbridge/hostservice`

#### Scenario:插件侧兼容路径仍可用
- **当** 动态插件 guest 代码继续调用`pluginbridge.NewGuestRuntime`、`pluginbridge.BindJSON`或`pluginbridge.Runtime()`
- **则** 系统继续提供兼容入口
- **且** 这些入口委托到 guest 子组件或内部 transport 实现

## ADDED Requirements

### Requirement: pluginbridge 不得拥有宿主业务能力语义

系统 SHALL 将`pluginbridge`限定为动态插件 transport 和 ABI 层。宿主业务能力的契约、授权资源语义、降级行为和消费 service MUST 归属于`pluginservice`或`frameworkcap`；`pluginbridge`不得新增与源码插件平行的业务能力目录。

#### Scenario: 动态插件调用框架能力

- **WHEN** 动态插件调用`framework.org.v1`能力
- **THEN** `pluginbridge`只负责编码、传输和解码 host service envelope
- **AND** 能力授权、provider 激活、消费 DTO 和降级语义由`pluginservice`和`frameworkcap`处理

#### Scenario: 新增宿主能力时不修改 ABI 语义 owner

- **WHEN** 系统新增一个插件可消费宿主能力
- **THEN** 该能力首先定义在`pluginservice`或`frameworkcap`
- **AND** 动态插件仅在需要 transport 支持时扩展 bridge payload 或 handler 适配
