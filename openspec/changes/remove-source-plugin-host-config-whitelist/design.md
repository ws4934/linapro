## 决策：源码插件 HostConfig 不再使用 key 白名单

源码插件是宿主构建产物的一部分，其信任边界不同于动态插件。`HostServices.HostConfig()` 对源码插件应表达“只读读取宿主配置”的能力，而不是“读取少量公开配置”的能力。实现上保留 `HostConfigService` 只读接口，移除 `valueForKey` 中的固定 key switch 和根节点拒绝逻辑，改为按调用方提供的 key 直接读取 GoFrame 宿主配置；空 key 或 `.` 按 GoFrame 语义返回完整配置快照。

## 决策：动态插件授权模型保持不变

动态插件仍通过 `pluginbridge` 的 `hostServices` 声明和授权快照访问 `hostconfig.get`。本次只移除宿主配置适配器内部的公开 key 白名单；动态插件调用在到达适配器前仍由 `hostCallContext.hasHostServiceAccess` 校验 `resources.keys`。因此动态插件不能借由本次变更绕过 manifest 授权。

## 影响分析

- `i18n` 影响：无运行时用户可见文案、菜单、API 文档源文本、插件清单或语言包资源变更。
- 缓存一致性影响：读取静态宿主配置，不新增缓存、失效、热更新或跨节点一致性机制。
- 数据权限影响：不新增数据库读写、列表、详情、导出、聚合、下载或租户/组织可见性逻辑。
- 开发工具跨平台影响：不修改脚本、构建、CI、`linactl` 或跨平台入口。
- DI 影响：不新增运行期依赖；`HostConfigService` 仍由启动期传入并复用原有宿主配置服务实例，源码插件配置读取通过该实例的窄 `GetRaw(ctx, key)` 能力进入宿主配置源。
- 测试策略：新增 `hostconfig` 单元测试覆盖未预置白名单 key 可读取、根配置可读取、缺失 key 返回不存在；运行变更包 Go 测试和 OpenSpec 严格校验。
