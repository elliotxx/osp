# OSP 配置管理设计

## 配置文件结构

### 主配置文件
位置：`~/.config/osp/config.yml`

```yaml
# 基础配置
version: "1.0"
debug: false
quiet: false

# 认证配置
auth:
  token: "xxx"           # GitHub token
  host: "github.com"     # GitHub host
  keyring: true         # 是否使用系统 keyring

# 仓库管理
repos:
  - name: "cli/cli"      # 仓库全名
    alias: "gh-cli"      # 别名
    path: "/local/path"  # 本地路径
    config:              # 仓库特定配置
      auto_sync: true
      sync_interval: "1h"
  
  - name: "org/repo"
    alias: "myrepo"
    config:
      auto_sync: false

# 当前选择的仓库
current_repo: "cli/cli"

# 默认设置
defaults:
  period: "7d"          # 默认时间周期
  format: "markdown"    # 默认输出格式
  auto_sync: true      # 默认是否自动同步
  cache_ttl: "24h"     # 缓存过期时间

# 自定义配置
custom:
  templates_dir: "~/.config/osp/templates"
  reports_dir: "~/.config/osp/reports"
```

## 数据存储结构

```
~/.config/osp/
  ├── config.yml          # 主配置文件
  ├── auth/              # 认证信息
  │   └── credentials    # 凭证文件
  ├── cache/            # 数据缓存
  │   ├── cli/cli/      # 按仓库组织缓存
  │   │   ├── stats/
  │   │   ├── stars/
  │   │   └── activities/
  │   └── org/repo/
  ├── templates/        # 模板文件
  │   ├── report/
  │   ├── plan/
  │   └── task/
  └── reports/         # 生成的报告
      ├── daily/
      ├── weekly/
      └── monthly/
```

## 配置管理机制

### 配置优先级
1. 命令行参数
2. 环境变量
3. 仓库特定配置
4. 用户配置文件
5. 默认配置

### 环境变量支持
- `OSP_TOKEN`: GitHub token
- `OSP_CONFIG`: 配置文件路径
- `OSP_DEBUG`: 调试模式
- `OSP_QUIET`: 静默模式
- `OSP_FORMAT`: 输出格式

### 配置验证
1. 启动时验证配置文件格式
2. 验证必要字段
3. 类型检查
4. 值范围检查

### 配置热重载
- 支持运行时重载配置
- 监听配置文件变化
- 平滑重载不影响运行中的命令

## 安全考虑

### 敏感信息处理
1. token 优先使用系统 keyring
2. 配置文件权限设置为 600
3. 敏感信息不写入日志

### 加密存储
- 使用系统 keyring 存储认证信息
- 支持自定义加密方案
- 加密密钥轮换机制

## 缓存策略

### 缓存配置
```yaml
cache:
  enabled: true
  ttl: "24h"
  max_size: "1GB"
  cleanup_interval: "1h"
```

### 缓存类型
1. API 响应缓存
2. 统计数据缓存
3. 报告缓存

### 缓存清理
- 定期清理过期缓存
- 超出大小限制时清理
- 支持手动清理

## 模板系统

### 模板位置
```
templates/
  ├── report/
  │   ├── daily.md
  │   ├── weekly.md
  │   └── monthly.md
  ├── plan/
  │   ├── roadmap.md
  │   └── milestone.md
  └── task/
      ├── good-first-issue.md
      └── help-wanted.md
```

### 模板变量
- 支持项目信息变量
- 支持统计数据变量
- 支持自定义变量

### 自定义模板
- 支持用户自定义模板
- 模板继承机制
- 模板验证
