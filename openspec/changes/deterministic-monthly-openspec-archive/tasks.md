## 1. Deterministic Archive Workflow

- [x] 1.1 Add a shared monthly OpenSpec deterministic archive composite action.
- [x] 1.2 Wire Codex, Claude Code, and GitHub Copilot reusable workflows to use deterministic archive before AI consolidation.
- [x] 1.3 Keep final failure signaling for completed active changes that remain after partial archive success.
- [x] 1.4 Upgrade artifact upload workflow actions away from the Node 20 runtime generation.

## 2. OpenSpec Archive Blocker Fix

- [x] 2.1 Fix `remove-sqlite-support` so `openspec archive -y remove-sqlite-support` can apply against the current baseline.

## 3. Verification

- [x] 3.1 Run OpenSpec validation for this change and `remove-sqlite-support`.
- [x] 3.2 Run workflow YAML/action validation and shell syntax checks for modified CI files.
- [x] 3.3 Run a temporary-copy deterministic archive smoke that covers partial success and the fixed `remove-sqlite-support` blocker.
- [x] 3.4 Record i18n, cache, data permission, REST API, E2E, and Go production code impact.
- [x] 3.5 Run `lina-review` for the CI/OpenSpec governance change.

## Feedback

- [x] **FB-1**: Monthly OpenSpec archive failed because AI auto archive returned success while completed active changes remained unarchived.
- [x] **FB-2**: `remove-sqlite-support` cannot be archived because its OpenSpec delta removes a requirement header that no longer exists in the baseline spec.
- [x] **FB-3**: Manual monthly OpenSpec archive dispatch is skipped when triggered from a non-default branch.
- [x] **FB-4**: Manual monthly OpenSpec archive dispatch should run against the selected source branch and create a PR back to that same branch.
- [x] **FB-5**: Copilot archive consolidation can produce invalid OpenSpec and block the already validated deterministic archive PR.
- [x] **FB-6**: Archive branch push succeeds but repository policy blocks GitHub Actions from creating the pull request.

## Verification Notes

- FB-1 修复：新增 `.github/actions/monthly-openspec-auto-archive`，使用固定 OpenSpec CLI 版本执行 `openspec list --json`，按名称稳定排序 completed active changes，逐个运行 `openspec archive -y <change>`。当 OpenSpec 提供任务计数且 `completedTasks != totalTasks` 时直接记录失败，不执行归档。每个候选归档后重新执行 `openspec list --json`，确认该 change 已离开活跃列表；即使 OpenSpec CLI 在 `Aborted` 场景下返回 0，也会被识别为失败。
- FB-1 workflow 接入：`.github/workflows/monthly-openspec-archive-{codex,cc,copilot}.yml` 先执行确定性归档，再检测 OpenSpec diff；只有存在 diff 时才准备对应 AI runtime 并执行 archive consolidation。`monthly-openspec-assert-archive-complete` 移到 PR finalization 之后，并仅在 deterministic archive 报告失败时运行，从而允许成功归档部分先写入归档 PR，同时保留失败 job 信号。
- FB-1 prompt 清理：基础自动归档不再通过 AI prompt 执行，已删除废弃的 `.github/prompts/monthly-openspec-auto-archive.zh-CN.md`；保留 archive consolidation prompt 供 AI 聚合阶段使用。
- FB-1 artifact 升级：所有 `actions/upload-artifact@v4` 已升级为 `actions/upload-artifact@v7`，静态扫描确认 `.github` 中不再存在 `upload-artifact@v4`。
- FB-2 修复：将 `remove-sqlite-support` 中与当前主规范不匹配的 REMOVED/MODIFIED header 调整为现有 baseline requirement，删除已经被当前 baseline 吸收且不存在的 SQLite 专属 REMOVED delta；新增 header mismatch 检查确认该变更所有 MODIFIED/REMOVED requirement 标题均存在于当前主规范。
- 验证通过：`openspec validate remove-sqlite-support --strict`。
- 验证通过：`openspec validate deterministic-monthly-openspec-archive --strict`。
- 验证通过：`ruby -e 'require "yaml"; ARGV.each { |f| YAML.load_file(f); puts "ok #{f}" }' .github/workflows/*.yml .github/actions/*/action.yml`。
- 验证通过：`go run github.com/rhysd/actionlint/cmd/actionlint@latest .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml .github/workflows/reusable-e2e-tests.yml .github/workflows/reusable-host-only-build-smoke.yml .github/workflows/reusable-redis-cluster-smoke.yml`。
- 验证通过：对 `remove-sqlite-support` 执行 Node header mismatch 扫描，所有 MODIFIED/REMOVED requirement 标题均匹配当前 `openspec/specs/<capability>/spec.md`。
- 验证通过：在 `/tmp/linapro-action-smoke` 临时副本中抽取 `.github/actions/monthly-openspec-auto-archive/action.yml` 的实际 `run` 脚本执行。结果：`archived-count=8`、`failed-count=0`、`had-failures=false`、`failed-json=[]`、`openspec list --json` 剩余 `changes=[]`、`openspec validate --all` 为 `92 passed, 0 failed`。
- 验证通过：`git diff --check -- .github openspec/changes/deterministic-monthly-openspec-archive openspec/changes/remove-sqlite-support/specs`。
- i18n 影响：本次仅修改 GitHub Actions workflow/composite action、OpenSpec 变更文档和 OpenSpec delta，不新增、修改或删除用户运行时可见文案、前端语言包、宿主/插件 `manifest/i18n` 或 apidoc i18n JSON。
- 缓存一致性影响：本次不修改运行时业务缓存、缓存键、失效触发、分布式协调或跨实例一致性逻辑；新增的 CI action 只在 GitHub runner 工作区内执行 OpenSpec 归档，不涉及生产缓存。
- 数据权限影响：本次不新增或修改 HTTP/API 数据操作接口、服务数据访问路径、插件宿主服务适配器或聚合统计，不影响角色数据权限边界。
- REST API 影响：本次不新增或修改 REST API。
- E2E 影响：本次为 CI/OpenSpec 治理修复，不涉及用户可观察页面、路由、表单、表格或端到端业务流程；使用 OpenSpec、workflow/actionlint、header mismatch 扫描和临时归档 smoke 作为治理验证。
- Go 生产代码影响：本次不新增或修改 Go 生产代码，不触发后端 Go 编译门禁。`actionlint` 通过 `go run` 执行属于外部验证工具，不改变仓库 Go 源码。
- Review：已按 `lina-review` 口径完成审查。审查范围来源包括 `git status --short`、`git ls-files --others --exclude-standard`、`openspec status --change deterministic-monthly-openspec-archive --json`、`openspec status --change remove-sqlite-support --json`、`.github` 与目标 OpenSpec 文件 diff、OpenSpec strict 校验、workflow/action YAML 解析、actionlint、`remove-sqlite-support` header mismatch 扫描和临时 action smoke。确认 `.github/actions/monthly-openspec-auto-archive` 会在归档后复查 active list，能识别 OpenSpec CLI `Aborted` 但退出码为 0 的场景；三条 monthly reusable workflow 均先执行确定性归档，再按需准备 AI runtime 做聚合，且 finalization 后保留剩余 completed active change 的失败信号；`upload-artifact@v4` 已清理。严重问题 0；警告 0。当前工作区仍存在与本次无关的 Go、前端、测试与其他 OpenSpec 改动，本次未修改或回退。
- FB-3 修复：`.github/workflows/monthly-openspec-archive.yml` 的 router `detect` job 现在允许 `workflow_dispatch` 从任意 branch ref 进入，同时继续只让定时触发在默认分支运行。
- FB-3 规范更新：`Manual archive run` 场景明确手动触发可来自任意暴露该 workflow 的 branch ref。
- FB-3 验证通过：`go run github.com/rhysd/actionlint/cmd/actionlint@latest .github/workflows/monthly-openspec-archive.yml`。
- FB-3 验证通过：`ruby -e 'require "yaml"; YAML.load_file(ARGV[0]); puts "ok #{ARGV[0]}"' .github/workflows/monthly-openspec-archive.yml`。
- FB-3 验证通过：`openspec validate deterministic-monthly-openspec-archive --strict`。
- FB-3 验证通过：`git diff --check -- .github/workflows/monthly-openspec-archive.yml openspec/changes/deterministic-monthly-openspec-archive`。
- FB-3 影响评估：本次仅修改 GitHub Actions router 条件和 OpenSpec 文档；不新增或修改 Go 生产代码、前端页面、REST API、运行时 i18n 资源、业务缓存、数据权限逻辑或用户可观察应用流程，因此不触发后端 Go 编译门禁和 E2E 测试。验证方式采用 workflow/actionlint、YAML 解析、OpenSpec strict 校验和 diff whitespace 检查。
- FB-3 Review：已按 `lina-review` 口径完成审查。确认手动触发入口不再受 `github.ref == default_branch` 限制；定时触发仍保留默认分支限制；手动触发限定为 branch ref，避免 tag ref 进入 PR base 语义。严重问题 0；警告 0。
- FB-4 修复：router 新增 `Resolve Target Branch` 步骤，从 `github.ref_name` 解析本次触发分支，输出 `target_branch` 和带安全化源分支标识的 `pr_branch`（格式为 `automation/monthly-openspec-archive-<branch-slug>`）。router 的检测 checkout 改为 `target_branch`，并将 `target_branch`、`pr_branch` 传入 Codex、Claude Code 和 Copilot reusable workflow。
- FB-4 工具 workflow 接入：`.github/workflows/monthly-openspec-archive-{codex,cc,copilot}.yml` 新增 required inputs `target_branch` 和 `pr_branch`；三个 workflow 均 checkout `target_branch` 执行确定性归档和可选聚合，并在 `Finalize Archive Pull Request` 中使用 `base-branch: target_branch`、`pr-branch: pr_branch`。
- FB-4 规范更新：`Manual archive run` 场景明确手动触发分支就是检测、执行和 PR base 分支，PR 来源分支必须包含该触发分支的安全化标识。
- FB-4 影响评估：本次仅修改 GitHub Actions workflow 和 OpenSpec 文档；不新增或修改 Go 生产代码、前端页面、REST API、运行时 i18n 资源、业务缓存、数据权限逻辑或用户可观察应用流程，因此不触发后端 Go 编译门禁和 E2E 测试。验证方式采用 workflow/actionlint、YAML 解析、OpenSpec strict 校验和 diff whitespace 检查。
- FB-4 验证通过：`go run github.com/rhysd/actionlint/cmd/actionlint@latest .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml`。
- FB-4 验证通过：`ruby -e 'require "yaml"; ARGV.each { |f| YAML.load_file(f); puts "ok #{f}" }' .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml`。
- FB-4 验证通过：`openspec validate deterministic-monthly-openspec-archive --strict`。
- FB-4 验证通过：`git diff --check -- .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml openspec/changes/deterministic-monthly-openspec-archive`。
- FB-4 验证通过：Node slug smoke 确认 `john-e2e-enhance -> automation/monthly-openspec-archive-john-e2e-enhance`、`feature/archive test -> automation/monthly-openspec-archive-feature-archive-test`、`release/2026.05 -> automation/monthly-openspec-archive-release-2026.05`。
- FB-4 Review：已按 `lina-review` 口径完成审查。确认手动触发时 workflow 使用触发 branch ref 作为检测 checkout、归档 checkout 和归档 PR base；PR head 分支使用 `automation/monthly-openspec-archive-<branch-slug>`，包含源分支标识并替换非法字符；定时触发仍仅允许默认分支，并自然生成默认分支对应 PR head。严重问题 0；警告 0。
- FB-5 失败分析：GitHub Actions run `26200881048` 在 `john-e2e-enhance` 手动触发，`Run Deterministic OpenSpec Auto Archive` 与 `Validate OpenSpec After Auto Archive` 均成功，`Run Lina Archive Consolidate` 返回成功，但 `Validate OpenSpec After Archive Consolidate` 失败，导致 `Finalize Archive Pull Request` 被跳过。GitHub job log 和 artifact API 当前返回 403，无法读取私有详细日志；从 job step 结论可确认失败范围只在可选 AI 聚合后的 OpenSpec 校验。
- FB-5 修复：`.github/workflows/monthly-openspec-archive-{codex,cc,copilot}.yml` 在确定性归档校验通过后创建 `openspec` 快照；`Run Lina Archive Consolidate` 和聚合后校验改为 `continue-on-error: true`；当 AI 聚合命令失败或聚合后 OpenSpec 校验失败时，workflow 记录 `git status`、`git diff -- openspec` 和 `openspec validate --all` 日志到对应 AI 日志目录，随后恢复确定性归档快照并重新校验，通过后继续执行 `Finalize Archive Pull Request`。
- FB-5 日志复核：用户提供的 `/Users/john/Downloads/job-logs.txt` 确认 Copilot 在聚合阶段创建了临时活跃变更 `openspec/changes/archive-consolidation`，但最终没有清理；`openspec validate --all` 输出 `✗ change/archive-consolidation`，总计 `92 passed, 1 failed (93 items)`。已新增 `Check Archive Consolidate Temporary Change Cleanup` 步骤，显式检测该临时变更残留并触发诊断与回退。
- FB-5 规范更新：AI archive consolidation 被定义为可选增强阶段；失败或产生无效 OpenSpec 时不得阻塞已经通过校验的 deterministic archive PR，必须记录诊断、回滚聚合结果并继续归档。
- FB-5 影响评估：本次仅修改 GitHub Actions workflow 和 OpenSpec 文档；不新增或修改 Go 生产代码、前端页面、REST API、运行时 i18n 资源、业务缓存、数据权限逻辑或用户可观察应用流程，因此不触发后端 Go 编译门禁和 E2E 测试。验证方式采用 workflow/actionlint、YAML 解析、OpenSpec strict 校验和 diff whitespace 检查。
- FB-6 失败分析：GitHub Actions run `26203378371` 已成功 push `automation/monthly-openspec-archive-john-e2e-enhance`，但 `gh pr create` 返回 `GraphQL: GitHub Actions is not permitted to create or approve pull requests (createPullRequest)`，说明仓库 Actions 设置不允许 `GITHUB_TOKEN` 创建 PR。
- FB-6 修复：`.github/actions/monthly-openspec-finalize-pr/action.yml` 在 `gh pr create` 或 `gh pr edit` 失败时检查该仓库策略错误；若归档分支已经成功 push，则输出 base/head 分支和手动 PR URL 到日志与 step summary，并让 workflow 成功结束。其他 PR 命令错误仍按失败处理。
- FB-6 规范更新：新增仓库策略阻止 PR 创建场景，明确已推送有效归档分支时 workflow 应输出手动 PR 链接并成功结束。
- FB-6 影响评估：本次仅修改 GitHub Actions composite action 和 OpenSpec 文档；不新增或修改 Go 生产代码、前端页面、REST API、运行时 i18n 资源、业务缓存、数据权限逻辑或用户可观察应用流程，因此不触发后端 Go 编译门禁和 E2E 测试。验证方式采用 actionlint、YAML 解析、OpenSpec strict 校验、diff whitespace 检查和 shell 语法检查。
- FB-6 验证通过：`go run github.com/rhysd/actionlint/cmd/actionlint@latest .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml`。
- FB-6 验证通过：`ruby -e 'require "yaml"; ARGV.each { |f| YAML.load_file(f); puts "ok #{f}" }' .github/actions/monthly-openspec-finalize-pr/action.yml .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml`。
- FB-6 验证通过：`openspec validate deterministic-monthly-openspec-archive --strict`。
- FB-6 验证通过：`git diff --check -- .github/actions/monthly-openspec-finalize-pr/action.yml .github/workflows/monthly-openspec-archive.yml .github/workflows/monthly-openspec-archive-codex.yml .github/workflows/monthly-openspec-archive-cc.yml .github/workflows/monthly-openspec-archive-copilot.yml openspec/changes/deterministic-monthly-openspec-archive`。
- FB-6 验证通过：从 `.github/actions/monthly-openspec-finalize-pr/action.yml` 提取 `Create or Update Archive Pull Request` 的 bash 脚本并执行 `bash -n`。
- FB-6 验证通过：模拟 `GraphQL: GitHub Actions is not permitted to create or approve pull requests (createPullRequest)` 输出，确认策略错误匹配分支可识别该错误。
- FB-6 Review：已按 `lina-review` 口径完成审查。确认归档分支 push 仍是硬失败边界；只有 push 成功后的 PR 创建/更新被仓库策略拒绝时才降级为手动 PR 链接；其他 `gh pr` 错误仍会 `exit 1`。严重问题 0；警告 0。
