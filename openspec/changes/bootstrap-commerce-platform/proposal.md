## Why

`LinaPro`当前定位是面向可持续交付的 AI 原生全栈框架，默认入口是管理工作台。要在其上二开出类 Shopify 2.0 的多语言、多货币、多插件国际化 SaaS 商城（目标兼容任意 Shopify 主题与 theme app extension），需要跨越 storefront 渲染、主题运行时、商业域、国际化与插件 core 扩展的大规模工程。这种规模必须先有一份程序级基线，把方向、宿主边界与分阶段路线一次性钉死，避免后续每个模块各自决策、污染核心宿主或重复返工。

本基线只确立方向与治理边界，不实现任何商业功能；所有商业能力由后续派生的独立 OpenSpec 变更落地。

## What Changes

- 确立北极星架构：`store = tenant`（复用`linapro-tenant-core`provider）、storefront 作为与 admin 平级的新 host 渲染面、app 与主题与 function 统一收敛在插件 core 之下。
- 确立插件 core 需新增的 6 个扩展原语方向：Functions（typed 计算扩展点，承重墙）、storefront 贡献面、`metafields`/`metaobjects`、events/webhooks、admin UI 扩展 slot、app 安装与计费；并明确这些属于核心宿主能力增强，与「业务逻辑进插件」不冲突。
- 确立主题运行时形态：基于`liquidgo`的 SSR Liquid + 自建 Shopify 主题运行时（Drop 对象模型、Shopify filter/tag、`OS2.0`section 架构、在线可视化装修），并接受按需 fork`liquidgo`。
- 确立 P0–P7 分阶段路线与「原语前置」依赖，以及每个模块以独立 OpenSpec 变更派生、并遵守本基线治理规范的协作模式。
- 确立`linapro-tenant-core`二开边界（resolver 链可扩展、域名解析、公开 storefront 解析）与店铺业务属性另置策略；记录已就位的双 fork 环境（`linapro`与`official-plugins`各自 fork）。

## Capabilities

### New Capabilities

- `commerce-platform-architecture`：商城二开的北极星架构与宿主边界，覆盖`store = tenant`、storefront host 面、插件 core 扩展模型、主题运行时形态，以及 P0–P7 模块分解与派生变更关系。
- `commerce-extension-primitives-governance`：插件 core 6 个扩展原语（以 Functions 为核心）的边界、契约形态与确定性、性能、数据权限要求，约束后续所有原语类派生变更。

### Modified Capabilities

本基线为增量方向确立，不修改任何既有规范的需求。既有治理规范`core-host-boundary-governance`、`project-positioning-governance`、`module-decoupling`、`plugin-host-domain-capabilities`、`plugin-hook-slot-extension`、`tenant-platform-access-control`作为遵循依据被引用，不在本变更中修改。

## Impact

- 改动边界：确立`apps/lina-core`（核心宿主：storefront host 面、插件 core 原语）与`apps/lina-plugins`（`tenant-core`二开、commerce 源插件、第三方动态 app）两条改动线，以及前端`apps/lina-vben`（admin 工作台与可视化装修）的职责划分。
- 外部依赖：引入`liquidgo`（MIT）作为 storefront 模板引擎，可能按需 fork；不在本变更引入运行期代码依赖。
- 治理影响：本基线升档后成为后续所有商城派生变更的约束基线；命中核心宿主边界、插件能力、数据权限、`i18n`、缓存一致性与接口性能等规则域，由各派生变更在其范围内分别满足。
- 不涉及：本变更不改动运行期依赖、服务构造、启动装配、数据库迁移或对外 API；无 DI 来源变化，属纯文档与治理基线。
