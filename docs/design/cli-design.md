# OSP 命令行设计

## 命令概述

OSP 提供以下主要命令：

### 认证管理
```bash
osp auth login   # 登录 GitHub
osp auth logout  # 登出
osp auth status  # 查看认证状态
```

### 仓库管理
```bash
osp repo add <owner/repo>     # 添加仓库
osp repo remove <owner/repo>  # 移除仓库
osp repo list                 # 列出所有管理的仓库
osp repo switch <owner/repo>  # 切换当前操作的仓库
osp repo current             # 显示当前操作的仓库
```

### 规划管理
```bash
osp plan <milestone-number>                   # 基于里程碑生成规划文档
  --label, -l planning                        # 自定义规划标签
  --categories, -c bug,documentation          # 自定义分类
  --priorities, -p high,medium,low            # 自定义优先级标签列表，按优先级从高到低排序，前两个标签将显示在高优先级任务部分
  --exclude-pr, -e=false                      # 排除 PR（默认：false）
  --dry-run, -d                               # 仅显示内容，不更新
  --yes, -y                                   # 跳过确认，自动更新
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

### 任务管理
```bash
osp task           # 查看/生成任务
  --type good-first-issue  # 新手任务
  --type help-wanted      # 需要帮助
  --auto                 # 自动生成
```

### 活动追踪
```bash
osp activity         # 查看活动
  --type issue      # Issue 活动
  --type pr        # PR 活动
  --period 7d     # 指定统计周期
```

## 命令设计原则

1. 一致性
   - 所有命令都遵循相同的命名和参数格式
   - 使用统一的输出格式和错误处理

2. 简洁性
   - 命令名称简短明了
   - 参数名称易于理解和记忆

3. 友好性
   - 提供详细的帮助信息
   - 错误信息清晰有用

4. 可扩展性
   - 命令结构支持未来扩展
   - 参数设计预留扩展空间

## 命令实现

### 认证管理

认证命令使用 GitHub CLI 的认证机制，主要功能：

1. `auth login`
   - 使用 GitHub CLI 的令牌
   - 验证令牌有效性

2. `auth status`
   - 检查认证状态
   - 显示令牌信息

3. `auth logout`
   - 清除认证信息

### 仓库管理

仓库管理命令用于维护仓库列表，主要功能：

1. `repo add`
   - 验证仓库是否存在
   - 添加到配置文件

2. `repo list`
   - 显示所有仓库
   - 标记当前仓库

3. `repo switch`
   - 验证仓库是否在列表中
   - 更新当前仓库

### 规划管理

规划管理命令用于生成和更新规划文档，主要功能：

1. 里程碑管理
   - 获取里程碑信息
   - 统计进度

2. 文档生成
   - 按类别分组
   - 生成进度条
   - 统计贡献者

## 错误处理

1. 输入验证
   - 验证必要参数
   - 检查参数格式

2. 运行时错误
   - API 调用失败
   - 配置文件错误

3. 错误信息
   - 清晰的错误描述
   - 可能的解决方案

## 输出格式

1. 成功信息
   - 使用 ✓ 表示成功
   - 简洁的成功消息

2. 错误信息
   - 使用 × 表示错误
   - 详细的错误原因

3. 列表输出
   - 清晰的缩进
   - 当前项标记

## 全局参数

OSP 支持以下全局参数：

- `--verbose, -v`: 显示详细的执行日志
  - 默认值：false
  - 当启用时，会显示所有级别的日志信息，包括追踪信息
  - 当禁用时，只显示操作、成功和错误信息

- `--no-color`: 禁用颜色输出
  - 默认值：false
  - 当启用时，会禁用所有颜色输出

## 日志级别

OSP 使用不同的符号和颜色来标记不同级别的日志信息：

- `»` 表示追踪信息 (浅灰色)
  - 仅在 verbose 模式下显示
  - 用于展示详细的执行步骤和中间状态
  - 例如：`» Found milestone: v1.0.0 (#1)`

- `+` 表示正在执行的操作 (蓝色)
  - 默认显示
  - 用于提示用户当前正在进行的操作
  - 例如：`+ Updating existing planning issue #3`

- `✓` 表示操作成功完成 (绿色)
  - 默认显示
  - 用于确认操作已经成功完成
  - 例如：`✓ Successfully updated planning issue #3`

- `×` 表示操作失败 (红色)
  - 默认显示
  - 用于提示用户操作失败和错误信息
  - 例如：`× Failed to get milestone: 404 Not Found`

此外，OSP 还支持自定义日志格式：

- 缩进级别：通过 `L(level)` 设置，每级缩进两个空格
- 自定义前缀：通过 `P(prefix)` 设置，例如使用 `→` 表示处理步骤
- 自定义颜色：通过 `C(color)` 设置，支持以下颜色：
  - 红色 (ColorRed)
  - 绿色 (ColorGreen)
  - 黄色 (ColorYellow)
  - 蓝色 (ColorBlue)
  - 紫色 (ColorPurple)
  - 青色 (ColorCyan)
  - 灰色 (ColorGray)
  - 加粗样式 (StyleBold)

示例输出：
```
+ Found 2 items
  → Processing item 1
  ✓ Item 1 processed
  → Processing item 2
  × Failed to process item 2
```

## 未来扩展

计划添加的命令：
```bash
osp report        # 生成综合报告
osp contrib       # 贡献者分析
osp health        # 项目健康度
osp trend         # 发展趋势分析
osp compare       # 与其他项目对比
osp action        # GitHub Action 管理
