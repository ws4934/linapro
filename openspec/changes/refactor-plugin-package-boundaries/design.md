## Context

当前代码已经把插件能力边界做过一次收敛：`pluginhost`偏源码插件贡献入口，`pluginbridge`偏动态插件 ABI 和 transport，`pluginservice`承担插件消费宿主能力目录，`sourceupgrade`的真实执行逻辑已经下沉到`internal/service/plugin/internal/sourceupgrade`。这次重构继续推进命名和包边界，目标是消除`pluginservice`和`sourceupgrade`两个长期容易误解的公共入口。

当前需要特别遵守 Go `internal`可见性规则：`apps/lina-core/pkg/plugin/internal/**`只能被`apps/lina-core/pkg/plugin/...`子树下的包导入，不能被`apps/lina-plugins/...`源码插件、动态插件 guest 代码，也不能被`apps/lina-core/internal/service/plugin/...`宿主运行时直接导入。因此插件开发者需要 import 的契约、DTO、guest SDK、`bizctx`、`config`、`manifest`和`orgcap`/`tenantcap`等能力不能放入`pkg/plugin/internal`。

## Goals / Non-Goals

**Goals:**

- 建立`apps/lina-core/pkg/plugin/`统一命名空间，使插件相关公共组件按职责集中。
- 用`pkg/plugin/capability`替代`pkg/pluginservice`，表达“宿主发布给插件消费的稳定基础能力集合”，避免与普通业务 service 层混淆。
- 保持动态插件 manifest 中`hostServices`作为授权和 transport 概念，不将其作为 Go 公共组件名。
- 删除`apps/lina-core/pkg/sourceupgrade`，将源码插件升级治理保留在宿主插件运行时内部。
- 保持源码插件和动态插件访问同一宿主能力时共享语义契约、授权模型、错误语义、数据权限和缓存一致性策略。
- 不保留旧路径 facade 或兼容包，因为项目无历史兼容负担。

**Non-Goals:**

- 不改变动态插件`hostServices`manifest 字段名、授权快照语义或 bridge service/method 字符串。
- 不改变插件 HTTP API、插件管理页面、SQL 资源路径、语言包资源路径或插件生命周期资源目录。
- 不重写组织、租户、配置、缓存、通知、manifest、host config 等能力本身的业务语义。
- 不把源码插件升级治理设计成外部 SDK；外部管理能力仍通过插件管理 API 或宿主内部服务接口表达。

## Decisions

### 1. 公共命名空间收敛到`pkg/plugin`

目标结构如下：

```text
apps/lina-core/pkg/plugin/
  pluginhost/
  pluginbridge/
  capability/
    data/
    contract/
    guest/
    bizctx/
    config/
    hostconfig/
    manifest/
    orgcap/
    pluginlifecycle/
    pluginstate/
    tenantcap/
    tenantfilter/
    internal/
      capabilityregistry/
  internal/
```

`pkg/plugin/internal`只服务`pkg/plugin/pluginhost`、`pkg/plugin/pluginbridge`和`pkg/plugin/capability`这些公共组件自身。插件作者或宿主运行时需要直接导入的接口必须放在公开子包中。

备选方案是把`bizctx`、`guest`等迁入`pkg/plugin/internal`。该方案违反 Go `internal`导入限制，会让源码插件和动态插件无法编译，因此不采用。

### 2. 保留`capability/internal`作为 capability 私有实现边界

`pkg/plugin/capability/internal`和`pkg/plugin/internal`不是重复目录，二者表达不同可见性：

```text
pkg/plugin/capability/internal
  only importable by pkg/plugin/capability/...
  capability registry / provider lazy loading / conflict detection

pkg/plugin/internal
  importable by pkg/plugin/pluginhost
                pkg/plugin/pluginbridge
                pkg/plugin/capability
  cross-component shared internals only
```

因此`capability/internal/capabilityregistry`这类只服务能力目录的实现必须保留在`capability/internal`下，不迁入`pkg/plugin/internal`。如果迁入总`internal`，`pluginhost`和`pluginbridge`在编译层面也能导入 capability registry，会扩大内部实现可见范围，并增加后续职责重新耦合的风险。

`pkg/plugin/internal`可以作为预留共享边界，但只有当某个实现确实同时服务`pluginhost`、`pluginbridge`和`capability`中的多个公共组件时才使用。没有跨组件复用需求时，不要求创建或填充该目录。

### 3. `pluginservice`重命名为`capability`

`capability`比`hostservice`更合适，因为动态插件 manifest 中已经存在`hostServices`，继续使用`hostservice`作为 Go 包名会混淆授权 transport 面和语义能力集合。

迁移后的核心语义是：

```text
source plugin
  -> pkg/plugin/pluginhost registrar
  -> pkg/plugin/capability Services
  -> config / manifest / orgcap / tenantcap / cache / notify / ...

dynamic plugin
  -> pkg/plugin/capability/guest client / pkg/plugin/capability/data SDK
  -> pkg/plugin/pluginbridge transport
  -> hostServices authorization snapshot
  -> same capability Services
```

根能力集合命名为`Services`。原因是根包已经是`capability`，完整限定名`capability.Services`能直接表达“宿主基础能力集合”，同时避免`Directory`在插件资源目录、菜单目录和 guest 侧 client directory 语义之间产生额外歧义。动态插件 manifest 中的`hostServices`仍只表达授权和 transport，不作为 Go 公共包名或 capability 根接口 owner。

### 4. `pluginbridge`只保留协议和 transport 职责

迁移到`pkg/plugin/pluginbridge`后，`pluginbridge`继续负责 ABI、WASM request/response envelope、host call、host service envelope、codec、artifact section 和 dynamic guest runtime。`pkg/plugin/pluginbridge/protocol`是唯一公开协议出口，统一承载 bridge envelope、ABI 常量、WASM section、host call payload、host service payload 和协议 codec 的公开别名。`pluginbridge`根包只保留命名空间说明，不再通过 facade 重新导出协议 DTO、常量、codec 或 guest helper，避免形成第二个公开协议面。

动态插件的`hostServices`仍然存在于 manifest 和 bridge 授权层，它表达“这个动态插件被允许调用哪些宿主 service/method/resource”。真正的宿主能力语义由`pkg/plugin/capability`持有。

`pkg/plugin/pluginbridge/guest`必须保持薄层职责，只承载 guest runtime、controller route dispatcher、request/response binding 和 raw host call transport。该包可以在 runtime/helper 方法签名中直接使用`protocol`类型，但不得再为 request/response envelope、route snapshot、host call transport、response helper、常量或 codec 创建导出别名；动态 guest 业务代码需要 bridge 协议 DTO、常量或 codec 时应直接导入`pkg/plugin/pluginbridge/protocol`。

动态 guest 侧的`RuntimeHostService`、`StorageHostService`、`ConfigHostService`、`DataHostService`等能力 client 接口、默认实例、WASI 实现和非 WASI stub 归属`pkg/plugin/capability/guest`。`capability/guest`可以在能力 client 方法签名中使用`protocol`类型，但不得再为`protocol`的 DTO、常量或 codec 创建导出别名；动态插件业务代码需要日志等级、cron 合约、storage/network/cache/notify 等协议结构时应使用`protocol.*`。

### 5. `pluginhost`只保留源码插件贡献职责

迁移到`pkg/plugin/pluginhost`后，`pluginhost`继续负责源码插件注册、静态资源、HTTP route、hook、cron、lifecycle、provider factory declaration 和治理声明。源码插件需要消费宿主能力时，通过 registrar 拿到`capability.Services`，而不是通过`pluginhost`拥有自己的 host service 目录。

### 6. 删除公共`pkg/sourceupgrade`

`apps/lina-core/pkg/sourceupgrade`当前只是 delegate facade 和 contract，真实执行器已在`internal/service/plugin/internal/sourceupgrade`。迁移后：

- 源码插件升级发现、对比、执行、失败状态和 release 切换继续由`internal/service/plugin/internal/sourceupgrade`实现。
- 宿主内部服务接口在`internal/service/plugin`层暴露`ListSourceUpgradeStatuses`、`UpgradeSourcePlugin`和`ValidateSourcePluginUpgradeReadiness`等方法。
- DTO 可以定义在`internal/service/plugin`或其内部 sourceupgrade 组件中，再由宿主内部 facade 按需别名，不再通过公共`pkg/sourceupgrade/contract`发布。
- 旧`pkg/sourceupgrade`目录和测试删除，不保留兼容 facade。

不采用迁移到`pkg/plugin/internal/sourceupgrade`的方案，因为宿主插件运行时不在`pkg/plugin`子树下，无法直接导入该 internal 包。

### 7. `plugindb`迁移为`capability/data`

`apps/lina-core/pkg/plugindb`同时承担动态插件 guest-side ORM facade、typed data plan 契约和宿主侧受治理 DB wrapper。该职责本质上属于插件消费宿主 data 能力的公开子能力，应迁入`apps/lina-core/pkg/plugin/capability/data`，而不是继续保留为顶层`pkg/plugindb`。

不采用迁入`pkg/plugin/capability/internal`的方案，因为动态插件 guest 代码、宿主 datahost 和测试都需要直接导入该公开契约；放入`internal`会违反 Go `internal`可见性规则并导致动态插件或宿主运行时无法编译。

迁移后包名必须与目录职责保持一致，使用`package data`、主文件`data.go`和`data_*.go`文件前缀表达动态插件 data capability SDK 语义。调用方可在导入处使用`plugindata`等别名避免与局部变量冲突，但组件自身不得继续保留旧`plugindb`包名或文件前缀。内部`host`和`plan`子组件继续保留在`capability/data/internal/{host,plan}`，只允许`capability/data`公开 facade 暴露必要契约。宿主 datahost 继续只依赖公开 facade，不直接导入其内部实现。

### 8. 迁移验证以编译门禁和静态检索为主

本变更是包边界重构，必须以 Go 编译门禁验证所有调用方导入路径已经同步。重点验证：

- `apps/lina-core/pkg/plugin/...`
- `apps/lina-core/pkg/plugin/capability/data`
- `apps/lina-core/internal/service/plugin/...`
- `apps/lina-core/internal/cmd`
- 官方源码插件后端包
- 动态插件样例和 WASI guest 编译
- `hack/tools/linactl`中的 WASM 构建器和 smoke 测试

同时用静态检索确认旧路径无生产残留：`pkg/pluginservice`、`pkg/pluginhost`、`pkg/pluginbridge`、`pkg/plugindb`、`pkg/sourceupgrade`。

## Risks / Trade-offs

- **迁移面大，导入路径多** → 按“新目录建立、包移动、导入替换、旧目录删除、编译验证”的顺序推进，并优先使用 Go 编译暴露遗漏。
- **`internal`误用导致插件无法编译** → 在设计、任务和审查中显式检查`capability/guest`、`capability/bizctx`、`capability/contract`等插件可导入包不在`internal`下。
- **动态`hostServices`和`capability`语义混淆** → 保留 manifest 字段名不变，在 Go 包和注释中明确`hostServices`是授权 transport 面，`capability`是语义能力目录。
- **`plugindb`需要同时被 guest 与宿主 datahost 导入** → 放入`pkg/plugin/capability/data`公开子包，避免`internal`可见性破坏动态插件和宿主运行时编译。
- **删除`pkg/sourceupgrade`后宿主内部 DTO 归属不清** → 将 DTO 归属到宿主插件服务或内部 sourceupgrade 组件，由`internal/service/plugin`对宿主内部调用方暴露稳定服务方法。
- **规范与代码路径不同步** → 将 OpenSpec 增量规范作为任务项之一，迁移完成后运行`openspec validate refactor-plugin-package-boundaries --strict`并静态检索旧路径。

## Migration Plan

1. 建立`apps/lina-core/pkg/plugin/`目标目录结构，先迁移`pluginbridge`和`pluginhost`到新路径，保留其内部职责不变。
2. 将`pkg/pluginservice`迁移到`pkg/plugin/capability`，同步调整根目录接口命名、包注释、子包路径和测试。
3. 更新宿主插件运行时、WASM host service、启动装配、官方源码插件、动态插件样例和`linactl`构建器的 import 路径。
4. 将`apps/lina-core/pkg/plugindb`迁移到`apps/lina-core/pkg/plugin/capability/data`，同步动态插件样例、宿主 datahost、CI smoke fixture 和`linactl`测试。
5. 删除`apps/lina-core/pkg/sourceupgrade`，把公共 contract 使用点改为宿主插件服务内部类型或内部 sourceupgrade 类型。
6. 删除旧公共目录，不保留兼容 facade；使用静态检索确认生产代码不再 import 旧路径。
7. 运行覆盖变更范围的 Go 编译门禁、动态插件 WASI 编译 smoke、`linactl`测试、OpenSpec 严格校验和格式检查。

回滚策略：本项目无兼容负担，不通过旧路径 facade 回滚。若迁移过程中发现某个能力目录职责划分不清，应修正`pkg/plugin/capability`内部结构后继续，而不是恢复`pkg/pluginservice`。

## Open Questions

暂无需要用户继续确认的阻断问题。用户反馈已确认`plugindb`需要随本轮公共命名空间收敛一并迁移；迁移目标为`pkg/plugin/capability/data`公开子包，保持 data service 授权、typed plan、宿主治理和动态 guest 行为不变。
