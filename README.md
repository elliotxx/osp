# OSP - Open Source Pilot

[![Go Report Card](https://goreportcard.com/badge/github.com/elliotxx/osp)](https://goreportcard.com/report/github.com/elliotxx/osp)
[![GoDoc](https://godoc.org/github.com/elliotxx/osp?status.svg)](https://godoc.org/github.com/elliotxx/osp)
[![License](https://img.shields.io/github/license/elliotxx/osp.svg)](https://github.com/elliotxx/osp/blob/main/LICENSE)

> Automated Open Source Software Management

OSP (Open Source Pilot) 是一个自动化的开源软件管理工具，它帮助开源项目维护者更高效地管理项目、跟踪进展、生成报告。

## 特性

- 🔑 GitHub 认证管理
- 📊 项目数据统计和分析
- 📝 自动生成项目规划
- ✨ 社区任务管理
- 📈 Star 历史追踪
- 📅 活动报告生成

## 快速开始

### 安装

```bash
go install github.com/elliotxx/osp@latest
```

### 基本使用

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
osp repo switch owner/repo

# 查看当前仓库
osp repo current
```

3. 生成规划
```bash
# 基于里程碑生成规划
osp plan <milestone-number>

# 使用自定义标签和分类
osp plan <milestone-number> --label planning --categories bug,documentation,enhancement
```

更多使用说明请参考 [使用文档](docs/usage/README.md)。

## 文档

- [使用文档](docs/usage/README.md) - 安装和使用指南
- [设计文档](docs/design/README.md) - 架构和实现细节
- [API 文档](docs/api/README.md) - API 参考

## 贡献

欢迎贡献代码和提出建议！请参考我们的[贡献指南](CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。
