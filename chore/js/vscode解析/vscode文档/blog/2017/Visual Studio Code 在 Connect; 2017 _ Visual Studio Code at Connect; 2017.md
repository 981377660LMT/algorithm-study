# Visual Studio Code 在 Connect; 2017

链接：https://code.visualstudio.com/blogs/2017/11/06/connect-2017

## 深入分析

Connect是微软的虚拟开发者大会（对标Build），VS Code团队用这个舞台宣布了什么？

### 主要公告
1. **Remote Development的预告** - VS Code即将支持远程开发，改变整个IDE的架构
2. **更多语言支持** - Python、Go、Rust等的官方扩展发布或改进
3. **性能优化** - 着手解决大型项目（>10000文件）的性能问题

### Remote Development的意义
- 这是VS Code历史上最雄心勃勃的功能之一
- 允许开发者在本地使用VS Code UI，但在远程机器（Linux服务器、容器等）中运行后端
- 解决的问题：跨平台开发、环境一致性、团队工作流
- 实现难度：需要彻底重构VS Code的架构，分离前端UI和后端运行时

### 技术远见
- 2017年的云计算和容器化浪潮，VS Code团队预见到了Remote Development的未来
- 这一决策，让VS Code在2019-2020年的远程开发热潮中处于绝对领先地位
