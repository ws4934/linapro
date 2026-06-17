## ADDED Requirements

### Requirement: 自定义域名解析

系统 SHALL 提供`domainResolver`，按规范化后的`request.Host`查询域名映射表，仅当命中域名`is_verified = true`且其租户`status = active`时解析到对应`tenant_id`；未命中、未验证或租户停用时 MUST 返回不匹配，且 MUST NOT 返回平台`tenant_id = 0`。

#### Scenario: 已验证域名命中

- **WHEN** 请求`Host`等于某已验证且租户处于`active`的映射域名
- **THEN** `domainResolver`返回该域名映射的`tenant_id`并标记来源为`domain`
- **AND** 不进行按租户`code`的子域名匹配

#### Scenario: 未验证或停用域名不匹配

- **WHEN** 请求`Host`对应的域名映射`is_verified = false`，或其租户`status = suspended`
- **THEN** `domainResolver`返回不匹配
- **AND** 不返回任何`tenant_id`，也不返回平台`0`

### Requirement: host-only 解析不回落平台

系统 SHALL 提供仅包含域名与子域名 resolver、排除`default`的 host-only 解析链。该链在未匹配任何租户时 MUST 返回空结果，MUST NOT 回落平台`tenant_id = 0`，以供顾客匿名 storefront 解析使用。

#### Scenario: host-only 链未匹配

- **WHEN** 以 host-only 链解析的请求`Host`无任何已验证域名或有效子域名匹配
- **THEN** 解析返回空结果
- **AND** 不返回平台`tenant_id = 0`，由调用方按未找到处理

#### Scenario: host-only 链命中

- **WHEN** 以 host-only 链解析的请求`Host`命中一个已验证域名或有效子域名
- **THEN** 解析返回对应`tenant_id`
- **AND** 解析过程不经过`default`回落

### Requirement: 可配置子域名解析

系统 SHALL 从插件配置读取`RootDomain`与保留子域名集合驱动子域名解析。`RootDomain`为空时子域名解析 MUST 保持禁用；非空时按首段 label 匹配租户`code`，且保留子域名 MUST NOT 被解析为租户。

#### Scenario: 配置 RootDomain 后子域名命中

- **WHEN** `RootDomain`配置为有效根域且请求`Host`为`<code>.<RootDomain>`
- **THEN** 子域名解析按 label 匹配租户`code`并返回其`tenant_id`

#### Scenario: 保留子域名不解析

- **WHEN** 请求`Host`首段 label 属于保留子域名（如`www`、`api`、`admin`）
- **THEN** 子域名解析返回不匹配

#### Scenario: RootDomain 为空时禁用

- **WHEN** `RootDomain`为空
- **THEN** 子域名解析对任意`Host`返回不匹配

### Requirement: 域名映射管理与数据权限

系统 SHALL 提供平台作用域的域名映射管理（创建、列表、删除、验证标记）。一个域名 MUST 唯一映射到一个租户；读取类接口 MUST 在数据库查询阶段注入数据权限，创建、删除与验证类操作 MUST 在操作前校验目标域名与租户可见性。

#### Scenario: 创建域名映射

- **WHEN** 平台管理员为某租户创建一个域名映射
- **THEN** 系统校验该域名未被占用后建立映射
- **AND** 重复域名被唯一约束拒绝

#### Scenario: 列表注入数据权限

- **WHEN** 调用方查询域名映射列表
- **THEN** 查询在数据库阶段注入数据权限过滤
- **AND** 不先返回全量再在内存过滤
