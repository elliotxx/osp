# OSP 认证设计

## 概述

OSP 使用 GitHub OAuth 设备流程进行认证，支持：
- 基于 GitHub OAuth 的设备流程认证，复用 GitHub CLI 的凭证，但获取的密钥独立存储，和 Github CLI 互不影响
- 安全的令牌存储，存储在系统级的密钥管理中
- 令牌验证和作用域检查

## 认证流程

### 1. 登录流程

1. 初始化 OAuth 设备流程
   - 使用 GitHub CLI 的 OAuth 应用凭证
   - 请求作用域：`repo` 和 `read:org`

2. 用户交互
   - 显示一次性验证码
   - 自动打开浏览器到 GitHub 认证页面
   - 等待用户完成认证

3. 令牌处理
   - 获取访问令牌
   - 获取用户信息
   - 安全存储令牌

### 2. 令牌获取顺序

按以下顺序查找 GitHub 令牌：

1. 环境变量
   - `GH_TOKEN`
   - `GITHUB_TOKEN`

2. 系统密钥环
   - 使用 `zalando/go-keyring` 包
   - 服务名：`osp:github.com`

## 核心功能

### 1. 认证管理

```go
// 执行 GitHub OAuth 设备流程登录
Login() (string, error)

// 登出并移除存储的凭证
Logout() error

// 获取存储的 GitHub 令牌
GetToken() (string, error)

// 获取当前认证状态
GetStatus() ([]*Status, error)
```

### 2. 令牌管理

```go
// 将令牌保存到密钥环
SaveToken(username, token string) error

// 从密钥环移除令牌
RemoveToken() error

// 从密钥环获取存储的令牌
getStoredToken() (string, error)
```

### 3. 令牌验证

```go
// 验证令牌有效性
validateToken(token string) error

// 获取令牌作用域
getTokenScopes(token string) ([]string, error)

// 获取用户信息
getUserInfo(token string) (string, error)
```

## 状态管理

### Status 结构

```go
type Status struct {
    Username     string   // GitHub 用户名
    Token        string   // 访问令牌
    TokenDisplay string   // 用于显示的令牌（部分隐藏）
    StorageType  string   // 存储类型（环境变量或密钥环）
    IsKeyring    bool     // 是否存储在密钥环中
    Scopes       []string // 令牌作用域
    Active       bool     // 令牌是否有效
}
```

## 安全特性

1. 令牌存储
   - 使用系统密钥环安全存储令牌
   - 支持从环境变量读取令牌
   - 令牌显示时部分隐藏

2. 令牌验证
   - 验证令牌有效性
   - 检查令牌作用域
   - 验证用户信息

3. 错误处理
   - 详细的错误信息
   - 优雅的错误恢复
   - 调试日志支持

## 最佳实践

1. 令牌管理
   - 优先使用环境变量
   - 使用系统密钥环存储长期令牌
   - 定期验证令牌有效性

2. 作用域管理
   - 使用最小必要权限
   - 验证令牌作用域
   - 提示缺失权限

3. 安全性
   - 不在日志中显示完整令牌
   - 使用 HTTPS 进行 API 调用
   - 及时清理失效令牌

## 使用示例

### 1. 登录
```go
token, err := auth.Login()
if err != nil {
    log.Error("登录失败：%v", err)
    return
}
log.Success("登录成功")
```

### 2. 检查认证状态
```go
statuses, err := auth.GetStatus()
if err != nil {
    log.Error("获取状态失败：%v", err)
    return
}
for _, status := range statuses {
    log.Info("用户：%s", status.Username)
    log.Info("令牌：%s", status.TokenDisplay)
    log.Info("作用域：%v", status.Scopes)
}
```

### 3. 登出
```go
if err := auth.Logout(); err != nil {
    log.Error("登出失败：%v", err)
    return
}
log.Success("登出成功")
```
