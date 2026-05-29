## Tasks

- [x] 更新插件配置规范，明确源码插件宿主配置读取不受 key 白名单限制。
- [x] 移除源码插件宿主配置适配器中的公开 key 白名单并更新注释。
- [x] 增加后端单元测试覆盖任意宿主配置 key 读取。
- [x] 运行 Go 编译/单元测试和 OpenSpec 严格校验。

## Feedback

- [x] **FB-1**: 源码插件读取宿主配置被公开 key 白名单错误限制

### FB-1 处理记录

- 根因：`apps/lina-core/pkg/plugin/capability/hostconfig` 的 `valueForKey` 通过固定 `switch` 只允许 `workspace.basePath`、`i18n.default`、`i18n.enabled` 三个宿主配置 key，其他 key 返回 `host config key is not public`。这使源码插件的 `HostServices.HostConfig()` 被错误套用了公开 key 白名单。
- 修复：移除固定 key 白名单，源码插件 `HostConfigService` 通过启动期注入的宿主配置服务窄 `GetRaw(ctx, key)` 能力按调用 key 读取 GoFrame 宿主配置；空 key 或 `.` 按 GoFrame 配置组件语义返回完整配置快照。同步更新 `HostConfig` 契约和动态插件注释，明确动态插件仍由 `hostServices.resources.keys` 授权快照约束。
- 修改文件：`apps/lina-core/internal/service/config/config_raw.go`、`apps/lina-core/pkg/plugin/capability/hostconfig/hostconfig_access.go`、`apps/lina-core/pkg/plugin/capability/hostconfig/hostconfig.go`、`apps/lina-core/pkg/plugin/capability/hostconfig/hostconfig_access_test.go`、`apps/lina-core/pkg/plugin/capability/contract/config.go`、`apps/lina-core/pkg/plugin/capability/**/*.go`、`apps/lina-core/pkg/plugin/pluginbridge/internal/hostservice/*.go`、`apps/lina-core/internal/service/plugin/internal/{hostservices,wasm,runtimeupgrade,catalog}`、`openspec/changes/remove-source-plugin-host-config-whitelist/**`。
- 影响分析：无 HTTP API、前端 UI、SQL、DAO、数据权限、列表/详情/导出/聚合接口或运行时配置写入影响；源码插件宿主配置读取能力由受白名单限制改为只读无限制读取；动态插件授权模型保持 `resources.keys` 校验不变。
- `i18n` 影响：无运行时用户可见文案、菜单、API 文档源文本、错误消息、插件清单或语言包资源变更。
- 缓存一致性影响：只读读取宿主静态配置，不新增缓存、失效、刷新、热更新或跨节点一致性机制。
- 数据权限影响：不新增数据库读写、业务数据可见性、租户/组织边界或下载/导出逻辑。
- 开发工具跨平台影响：不修改脚本、CI、构建、测试 runner、`linactl` 或跨平台入口。
- DI 来源检查：不新增运行期依赖；`HostConfigService` 继续由宿主启动期显式构造并传入源码插件 host services 和 WASM host service 配置入口，复用原有 `config.Service` 实例，并通过该实例的 `GetRaw(ctx, key)` 窄能力读取宿主配置，避免在插件适配器中构造独立服务图。
- 已读取规则文件：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/architecture.md`、`.agents/rules/plugin.md`、`.agents/rules/backend-go.md`、`.agents/rules/testing.md`、`.agents/rules/documentation.md`。
- 验证：`go test ./internal/service/config -count=1`；`go test ./pkg/plugin/capability/hostconfig -count=1`；`go test ./pkg/plugin/capability ./pkg/plugin/capability/contract ./pkg/plugin/capability/guest -count=1`；`go test ./internal/service/plugin/internal/hostservices ./internal/service/plugin/internal/wasm ./pkg/plugin/pluginbridge/internal/hostservice -count=1`；`go test ./internal/cmd -count=1`；`openspec validate remove-source-plugin-host-config-whitelist --strict`；`git diff --check`。
