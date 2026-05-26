## MODIFIED Requirements

### Requirement:插件配置服务提供通用只读配置访问

系统 SHALL 通过`apps/lina-core/pkg/plugin/capability/config`向源码插件提供业务无关的只读配置访问服务。该服务必须允许源码插件通过任意配置键读取宿主配置文件内容，不得为特定插件或业务模块暴露插件特定的`GetXxx()`配置方法。

#### Scenario:插件读取任意配置键

- **WHEN** 源码插件通过插件配置服务读取现有配置键
- **THEN** 系统返回该键的配置值
- **AND** 配置服务不要求键位于插件特定前缀下

#### Scenario:公共组件不包含插件业务配置方法

- **WHEN** 为源码插件添加或修改私有配置结构
- **THEN** 开发者在插件内部定义配置结构、默认值和验证逻辑
- **AND** 无需在`apps/lina-core/pkg/plugin/capability/config`中添加插件特定的`GetXxx()`方法或插件业务配置类型
