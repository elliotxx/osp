# OSP 命令行设计

## 命令概述

OSP 是一个开源项目管理工具，提供以下主要命令：

### 认证管理 (auth)
```bash
osp auth login    # 登录 GitHub
osp auth logout   # 登出
osp auth status   # 查看认证状态
```

### 配置管理 (config)
```bash
osp config edit   # 编辑配置文件
osp config clean  # 清理配置文件
osp config list   # 列出配置信息
```

### 仓库管理 (repo)
```bash
osp repo add <owner/repo>     # 添加仓库
osp repo remove <owner/repo>  # 移除仓库
osp repo list                 # 列出所有管理的仓库
osp repo switch <owner/repo>  # 切换当前操作的仓库
osp repo current             # 显示当前操作的仓库
```

### 社区规划 (plan)
```bash
osp plan update [flags]       # 更新社区规划
  --label, -l string         # 规划标签 (默认 "planning")
  --categories strings       # 自定义分类
  --priorities strings      # 优先级标签列表
  --exclude-pr            # 排除 PR
  --dry-run             # 仅显示内容，不更新
  --yes                # 跳过确认
```

### Star 管理 (star)
```bash
osp star                # Star 当前仓库
osp star history       # 查看 Star 历史
```

### 数据统计 (stats)
```bash
osp stats              # 显示仓库统计信息
```

### 新手引导 (onboard)
```bash
osp onboard           # 管理社区贡献者引导内容
```

## 全局选项

```bash
--no-color     # 禁用彩色输出
--verbose, -v  # 详细输出模式
--version, -V  # 显示版本信息
--help, -h     # 显示帮助信息
```

## 设计原则

1. 命令结构
   - 使用子命令组织功能
   - 保持命令层级简单清晰
   - 相关功能组合在同一子命令下

2. 命名规范
   - 命令名使用简单的动词或名词
   - 参数名使用完整的描述性词语
   - 保持命名一致性和可读性

3. 用户体验
   - 提供合理的默认值
   - 支持简写参数
   - 详细的帮助文档
   - 清晰的错误提示

4. 输出设计
   - 支持彩色输出
   - 支持详细模式
   - 结构化的输出格式
   - 清晰的进度反馈

## 实现细节

### 命令框架

使用 [spf13/cobra](https://github.com/spf13/cobra) 构建命令行应用：

1. 命令注册
```go
rootCmd.AddCommand(
    newAuthCmd(),    // 认证管理
    newConfigCmd(),  // 配置管理
    newRepoCmd(),    // 仓库管理
    newPlanCmd(),    // 规划管理
    newStarCmd(),    // Star 管理
    newStatsCmd(),   // 统计信息
    newOnboardCmd(), // 新手引导
)
```

2. 命令执行
```go
if err := rootCmd.Execute(); err != nil {
    os.Exit(1)
}
```

### 错误处理

1. 错误类型
   - 参数错误
   - 网络错误
   - 权限错误
   - 业务逻辑错误

2. 错误输出
   - 使用统一的错误格式
   - 提供错误上下文
   - 给出修复建议

### 输出格式

1. 普通输出
   - 简洁的单行输出
   - 结构化的表格输出
   - 支持 JSON 格式

2. 详细输出
   - 包含调试信息
   - 显示操作步骤
   - 显示 API 调用

3. 进度显示
   - 长时间操作显示进度条
   - 批量操作显示计数器
   - 可中断的操作支持取消

## 未来规划

1. 功能扩展
   - 支持更多的统计维度
   - 添加更多的自动化功能
   - 增强社区管理能力

2. 交互优化
   - 添加交互式配置
   - 支持命令补全
   - 优化帮助文档

3. 集成增强
   - 支持更多的 CI/CD 平台
   - 增加插件系统
   - 提供 API 接口
