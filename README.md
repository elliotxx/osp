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

1. 认证
```bash
osp auth login
```

2. 添加仓库
```bash
osp add owner/repo
```

3. 查看项目统计
```bash
osp stats
```

更多使用说明请参考 [使用文档](docs/usage/README.md)。

## 文档

- [设计文档](docs/design/README.md)
- [使用文档](docs/usage/README.md)
- [API 文档](docs/api/README.md)

## 贡献指南

欢迎贡献代码！请查看我们的 [贡献指南](CONTRIBUTING.md)。

## 许可证

[MIT License](LICENSE)
