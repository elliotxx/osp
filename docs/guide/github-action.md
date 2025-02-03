# GitHub Action 使用指南

本指南详细介绍了如何通过 GitHub Action 自动化使用 OSP。

## 目录
- [快速开始](#快速开始)
- [使用场景](#使用场景)
  - [新手任务自动化](#新手任务自动化)
  - [项目规划自动化](#项目规划自动化)
  - [数据统计自动化](#数据统计自动化)
- [配置说明](#配置说明)
- [最佳实践](#最佳实践)

## 快速开始

1. 在你的仓库中创建 `.github/workflows/` 目录
2. 根据需要创建相应的工作流文件（如 `community-task-updater.yml`）
3. 配置工作流权限：
   - 进入仓库的 Settings -> Actions -> General
   - 在 "Workflow permissions" 中启用 "Read and write permissions"
   - 保存配置更改
4. 合并到主干分支后工作流将自动生效

## 使用场景

### 新手任务自动化

自动生成和更新社区新手任务列表，当 Issue 发生变化时自动更新：

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
            --difficulty-labels 'difficulty/easy,difficulty/medium,difficulty/hard'
            --category-labels 'bug,documentation,enhancement'
            --target-title '社区新手任务 | Community Tasks 🎯'
```

### 项目规划自动化

自动生成和更新项目里程碑规划，当里程碑或相关 Issue 发生变化时自动更新：

```yaml
name: Community Task Updater

on:
  # Manually triggered
  workflow_dispatch:
  # Trigger on issue events
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
          # Optional: version of OSP to use (default: latest)
          version: 'latest'
          
          # Optional: working directory (default: project root)
          working-directory: '.'
          
          # Optional: GitHub token (default: ${{ github.token }})
          github-token: ${{ secrets.GITHUB_TOKEN }}
          
          # Optional: enable debug mode (default: false)
          debug: false
          
          # Optional: skip caching (default: false)
          skip-cache: false
          
          # Optional: additional OSP arguments
          args: >-
            plan
            --yes
            --category-labels bug,documentation,enhancement
```

### 数据统计自动化 (TODO)

定期更新项目统计数据，包括 Issue、PR、Star 等数据：

```yaml
name: Project Stats Updater

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 * * *'  # 每天更新

jobs:
  osp-run:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Update Project Stats
        uses: elliotxx/osp-action@main
        with:
          version: 'latest'
          args: >-
            stats
            --yes
            --output-format markdown
            --output-file 'docs/stats/README.md'
```

## 配置说明

### 输入参数

osp-action 支持以下输入参数：

| 参数名 | 说明 | 必填 | 默认值 |
|--------|------|------|---------|
| version | OSP 版本 | 否 | latest |
| working-directory | 工作目录 | 否 | . |
| github-token | GitHub Token | 否 | ${{ github.token }} |
| debug | 调试模式 | 否 | false |
| skip-cache | 跳过缓存 | 否 | false |
| args | OSP 命令参数 | 是 | - |

### 工作流权限

1. 进入仓库的 Settings -> Actions -> General
2. 在 "Workflow permissions" 中启用 "Read and write permissions"
3. 保存配置更改

## 最佳实践

### 1. 事件触发

根据实际需求选择合适的触发事件：
- `workflow_dispatch`: 支持手动触发，便于调试和临时更新
- `schedule`: 定时触发，适合定期更新的场景
- `issues`/`pull_request`: 监听特定事件，实时响应变化
- `milestone`: 监听里程碑变化，用于项目规划更新

### 2. 缓存优化

使用 GitHub Actions 缓存加速执行：

```yaml
jobs:
  osp-run:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Cache OSP data
        uses: actions/cache@v3
        with:
          path: ~/.osp
          key: osp-${{ runner.os }}-${{ hashFiles('**/*.yml') }}
      
      - name: Run OSP
        uses: elliotxx/osp-action@main
        with:
          args: 'onboard --yes'
```