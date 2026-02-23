# 推介 Insiders Build

链接：https://code.visualstudio.com/blogs/2016/11/21/introducing-insiders

## 深入分析

Insiders Build是VS Code持续创新的关键机制，值得深入讨论。

### 为什么需要Insiders Build？
1. **快速反馈循环** - 不能等待稳定版本（当时是每月发布），需要每日构建来验证新想法
2. **社区参与** - 早期用户（Power Users）可以更早体验新功能并反馈
3. **降低风险** - 在大规模用户前做beta测试，问题暴露更快

### Insiders Build的优势
- **独立安装** - 可与稳定版本并行运行，开发者可以"两手准备"
- **快速迭代** - 从提交代码到Insiders release平均仅需2-3小时
- **数据驱动** - 遥测数据（telemetry）帮助团队了解用户的使用模式

### 竞争优势
- JetBrains的IntelliJ虽然功能强大但更新周期长（每年2-3个主版本）
- VS Code通过Insiders，每周都有令人惊喜的新功能
- 这种高频创新，让VS Code在市场中总是"最新鲜"的

### 缺陷与权衡
- Insiders的稳定性难以保证，偶尔出现严重bug导致工作流中断
- 不是所有开发者都有容错能力，这限制了Insiders的用户基数
- 但正是这个"牺牲一部分用户"的权衡，让稳定版本更加可靠
