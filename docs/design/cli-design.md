# OSP 命令行设计

## 命令概述

OSP 提供以下主要命令：

### 认证管理
```bash
osp auth login     # 登录 GitHub
osp auth logout    # 登出
osp auth status    # 查看认证状态
```

### 仓库管理
```bash
osp add <repo>     # 添加仓库
osp remove <repo>  # 移除仓库
osp list           # 列出所有管理的仓库
osp switch <repo>  # 切换当前操作的仓库
osp current        # 显示当前操作的仓库
```

### 数据统计
```bash
osp stats          # 项目概况
  --period 7d      # 指定统计周期
  --format json    # 输出格式

osp star history   # Star 历史
  --from 2024-01   # 起始时间
  --to 2024-12     # 结束时间
  --trend          # 显示趋势分析
```

### 规划管理
```bash
osp plan           # 查看/生成规划
  --update         # 更新已有规划
  --auto          # 自动生成建议
  --format md     # 输出格式
```

### 任务管理
```bash
osp task           # 查看/生成任务
  --type good-first-issue  # 新手任务
  --type help-wanted      # 需要帮助
  --auto                 # 自动生成
```

### 活动追踪
```bash
osp activity       # 查看活动
  --type issue     # 按类型筛选
  --type pr
  --type discussion
  --period 7d      # 时间范围
  --format daily   # 日报格式
```

## 命令行参数规范

### 全局选项
- `--help`, `-h`: 显示帮助信息
- `--version`, `-v`: 显示版本信息
- `--config`: 指定配置文件路径
- `--debug`: 启用调试模式
- `--quiet`: 静默模式
- `--format`: 输出格式 (json, yaml, text)

### 时间范围参数
支持以下格式：
- 相对时间：1d, 7d, 1m, 1y
- 绝对时间：2024-01-01
- 特殊值：today, yesterday, last-week, last-month

### 输出格式
支持以下格式：
- text (默认)
- json
- yaml
- markdown
- csv

## 错误处理

所有命令遵循以下错误处理原则：
1. 使用有意义的错误代码
2. 提供清晰的错误信息
3. 在调试模式下显示详细的错误堆栈

## 交互设计

1. 进度显示
- 长时间操作显示进度条
- 可以通过 Ctrl+C 中断

2. 颜色支持
- 错误信息使用红色
- 警告信息使用黄色
- 成功信息使用绿色

3. 交互式操作
- 添加仓库时可以搜索
- 切换仓库时提供选择列表

## 未来扩展

计划添加的命令：
```bash
osp report        # 生成综合报告
osp contrib       # 贡献者分析
osp health        # 项目健康度
osp trend         # 发展趋势分析
osp compare       # 与其他项目对比
osp action        # GitHub Action 管理
```
