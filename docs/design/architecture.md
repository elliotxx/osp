# OSP 架构设计

## 整体架构

OSP 采用分层架构设计，主要包含以下几个层次：

1. 命令行接口层（CLI Layer）
2. 业务逻辑层（Business Layer）
3. 数据访问层（Data Access Layer）

### 命令行接口层

命令行接口层使用 [cobra](https://github.com/spf13/cobra) 库实现，主要包含以下命令：

- `auth`: 认证管理
  - `login`: 登录 GitHub
  - `status`: 查看认证状态
  - `logout`: 登出 GitHub

- `repo`: 仓库管理
  - `add`: 添加仓库
  - `remove`: 移除仓库
  - `list`: 列出仓库
  - `switch`: 切换仓库
  - `current`: 查看当前仓库

- `plan`: 规划管理
  - 生成和更新规划文档
  - 支持自定义标签和分类
  - 支持排除 PR

### 业务逻辑层

业务逻辑层包含以下主要模块：

1. `auth`: 认证管理
   - 使用 GitHub CLI 的认证机制
   - 管理令牌的获取和验证

2. `repo`: 仓库管理
   - 管理多个 GitHub 仓库
   - 维护仓库列表和当前仓库

3. `planning`: 规划管理
   - 里程碑信息获取
   - 规划文档生成
   - 进度统计和可视化

### 数据访问层

数据访问层主要包含：

1. `config`: 配置管理
   - 基于 YAML 的配置文件
   - 支持多平台配置路径

2. `api`: GitHub API 访问
   - 使用 GitHub API v3
   - 支持 REST 和 GraphQL

## 数据流

1. 用户通过命令行输入命令
2. CLI 层解析命令和参数
3. 业务逻辑层处理请求
4. 数据访问层与 GitHub API 交互
5. 结果返回给用户

## 依赖关系

主要的外部依赖：

1. `github.com/spf13/cobra`: 命令行框架
2. `github.com/cli/go-gh`: GitHub CLI 工具库
3. `gopkg.in/yaml.v3`: YAML 配置解析

## 扩展性设计

1. 接口抽象
   - 所有核心功能都通过接口定义
   - 便于替换具体实现

2. 模块化设计
   - 每个功能都是独立的模块
   - 可以方便地添加新功能

3. 配置驱动
   - 大部分行为可通过配置文件调整
   - 支持自定义模板和规则
