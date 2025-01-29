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
osp plan <milestone-number>                  # 基于里程碑生成规划文档
  --label planning                           # 自定义规划标签
  --categories bug,documentation,enhancement  # 自定义分类
  --exclude-pr=false                         # 包含 PR
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
osp activity       # 查看活动
  --type issue     # 按类型筛选
  --type pr
  --type discussion
  --period 7d      # 时间范围
  --format daily   # 日报格式
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
