# OSP 配置管理设计

## 配置文件

### 位置

OSP 的配置文件存储在以下位置：

- macOS: `~/Library/Application Support/osp/config.yml`
- Linux: `~/.config/osp/config.yml`
- Windows: `%AppData%\osp\config.yml`

### 结构

配置文件使用 YAML 格式，主要包含以下部分：

```yaml
# 认证配置
auth:
  token: ""  # GitHub 令牌，由 GitHub CLI 提供

# 当前仓库
current: ""  # 当前选中的仓库

# 仓库列表
repositories:  # 管理的仓库列表
  - owner/repo1
  - owner/repo2
```

## 配置管理

### 配置加载

1. 配置初始化
   - 检查配置目录是否存在
   - 创建默认配置文件

2. 配置读取
   - 读取 YAML 文件
   - 解析配置项

3. 配置验证
   - 验证必要字段
   - 检查字段格式

### 配置更新

1. 仓库管理
   - 添加/移除仓库
   - 更新当前仓库

2. 配置保存
   - 序列化为 YAML
   - 写入文件

## 实现细节

### 配置结构体

```go
// Config represents the application configuration
type Config struct {
    Auth struct {
        Token string `yaml:"token"`
    } `yaml:"auth"`
    Current      string   `yaml:"current"`
    Repositories []string `yaml:"repositories"`
}
```

### 配置操作

1. 加载配置
```go
func Load(path string) (*Config, error)
```

2. 保存配置
```go
func (c *Config) Save() error
```

### 错误处理

1. 文件操作错误
   - 文件不存在
   - 权限不足
   - IO 错误

2. 解析错误
   - YAML 格式错误
   - 字段类型错误

3. 验证错误
   - 必要字段缺失
   - 字段格式错误

## 最佳实践

1. 配置验证
   - 在加载时验证配置
   - 在修改时验证更改

2. 错误恢复
   - 保持配置文件备份
   - 支持重置为默认值

3. 安全性
   - 敏感信息使用系统 keyring
   - 配置文件权限控制

4. 兼容性
   - 支持配置版本升级
   - 保持向后兼容
