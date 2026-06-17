## ADDED Requirements

### Requirement: Functions 作为 typed 计算扩展点

系统 SHALL 提供 Functions 扩展原语：宿主在自有业务流水线到达命名扩展点时，以 typed 输入调用插件实现，并消费其 typed 输出来改变业务结果。Functions SHALL 区别于 event 风格的 Hook 与单实现的 Provider SPI，并复用现有`ExtensionPoint`/`ExtensionKind`抽象新增`ExtensionKindFunction`。源码插件与动态 WASM 插件 SHALL 都可注册 Functions 实现。

#### Scenario: 结账定价调用 Functions

- **WHEN** 结账定价流水线到达`commerce.discount.run`扩展点
- **THEN** 宿主以 typed 输入调用已注册的折扣 function 实现
- **AND** 流水线消费其 typed 输出应用折扣，而非仅接收事件通知

### Requirement: Functions 确定性与资源预算

系统 SHALL 以 function 执行 profile 约束 Functions：默认禁止 host-service 调用、禁止网络、禁止非确定性时钟与随机，并施加时间、指令与内存预算。动态 WASM function 超出预算时 MUST 被终止并按扩展点声明的失败策略处理。相同输入与相同实现集 SHALL 产生相同输出，以支持结果缓存。

#### Scenario: function 超出预算

- **WHEN** 某动态 function 执行超过其时间或指令预算
- **THEN** 宿主终止该 function 执行
- **AND** 按扩展点契约的`FailOpen`或`FailClosed`策略处理该次调用

#### Scenario: 确定性可缓存

- **WHEN** 以相同输入与相同实现集重复调用同一扩展点
- **THEN** 合成结果保持一致
- **AND** 宿主可缓存该结果用于结账等热路径

### Requirement: 多实现确定性合成

系统 SHALL 使每个扩展点的多实现合成策略由扩展点 owner 在契约中声明，包括确定性排序、合并语义、数量上限与失败语义；implementer MUST NOT 改变合成策略。当注册实现数超过`MaxImpl`时，系统 MUST 截断并记录被丢弃的实现，不得静默丢弃。

#### Scenario: 多个折扣 function 合成

- **WHEN** 同一扩展点存在多个已启用 function 实现
- **THEN** 宿主按 owner 声明的确定性顺序与合并策略合成输出
- **AND** 超过`MaxImpl`的实现被截断并记录，不被静默忽略

### Requirement: 扩展点契约版本化与注册表

系统 SHALL 以版本化的 typed 输入/输出契约登记每个扩展点于宿主注册表；扩展点 owner 声明契约与调用，implementer 注册实现，二者解耦。契约演进 SHALL 通过版本号保持兼容。

#### Scenario: 扩展点契约演进

- **WHEN** 某扩展点的输入/输出契约需要演进
- **THEN** 以新版本号登记新契约，保留既有版本的兼容路径
- **AND** implementer 按声明的契约版本注册实现

### Requirement: storefront 贡献面

系统 SHALL 提供受治理的 storefront 贡献面，允许插件向店面渲染贡献 Liquid 对象、filter、section、block 与 snippet。插件 MUST NOT 在插件内执行 Liquid 渲染，渲染 SHALL 由宿主`theme-engine`完成；贡献内容取数据 SHALL 经请求级 host-service 授权快照，保持租户隔离。

#### Scenario: 应用 block 渲染

- **WHEN** 店面渲染到某第三方应用贡献的 app block
- **THEN** 宿主回调插件经授权 host-service 取数据
- **AND** 宿主用插件声明的 Liquid 片段渲染该 block，数据按租户隔离

### Requirement: 自定义数据原语 metafields 与 metaobjects

系统 SHALL 提供`metafields`与`metaobjects`原语，允许应用与商户为宿主资源定义自定义字段与自定义数据对象，keyed by owner 与`tenant_id`，并受数据权限治理。storefront SHALL 经 Drop 读取`metafields`。

#### Scenario: 定义并读取 metafield

- **WHEN** 应用为商品定义一个`metafield`并写入值
- **THEN** 该值按`tenant_id`隔离存储
- **AND** 店面商品页经 Drop 的`metafields`链读取该值

### Requirement: events 与 webhooks 原语

系统 SHALL 将现有 reserved 的`event`与`queue`能力落地为业务事件原语，提供宿主业务事件发布、插件订阅投递与对外 webhook 投递（含签名、重试与幂等）。

#### Scenario: 订单事件投递

- **WHEN** 宿主发布`order.created`业务事件
- **THEN** 已订阅该事件的插件按其订阅方式收到投递
- **AND** 已注册 webhook 的应用经签名 HTTP 投递收到事件，失败按重试与幂等策略处理

### Requirement: 扩展原语属于核心宿主增强并同步文档

系统 SHALL 将本组扩展原语视为核心宿主能力增强，遵守`core-host-boundary-governance`；每次新增或修改插件 core 扩展原语时 MUST 同步审查并更新`apps/lina-core/pkg/plugin`下的`README.md`与`README.zh-CN.md`。

#### Scenario: 新增扩展原语

- **WHEN** 某派生变更新增或修改一个插件 core 扩展原语
- **THEN** 该变更同步更新`apps/lina-core/pkg/plugin`的中英文 README
- **AND** 审查确认其属于宿主能力增强而非把商业逻辑塞入核心
