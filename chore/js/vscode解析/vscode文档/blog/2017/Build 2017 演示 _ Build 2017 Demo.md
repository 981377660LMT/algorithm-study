# Build 2017 演示

链接：https://code.visualstudio.com/blogs/2017/05/10/build-2017-demo

## 深入分析

Build 2017是微软的年度开发者大会，VS Code团队为其精心准备了一场演示，体现了产品的最新进展。

### 演示的关键看点
1. **Live Share的首次亮相** - 允许多个开发者实时协作编辑同一个文件，类似Google Docs但针对代码
   - 技术挑战：实时同步、冲突解决、网络延迟
   - 竞争优势：没有其他主流编辑器提供这一功能
2. **调试体验的改进** - Debug Console的增强，支持更复杂的表达式求值
3. **多文件夹工作区** - 允许开发者在一个VS Code窗口中打开多个项目

### 战略意义
- Build大会是微软向全球开发者展示技术方向的舞台
- VS Code在这个舞台的亮相，标志着微软将其视为与Visual Studio并行的旗舰产品
- Live Share的演示，直接激发了远程协作开发的市场需求

### 技术创新点
- Live Share的实现依赖CRDT（Conflict-free Replicated Data Type）算法，确保协作数据的最终一致性
- 这套技术后来成为VS Code Liveshare扩展的核心，并最终被集成进vs.live.share的官方版本
