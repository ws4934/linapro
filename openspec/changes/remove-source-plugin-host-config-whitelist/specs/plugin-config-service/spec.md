## MODIFIED Requirements

### Requirement:宿主公开配置必须通过独立服务读取

系统 SHALL 通过 `HostServices.HostConfig()` 向源码插件暴露宿主配置只读读取能力。源码插件通过该服务读取宿主配置时不得受公开 key 白名单限制；空 key 或 `.` MUST 按宿主配置组件语义返回完整配置快照。该服务不得提供写入、保存、热重载或运行时修改宿主配置的方法。动态插件通过 `hostconfig.get` 读取宿主配置时，仍 MUST 先通过 `hostServices` 授权快照校验对应 key，动态插件不得绕过 manifest 授权读取宿主配置。

#### Scenario:源码插件读取任意宿主配置键

- **WHEN** 源码插件通过 `HostServices.HostConfig()` 读取宿主配置键 `database.default.link`
- **THEN** 系统按宿主当前配置源返回该键的配置值
- **AND** 该读取不要求 key 预先登记到公开白名单

#### Scenario:源码插件读取缺失宿主配置键

- **WHEN** 源码插件通过 `HostServices.HostConfig()` 读取不存在的宿主配置键
- **THEN** 系统返回未找到语义
- **AND** 不因 key 未登记到白名单而返回权限错误

#### Scenario:源码插件读取完整宿主配置快照

- **WHEN** 源码插件通过 `HostServices.HostConfig()` 读取空 key 或 `.`
- **THEN** 系统返回宿主当前配置源中的完整配置快照
- **AND** 该读取不要求逐个 key 预先登记到公开白名单

#### Scenario:动态插件宿主配置读取仍受授权快照限制

- **WHEN** 动态插件通过 `hostconfig.get` 读取宿主配置键
- **THEN** 宿主先按当前 release 的 `hostServices` 授权快照校验该 key
- **AND** 未授权 key 的读取必须被拒绝
