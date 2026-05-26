## Why

当前插件相关公共包仍分散在`apps/lina-core/pkg/pluginservice`、`apps/lina-core/pkg/pluginhost`、`apps/lina-core/pkg/pluginbridge`和`apps/lina-core/pkg/sourceupgrade`等顶层路径中，`pluginservice`名称容易被理解为普通业务服务，`sourceupgrade`也不应继续作为公共`pkg`能力暴露。现在需要在没有兼容负担的前提下收敛插件公共命名空间，使源码插件贡献、动态插件桥接和插件消费宿主能力三类职责长期清晰。

## What Changes

- **BREAKING**：将插件相关公共命名空间收敛到`apps/lina-core/pkg/plugin/`下，公开顶层组件以`pluginhost`、`pluginbridge`和`capability`为主。
- **BREAKING**：将`apps/lina-core/pkg/pluginservice`迁移为`apps/lina-core/pkg/plugin/capability`，并保留`capability.Services`作为宿主基础能力集合的根接口名称。
- **BREAKING**：将`apps/lina-core/pkg/plugindb`迁移为`apps/lina-core/pkg/plugin/capability/data`，作为动态插件 data 宿主能力的公开 SDK、typed plan 契约和宿主治理适配入口。
- **BREAKING**：删除`apps/lina-core/pkg/sourceupgrade`，源码插件升级发现、对比、执行和结果 DTO 归入宿主插件运行时内部治理边界；宿主内部需要暴露的方法继续由`internal/service/plugin`服务接口承载。
- 保留动态插件 manifest 中的`hostServices`概念，但明确其只表达动态插件授权和 bridge transport 面，不作为 Go 公共包名或插件消费宿主能力的语义 owner。
- 保留`pluginbridge/protocol`作为动态插件 ABI、WASM transport、host call、host service envelope 和 codec 的唯一公开协议出口；`pluginbridge`根包只保留命名空间说明，不得把业务能力目录或协议 facade 塞回根包。
- 保留`pluginhost`作为源码插件贡献入口，负责 route、hook、cron、lifecycle、provider factory 等声明；不得让`pluginhost`重新拥有宿主能力消费目录。
- 明确`apps/lina-core/pkg/plugin/internal`只能承载`pkg/plugin/...`公共组件自身共享的内部实现，不能放置插件开发者、动态 guest SDK 或宿主插件运行时需要直接 import 的契约。
- 更新相关 OpenSpec 规范中的公共路径和边界描述，避免继续沉淀`pkg/pluginservice`、`pkg/sourceupgrade`等旧命名。

## Capabilities

### New Capabilities

- `plugin-package-boundary-governance`：定义`apps/lina-core/pkg/plugin/`命名空间下`pluginhost`、`pluginbridge`、`capability`和`internal`的职责边界、可导入规则和旧公共包移除要求。

### Modified Capabilities

- `plugin-host-service-extension`：将源码插件和动态插件消费宿主能力的统一服务集合从`pkg/pluginservice`调整为`pkg/plugin/capability`，并明确动态`hostServices`只承担授权与 transport 语义。
- `plugin-package-boundary-governance`：补充动态插件 data SDK 不得继续作为顶层`pkg/plugindb`暴露，必须归入`pkg/plugin/capability/data`公开子包。
- `pluginbridge-subcomponent-architecture`：将`pluginbridge`目标路径从`apps/lina-core/pkg/pluginbridge`调整为`apps/lina-core/pkg/plugin/pluginbridge`，并保持协议行为不变。
- `plugin-config-service`：将插件配置服务的公共契约路径从`pkg/pluginservice/config`调整为`pkg/plugin/capability/config`，保持只读配置语义不变。
- `plugin-upgrade-governance`：补充源码插件升级治理不得通过公共`pkg/sourceupgrade`暴露的要求。

## Impact

- 影响 Go 公共包导入路径：`apps/lina-core/pkg/pluginservice/**`、`apps/lina-core/pkg/pluginhost/**`、`apps/lina-core/pkg/pluginbridge/**`、`apps/lina-core/pkg/plugindb/**`和`apps/lina-core/pkg/sourceupgrade/**`。
- 影响源码插件后端入口和 provider adapter：官方插件需要改为 import 新的`pkg/plugin/pluginhost`和`pkg/plugin/capability/**`路径。
- 影响动态插件 guest SDK 和构建器：动态插件样例、WASM 构建器、bridge route 声明、data SDK 和 guest capability client 需要改为新路径，并直接使用`pluginbridge/protocol`中的协议 DTO、常量和 codec。
- 影响宿主插件运行时、WASM host service、插件管理和启动装配：需要统一注入新的`capability.Services`实例，并删除旧`sourceupgrade`facade。
- 不新增 HTTP API、SQL、前端 UI、插件 manifest 字段或运行时用户可见文案；`i18n`资源默认无影响。
- 不新增数据操作面；现有`capability`中的组织、租户、数据和缓存边界必须在迁移后保持原数据权限、批量投影和缓存一致性策略。
