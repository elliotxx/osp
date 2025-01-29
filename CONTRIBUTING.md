# 贡献指南

感谢你考虑为 OSP (Open Source Pilot) 做出贡献！

## 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat`: 新功能
- `fix`: 修复问题
- `docs`: 文档修改
- `style`: 代码格式修改
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat: 添加自动生成周报功能
fix: 修复认证过期问题
docs: 更新安装说明
```

## 代码风格

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 建议
- 添加必要的注释和文档
- 确保测试覆盖率

## 开发设置

1. 克隆项目
```bash
git clone https://github.com/yourusername/osp.git
```

2. 安装依赖
```bash
go mod download
```

3. 运行测试
```bash
go test ./...
```

4. 构建项目
```bash
go build ./cmd/osp
```

## 提交 PR 前的检查清单

- [ ] 通过所有测试
- [ ] 更新相关文档
- [ ] 添加必要的测试用例
- [ ] 遵循代码规范
- [ ] 提交信息符合规范

## 报告问题

报告问题时，请包含以下信息：

1. 问题描述
2. 复现步骤
3. 期望行为
4. 实际行为
5. 环境信息
   - OSP 版本
   - Go 版本
   - 操作系统
   - 其他相关信息

## 功能建议

我们欢迎新功能建议！请在提出建议时：

1. 检查现有 issues 避免重复
2. 详细描述新功能
3. 说明使用场景
4. 考虑实现方案

## 行为准则

请参阅我们的 [行为准则](CODE_OF_CONDUCT.md)。

## 许可证

通过贡献代码，你同意将代码以 [MIT 许可证](LICENSE) 授权。
