# OSP - Open Source Software Pilot

[![Go Report Card](https://goreportcard.com/badge/github.com/elliotxx/osp)](https://goreportcard.com/report/github.com/elliotxx/osp)
[![GoDoc](https://godoc.org/github.com/elliotxx/osp?status.svg)](https://godoc.org/github.com/elliotxx/osp)
[![License](https://img.shields.io/github/license/elliotxx/osp.svg)](https://github.com/elliotxx/osp/blob/main/LICENSE)

OSP (Open Source Software Pilot) 是一个自动化的开源软件社区管理工具，内置多种开源社区管理的最佳实践，帮助开源项目维护者更高效地管理项目、跟踪进展、生成报告。

OSP 有两种形态： CLI 和 Github Action。CLI 适应于本地主动维护管理，而 Action 可以通过订阅 Github 事件实现自动化管理，一键配置，自动维护。

## 特性

- [x] 🔑 GitHub 认证管理
- [x] 📊 项目数据统计和分析
- [x] 📝 自动生成项目规划，支持动态更新
- [x] 📝 自动生成新手任务，支持动态更新
- [x] 📈 Star 趋势统计
- [ ] 📝 自动生成 Roadmap，支持动态更新
- [ ] 📅 聚合社区近期动态
- [ ] 📝 基于 LLM 的 PR Review，支持自动化评论
- [ ] 📝 基于 LLM 的一句话创建 Issue
- [ ] 🤖 Github App 集成

## 🚀 安装

更多安装方式请参考 [高级安装指南](docs/guide/advanced-installation.md)。

### 🐙 Go 安装

```bash
go install github.com/elliotxx/osp@latest
```

### 🍺 Homebrew 安装

通过 Homebrew 安装：
```bash
brew tap elliotxx/tap
brew install osp
```

## 使用方法

### 🖥️ 本地安装使用

1. 登录 GitHub
```bash
# 使用 GitHub CLI 登录
gh auth login

# 验证 OSP 认证状态
osp auth status
```

2. 管理仓库
```bash
# 添加仓库
osp repo add owner/repo

# 切换仓库
osp repo

# 查看当前仓库
osp repo current
```

3. 使用功能
```bash
# 生成项目版本规划
osp plan

# 生成项目新手任务
osp onboard

# 查看项目统计
osp stats

# 查看 Star 趋势
osp stats star-history
```

更多使用说明请参考 [CLI 使用文档](docs/guide/cli.md)。

### 🤖 GitHub Action 使用

> osp-action 的代码仓库见 [osp-action](https://github.com/elliotxx/osp-action)

1. 在你的仓库中创建 `.github/workflows/osp.yml` 文件：

```yaml
TODO
```

2. 配置权限

- 进入仓库的 Settings -> Actions -> General
- 在 "Workflow permissions" 部分，选择 "Read and write permissions"
- 保存更改

3. 使用方式

- 自动运行：Action 会按照 cron 设置的时间自动运行
- 手动运行：
  1. 进入仓库的 Actions 页面
  2. 选择 "OSP Automation" workflow
  3. 点击 "Run workflow"

更多使用说明请参考 [Github Action 使用文档](docs/guide/github-action.md)。

## 文档

- [使用指南](docs/guide/README.md) -  使用指南
- [设计文档](docs/design/README.md) - 架构和实现细节
- [CLI 参考文档](docs/cli/osp.md) - CLI 参考文档

## 贡献

欢迎贡献代码和提出建议！请参考我们的[贡献指南](CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。
