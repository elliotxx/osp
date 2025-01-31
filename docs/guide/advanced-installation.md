# 高级安装指南

本指南提供了多种安装 OSP 的方法，包括从源码编译、二进制下载和包管理器安装等。

## 目录
- [从源码安装](#从源码安装)
- [二进制安装](#二进制安装)
- [包管理器安装](#包管理器安装)
  - [Homebrew](#homebrew)
  - [Go](#go)

## 从源码安装

### 前置条件
- Go 1.22+
- Git
- Make

### 步骤

1. 克隆仓库
```bash
git clone https://github.com/elliotxx/osp.git
cd osp
```

2. 编译安装

```bash
# 默认编译（amd64 架构）
make build

# 指定架构编译（例如 arm64）
GOARCH=arm64 make build

# 启用 CGO
CGO_ENABLED=1 make build

# 安装到系统
sudo mv _build/{GOOS}/osp /usr/local/bin/
```

> 注意：如果在 macOS 上遇到 "panic: permission denied" 错误，请访问 [macos-golink-wrapper](https://github.com/eisenxp/macos-golink-wrapper) 查看解决方案。

## 二进制安装

我们为主流操作系统和架构提供预编译的二进制文件。

### Linux/macOS

1. 从 [Release 页面](https://github.com/elliotxx/osp/releases) 下载对应平台的压缩包

2. 解压并安装
```bash
# 示例：Linux AMD64
tar -zxvf osp_linux_amd64.tar.gz
sudo mv osp /usr/local/bin/

# 示例：macOS ARM64
tar -zxvf osp_darwin_arm64.tar.gz
sudo mv osp /usr/local/bin/
```

### Windows

1. 从 [Release 页面](https://github.com/elliotxx/osp/releases) 下载 Windows 版本
2. 解压 zip 文件
3. 将解压后的 `osp.exe` 添加到系统环境变量 PATH 中

## 包管理器安装

### Homebrew

macOS 和 Linux 用户可以使用 Homebrew 安装：

```bash
# 添加仓库
brew tap elliotxx/tap

# 安装最新版
brew install osp

# 安装指定版本
brew install osp@0.1.1

# 升级到最新版
brew upgrade osp

# 切换版本
brew unlink osp
brew link osp@0.1.1
```

### Go

使用 Go 1.22+ 可以直接安装：

```bash
# 安装最新版
go install github.com/elliotxx/osp@latest

# 安装指定版本
go install github.com/elliotxx/osp@v0.1.1
```

## 验证安装

安装完成后，运行以下命令验证：

```bash
osp --version
```

## 卸载

### Homebrew 卸载
```bash
brew uninstall osp
brew untap elliotxx/tap  # 可选
```

### 手动卸载
```bash
sudo rm /usr/local/bin/osp
```
