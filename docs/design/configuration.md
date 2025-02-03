# OSP 配置管理设计

## 目录结构

OSP遵循XDG Base Directory Specification：

- 配置：`$XDG_CONFIG_HOME/osp`（默认：`~/.config/osp`）
- 状态：`$XDG_STATE_HOME/osp`（默认：`~/.local/state/osp`）
- 数据：`$XDG_DATA_HOME/osp`（默认：`~/.local/share/osp`）
- 缓存：`$XDG_CACHE_HOME/osp`（默认：`~/.cache/osp`）

各平台的默认路径参考：

### Linux
- 配置：`~/.config/osp/`
- 状态：`~/.local/state/osp/`
- 数据：`~/.local/share/osp/`
- 缓存：`~/.cache/osp/`

### macOS
- 配置：`~/Library/Application Support/osp/`
- 状态：`~/Library/Application Support/osp/`
- 数据：`~/Library/Application Support/osp/`
- 缓存：`~/Library/Caches/osp/`

### Windows
- 配置：`%AppData%\osp\`
- 状态：`%AppData%\osp\`
- 数据：`%AppData%\osp\`
- 缓存：`%LocalAppData%\osp\cache\`

## 文件结构

### 配置文件

配置文件（`config.yaml`）存储应用程序配置：

```yaml
# 暂时为空，保留用于未来使用
```

### 状态文件

状态文件（`state.yaml`）存储运行时状态：

```yaml
# GitHub用户名用于认证
username: ""

# 当前选中的仓库
current: ""

# 管理的仓库列表
repositories:
  - owner/repo1
  - owner/repo2
```

## 实现细节

### 核心类型

```go
// Config代表应用程序配置
type Config struct {
    // 暂时为空，保留用于未来使用
}

// State代表应用程序状态
type State struct {
    // 用于认证的用户名
    Username string `yaml:"username,omitempty"`

    // 当前仓库
    Current string `yaml:"current,omitempty"`

    // 仓库列表
    Repositories []string `yaml:"repositories,omitempty"`
}
```

### 关键函数

1. 目录管理
```go
// 获取XDG目录路径
GetConfigHome() string
GetStateHome() string
GetDataHome() string
GetCacheHome() string

// 获取OSP特定目录
GetConfigDir() string  // 返回$XDG_CONFIG_HOME/osp
GetStateDir() string   // 返回$XDG_STATE_HOME/osp
```

2. 文件操作
```go
// 获取文件路径
GetConfigFile() string // 返回config.yaml路径
GetStateFile() string  // 返回state.yaml路径

// 加载和保存配置
Load(path string) (*Config, error)
Save() error

// 加载和保存状态
LoadState() (*State, error)
SaveState(state *State) error
```

3. 状态管理
```go
// 用户名操作
GetUsername() (string, error)
SaveUsername(username string) error
RemoveUsername() error

// 仓库操作
GetCurrentRepo() (string, error)
SaveCurrentRepo(current string) error
GetRepositories() ([]string, error)
SaveRepositories(repos []string) error
```

### 安全特性

1. 文件权限
   - 目录：0700（用户读写执行权限）
   - 文件：0600（用户读写权限）

2. 错误处理
   - 缺失文件的优雅处理
   - 自动目录创建
   - 清晰的错误消息

3. 调试日志
   - 文件操作的详细调试日志
   - 路径解析日志
   - 错误上下文日志

## 最佳实践

1. XDG兼容性
   - 遵循XDG Base Directory Specification
   - 使用标准目录结构
   - 支持平台特定路径

2. 状态管理
   - 将配置与运行时状态分离
   - 原子状态更新
   - 缺失状态的优雅处理

3. 安全性
   - 严格的文件权限
   - 用户只读访问
   - 不存储敏感数据

4. 错误处理
   - 描述性的错误消息
   - 优雅的回退
   - 调试日志支持
