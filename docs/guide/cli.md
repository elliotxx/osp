# CLI 使用指南

本指南详细介绍了 OSP CLI 工具的使用方法。

## 目录
- [认证配置](#认证配置)
- [仓库管理](#仓库管理)
- [核心功能](#核心功能)
  - [项目规划](#项目规划)
  - [新手任务](#新手任务)
  - [数据统计](#数据统计)

## 认证配置

### 前提条件
- 拥有 GitHub 账号

### 工作原理
OSP 使用 GitHub CLI 的认证机制，但认证信息和 GITHUB CLI 是隔离的，独立存储，互不影响。认证信息会被安全地存储在本地系统的凭证管理器中：
- macOS: Keychain
- Linux: Secret Service API/libsecret
- Windows: Windows Credential Manager

### 使用方法

1. 使用 GitHub CLI 登录
```bash
# 使用 GitHub CLI 登录
gh auth login

# 验证认证状态
osp auth status
```

2. 使用 Token 登录
```bash
# 使用环境变量
export GITHUB_TOKEN=your_token
# or export GH_TOKEN=your_token

# 查看登录状态
osp auth status
```

## 仓库管理

### 前提条件
- 无

### 工作原理
OSP 将仓库信息存储在本地状态目录中（$XDG_STATE_HOME），支持管理多个仓库，但同一时间只能操作一个活跃仓库。

### 使用方法

#### 添加仓库
```bash
# 添加单个仓库
osp repo add owner/repo
```

#### 切换仓库
```bash
# 交互式切换
osp repo

# 直接切换
osp repo switch owner/repo
```

#### 查看仓库
```bash
# 查看当前选中的仓库
osp repo current

# 列出所有仓库
osp repo list
```

## 核心功能

### 项目规划

#### 前提条件
- 已完成仓库配置
- 仓库中已创建 Milestone 并关联 Issue
- 对仓库有写入权限（用于更新规划文档）

#### 工作原理
OSP 会扫描仓库中的所有 Issue，根据标签和里程碑信息自动生成项目规划文档。生成的文档会作为一个新的 Issue 或更新现有的规划 Issue。具体步骤：
1. 获取仓库所有 Issue
2. 根据里程碑和标签分类整理
3. 生成规划文档
4. 创建或更新规划 Issue

#### 使用方法
```bash
# 基础用法，默认会先预览生成的内容，确认后才会更新到远端
osp plan

# 指定里程碑
osp plan 1

# 自定义分类标签
osp plan --category-labels bug,enhancement,documentation

# 自定义优先级标签
osp plan --priority-labels priority/high,priority/medium,priority/low

# 自定义目标 Issue 标题
osp plan --target-title "Planning: {{ .Title }}"

# 排除 PR
osp plan --exclude-pr

# 模拟执行，不会更新任何内容
osp plan --dry-run

# 自动确认
osp plan --yes
```

### 新手任务

#### 前提条件
- 已完成仓库配置
- 仓库中已创建适合新手的 Issue 并设置了特定标签
- 对仓库有写入权限（用于更新任务列表）

#### 工作原理
OSP 通过以下步骤生成新手任务列表：
1. 扫描仓库中带有特定标签的 Issue（如 "good first issue"）
2. 分析 Issue 的难度（根据标签）
3. 按类别组织任务（如 bug 修复、文档改进等）
4. 生成任务列表文档
5. 创建或更新任务列表 Issue

#### 使用方法
```bash
# 基础用法，默认会先预览生成的内容，确认后才会更新到远端
osp onboard

# 自定义新手任务标签
osp onboard --onboard-labels "help wanted,good first issue"

# 自定义难度等级标签
osp onboard --difficulty-labels "difficulty/easy,difficulty/medium,difficulty/hard"

# 自定义分类标签
osp onboard --category-labels bug,enhancement,documentation

# 自定义目标 Issue 标签
osp onboard --target-label getting-started

# 自定义目标 Issue 标题
osp onboard --target-title "社区新手任务"

# 模拟执行，不会更新任何内容
osp onboard --dry-run

# 自动确认
osp onboard --yes
```

### 数据统计

#### 前提条件
- 已完成仓库配置
- 对仓库有读取权限

#### 工作原理
OSP 通过 GitHub API 收集以下数据：
1. 基础统计：Issue 数量、PR 数量、贡献者数量等
2. Star 历史：通过 Star 事件 API 获取增长趋势

#### 使用方法
```bash
# 基础统计
osp stats

# Star 历史
osp star history
```

## 全局选项

所有命令都支持以下选项：

```bash
--no-color   # 禁用彩色输出
-v, --verbose # 详细输出
-V, --version # 显示版本信息
