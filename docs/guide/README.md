# OSP 使用指南

OSP (Open Source Planning) 是一个用于管理开源社区规划的命令行工具。它可以帮助你基于 GitHub 里程碑生成和更新规划文档。

## 安装

```bash
go install github.com/elliotxx/osp@latest
```

## 使用方法

### 认证

在使用 OSP 之前，你需要先登录 GitHub。OSP 使用 GitHub CLI 的认证机制，所以你需要先安装并登录 GitHub CLI：

```bash
# 安装 GitHub CLI
brew install gh

# 登录 GitHub
gh auth login
```

### 仓库管理

OSP 支持管理多个仓库，你可以使用以下命令来管理仓库：

```bash
# 添加一个仓库
osp repo add owner/repo

# 列出所有仓库
osp repo list

# 切换当前仓库
osp repo switch owner/repo

# 查看当前仓库
osp repo current

# 移除一个仓库
osp repo remove owner/repo
```

### 规划管理

一旦你选择了一个仓库，你就可以使用以下命令来生成和更新规划文档：

```bash
# 基于里程碑生成规划文档
osp plan <milestone-number>

# 使用自定义标签
osp plan <milestone-number> --label planning

# 使用自定义分类
osp plan <milestone-number> --categories bug,documentation,enhancement

# 包含 PR
osp plan <milestone-number> --exclude-pr=false
```

规划文档会包含以下内容：

1. 里程碑的基本信息
2. 进度条和统计信息
3. 按类别分组的 issues
4. 贡献者列表

## 配置

OSP 的配置文件存储在以下位置：

- macOS: `~/Library/Application Support/osp/config.yml`
- Linux: `~/.config/osp/config.yml`
- Windows: `%AppData%\osp\config.yml`

配置文件的格式如下：

```yaml
auth:
  token: ""  # GitHub 令牌，由 GitHub CLI 提供
current: ""  # 当前选中的仓库
repositories:  # 管理的仓库列表
  - owner/repo1
  - owner/repo2
```
