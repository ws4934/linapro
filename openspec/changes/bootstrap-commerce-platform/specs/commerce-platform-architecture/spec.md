## ADDED Requirements

### Requirement: 商城商户建模为租户

系统 SHALL 将每个商城商户（store）建模为一个租户，复用`linapro-tenant-core`provider 与行级`tenant_id`隔离，不得另建与租户并行的 store 身份与隔离体系。店铺业务属性 SHALL 落在独立 commerce 插件或`tenant-core`的`config_override`，keyed by`tenant_id`，不得污染`tenant-core`的租户主体模型。

#### Scenario: 新建商户

- **WHEN** 平台新建一个商城商户
- **THEN** 系统创建对应租户，并使该商户全部业务数据按`tenant_id`行级隔离
- **AND** 不创建独立于租户体系的第二套商户身份

#### Scenario: 店铺属性归属

- **WHEN** 需要为商户保存货币、套餐、主题绑定或品牌等店铺业务属性
- **THEN** 这些属性存放在独立 commerce 插件或`config_override`，并以`tenant_id`关联
- **AND** 不向`tenant-core`租户主体表新增商业语义字段

### Requirement: storefront 作为独立宿主渲染面

系统 SHALL 将 storefront 作为与 admin 工作台平级的独立 host 渲染面，具备公开访问、顾客鉴权、按域名路由、SSR 渲染与多维缓存能力。storefront MUST NOT 复用带 staff 鉴权的`/x/<plugin-id>`插件路由来承载顾客访问，顾客身份 MUST 与后台 staff 身份相互隔离。

#### Scenario: 顾客匿名访问店面

- **WHEN** 未登录顾客通过商户域名访问店面页面
- **THEN** 请求经 storefront 面的公开解析链按域名定位到对应租户
- **AND** 不要求后台 staff 登录，也不进入带 staff 鉴权的插件路由

#### Scenario: 顾客与 staff 身份隔离

- **WHEN** 系统签发顾客会话
- **THEN** 顾客身份域与`sys_user`后台 staff 身份相互独立
- **AND** 顾客凭据不得获得后台管理权限

### Requirement: 商业能力归属插件不污染核心宿主

系统 SHALL 将商业域能力（catalog、cart、order、customer、discount、pricing 等）落为插件实现，不纳入`apps/lina-core`核心领域；核心宿主仅提供通用宿主能力与扩展原语。商业引擎 SHALL 采用源码插件，第三方应用 SHALL 采用动态 WASM 插件。

#### Scenario: 新增商业模块

- **WHEN** 团队规划新增一个商业域模块
- **THEN** 该模块以源码插件或动态插件实现，并通过受治理契约接入宿主
- **AND** 不直接扩展`lina-core`核心领域契约、通用 service 语义或存储模型

#### Scenario: 判断能力是否属于核心

- **WHEN** 某能力被多个模块复用并承担框架级统一治理职责
- **THEN** 系统将其作为宿主能力或扩展原语保留在核心
- **AND** 仅服务单一商业场景的能力不进入核心宿主

### Requirement: 主题运行时形态

系统 SHALL 采用 SSR Liquid 渲染主题，以`liquidgo`为语言层、宿主`theme-engine`为渲染引擎。主题包 SHALL 作为插件 artifact 类型分发，复用插件 catalog、安装、版本与按租户启停；渲染引擎 SHALL 留在宿主，插件 MUST NOT 在插件内部执行 Liquid 渲染。

#### Scenario: 上传主题

- **WHEN** 商户上传一个主题包
- **THEN** 系统将其作为插件 artifact 按租户存储并版本化
- **AND** 主题文件经宿主`theme-engine`的按租户`FileSystem`读取

#### Scenario: 应用注入店面

- **WHEN** 第三方应用需要向主题注入 section 或 block
- **THEN** 注入经 storefront 贡献面完成，应用只提供 schema、数据端点与 Liquid 片段
- **AND** Liquid 渲染由宿主执行，应用不渲染 Liquid

### Requirement: storefront 租户隔离红线

系统 SHALL 在 storefront 路径强制租户隔离。平台`tenant_id = 0`的 bypass 语义 MUST NOT 作用于顾客请求；未解析到有效租户的顾客请求 MUST 被拒绝或返回 404，不得落入平台上下文导致跨租户数据暴露。

#### Scenario: storefront 请求未解析到租户

- **WHEN** 顾客请求的域名无法解析到任一有效租户
- **THEN** 系统拒绝该请求或返回 404
- **AND** 不以`tenant_id = 0`平台上下文继续处理该请求

### Requirement: 商城模块以独立派生变更落地

系统 SHALL 使每个商城模块以独立 OpenSpec 变更派生，并遵守本基线确立的治理规范；本基线 MUST NOT 直接实现任何商业功能。

#### Scenario: 实现某商城模块

- **WHEN** 团队开始实现某个商城模块（如域名解析、Functions 原语、storefront 渲染面）
- **THEN** 新建一个独立 OpenSpec 变更并 ADD 其自有 capability 规范
- **AND** 该变更显式遵守`commerce-platform-architecture`与`commerce-extension-primitives-governance`
