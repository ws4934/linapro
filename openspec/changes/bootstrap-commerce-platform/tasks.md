## 1. 基线治理规范确立

- [ ] 1.1 确认`commerce-platform-architecture`治理规范的需求与场景完整、可测
- [ ] 1.2 确认`commerce-extension-primitives-governance`治理规范的需求与场景完整、可测
- [ ] 1.3 核对两份规范与既有治理规范（`core-host-boundary-governance`、`project-positioning-governance`、`module-decoupling`、`plugin-host-domain-capabilities`、`plugin-hook-slot-extension`、`tenant-platform-access-control`）无冲突

## 2. 架构决策与路线记录

- [ ] 2.1 确认`design.md`记录的架构决策 D1–D8（`store = tenant`、storefront host 面、`liquidgo`SSR、6 原语与 Functions、源码/动态分工、主题打包为插件、`tenant-core`二开边界、派生模型）
- [ ] 2.2 确认 P0–P7 路线与原语前置依赖表完整
- [ ] 2.3 记录双 fork 环境与 submodule 指向（`origin`指向各自 fork、`upstream`指向官方，已就位）

## 3. 派生变更待办登记

- [ ] 3.1 登记 P0 派生变更：`tenant-store-domain-resolution`、`commerce-store-profile`
- [ ] 3.2 登记原语派生变更：`plugin-functions-primitive`、`plugin-storefront-contribution-surface`、`plugin-metafields-metaobjects`、`plugin-events-webhooks`
- [ ] 3.3 登记店面与主题派生变更：`storefront-host-surface`、`theme-engine-liquid-runtime`、`theme-visual-editor`
- [ ] 3.4 登记商业域派生变更：`catalog-collection-engine`、`cart-checkout-order`、`customer-identity`
- [ ] 3.5 登记国际化与生态派生变更：`markets-currency`、`content-translation`、`app-marketplace`、`theme-app-extensions`、`admin-ui-extensions`、`merchant-billing`

## 4. 验证与归档准备

- [ ] 4.1 运行`openspec validate bootstrap-commerce-platform --strict`并通过
- [ ] 4.2 记录 DI 来源检查结论：本基线为纯文档与治理，无新增运行期依赖、服务构造或启动装配
- [ ] 4.3 记录影响分析：`i18n`、缓存一致性、数据权限、开发工具跨平台与测试策略在基线层无直接改动，均由后续派生变更在其范围内承担
- [ ] 4.4 `lina-review`审查通过后准备归档，升档两份治理规范为`openspec/specs/`活规范
