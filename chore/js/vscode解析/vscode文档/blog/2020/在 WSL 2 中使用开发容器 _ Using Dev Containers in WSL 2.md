# 在 WSL 2 中使用开发容器

链接：https://code.visualstudio.com/blogs/2020/08/17/wsl2-dev-containers

## 深入分析

这是一个深层的技术整合，代表了"开发环境容器化"的趋势。

### 背景

- WSL 2（Windows Subsystem for Linux 2）于2019年发布，提供了真正的Linux内核
- DevContainers是VS Code的一个特性，允许在Docker容器中进行开发
- 将两者结合，可以在Windows上获得完整的Linux开发体验

### 架构

```
Windows (UI)
    ↓
VS Code Remote Container Extension
    ↓
WSL 2 (Linux Kernel)
    ↓
Docker (Container Runtime)
    ↓
开发环境（Node、Python等）
```

### 优势

1. **环境一致性** - 开发环境、测试环境、生产环境完全相同
2. **隔离性** - 多个项目可以使用不同的容器，不会冲突
3. **可复现性** - Dockerfile描述的容器可以被任何人复现

### 对开发者的影响

- 之前：开发者需要在本地安装各种依赖，容易遇到版本冲突、平台差异等问题
- 之后：一个Dockerfile解决所有问题，开发环境就像代码一样可版本控制

### 行业意义

- 这一实践的推广，推动了DevOps文化的深入应用
- 容器化开发变成了新的标准
