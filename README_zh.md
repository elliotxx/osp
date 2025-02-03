![OSP](https://socialify.git.ci/elliotxx/osp/image?font=Raleway&language=1&name=1&owner=1&pattern=Plus&theme=Light)

# OSP - Open Source Software Pilot

[![Go Report Card](https://goreportcard.com/badge/github.com/elliotxx/osp)](https://goreportcard.com/report/github.com/elliotxx/osp)
[![GoDoc](https://godoc.org/github.com/elliotxx/osp?status.svg)](https://godoc.org/github.com/elliotxx/osp)
[![License](https://img.shields.io/github/license/elliotxx/osp.svg)](https://github.com/elliotxx/osp/blob/main/LICENSE)

[English](README.md) | [简体中文](README_zh.md)

OSP (Open Source Software Pilot) 是一款专注于开源社区治理的自动化管理工具。它融合了多种开源社区治理的最佳实践，为开源项目维护者提供了一套完整的工具链，助力项目高效运营、精准跟踪、数据驱动决策。

OSP 提供两种使用方式：CLI 命令行工具和 GitHub Action 自动化工作流。CLI 工具适合本地交互式管理，而 GitHub Action 则可以通过订阅事件实现全自动化运维，一次配置，持续服务。

## ✨ 特性

### 已实现功能
- 🔑 GitHub 认证管理 - 安全可靠的身份认证，GITHUB CLI 同款
- 📊 项目数据统计 - 多维度的数据分析
- 📝 新手任务、项目规划生成 - 支持通过订阅 Github 事件自动化更新
- 📈 Star 趋势统计 - 项目增长数据追踪

### 开发路线
- 📋 Roadmap 生成 - 支持通过订阅 Github 事件自动化更新
- 📅 社区动态聚合 - 自动聚合近期评论、新建 PR/Issue/Discussion，支持通过 webhook 订阅长期未响应的社区动态
- 🤖 智能 PR Review - 基于 LLM 的代码审查，支持自动化评论
- 💡 智能 Issue 创建 - 一句话生成 Issue，提升创建任务/需求的效率
- 🔌 GitHub App 集成 - 更强大的集成能力

## 🚀 安装

更多安装方式请参考 [高级安装指南](docs/guide/advanced-installation.md)。

### 🐙 Go 安装

```bash
go install github.com/elliotxx/osp@latest
```

### 🍺 Homebrew 安装

```bash
brew tap elliotxx/tap
brew install osp
```

## 🚀 使用方法

### 🖥️ CLI 命令行

1. 配置 GitHub 认证
```bash
# 使用 GitHub CLI 登录
gh auth login

# 验证认证状态
osp auth status
```

2. 项目管理
```bash
# 添加项目
osp repo add owner/repo

# 切换项目
osp repo

# 查看当前项目
osp repo current
```

3. 核心功能
```bash
# 生成版本规划
osp plan

# 管理新手任务
osp onboard

# 查看项目统计
osp stats

# 分析 Star 趋势
osp stats star-history
```

更多详细说明请参考 [CLI 使用指南](docs/guide/cli.md)。

### 🤖 GitHub Action

> osp-action 实现请查看 [osp-action](https://github.com/elliotxx/osp-action)

下面以通过 osp-action 实现**社区新手任务自动化生成和更新**为例，更多 osp-action 的自动化使用场景请查看 [文档](docs/guide/github-action.md)。

1. 在主分支（main/master）创建工作流配置文件 `.github/workflows/community-task-updater.yml`：
```yaml
name: Community Task Updater

on:
  # 手动触发
  workflow_dispatch:
  # Issue 相关事件触发
  issues:
    types: [opened, edited, deleted, transferred, milestoned, demilestoned, labeled, unlabeled, assigned, unassigned]

jobs:
  osp-run:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Update Community Tasks
        uses: elliotxx/osp-action@main
        with:
          # 可选：指定 OSP 版本，默认使用最新版
          version: 'latest'
          
          # 可选：指定工作目录，默认为项目根目录
          working-directory: '.'
          
          # 可选：GitHub Token，默认使用 GITHUB_TOKEN
          github-token: ${{ secrets.GITHUB_TOKEN }}
          
          # 可选：启用调试模式
          debug: false
          
          # 可选：跳过缓存
          skip-cache: false
          
          # OSP 命令参数
          args: >-
            onboard
            --yes
            --onboard-labels 'help wanted,good first issue'
            --difficulty-labels 'good first issue,help wanted'
            --category-labels bug,documentation,enhancement
            --target-title '社区新手任务 🎯'
```

2. 配置必要权限
- 导航至 Settings -> Actions -> General
- 在 "Workflow permissions" 中启用 "Read and write permissions"
- 保存配置更改

3. 使用方式
- 自动执行：设定的 Github 事件触发时工作流会自动执行
- 手动触发：
  1. 进入 Actions 页面
  2. 选择 "Community Task Updater"
  3. 点击 "Run workflow"

## 📚 文档

- [使用指南](docs/guide/README.md) - 详细的使用说明
- [设计文档](docs/design/README.md) - 架构设计与实现
- [CLI 手册](docs/cli/osp.md) - 命令行工具参考

## 🤝 贡献

我们欢迎各种形式的贡献！无论是新功能、文档改进还是 bug 修复。详情请参考[贡献指南](CONTRIBUTING.md)。

## 📄 许可证

本项目采用 MIT 许可证，查看 [LICENSE](LICENSE) 了解详情。
