# 使用 VS Code 调试 Java 应用

链接：https://code.visualstudio.com/blogs/2017/09/28/java-debug

## 一针见血的分析

这篇文章标志着 VS Code 从“前端编辑器”向“全能 IDE”进化的关键一环：

1.  **强强联手的生态策略**：VS Code 并没有尝试自己发明一个 Java IDE，而是与 Java 生态的巨头 **Red Hat** 合作。通过复用 Red Hat 的 Java 语言支持和微软的 Java 调试器（基于 Java Debug Server），实现了“站在巨人肩膀上”的跨越式发展。
2.  **LSP 之外的调试协议（DAP）**：Java 调试功能的落地是 **Debug Adapter Protocol (DAP)** 成熟的标志。DAP 使得像 Java 这种重文本语言也能在原本为轻量级编辑器设计的 VS Code 中拥有像 Eclipse/IntelliJ 一样的断点调试、变量检查、多线程堆栈监控等功能。
3.  **扩展包（Extension Pack）的引入**：这是 VS Code 首次大规模推广“扩展包”概念。通过将语言支持和调试器打包，解决了用户“配置成本过高”的痛点。这种“开箱即用”的体验是 VS Code 后来在后端开发领域迅速替代传统 IDE 的利器。

## 摘要

2017年9月，微软宣布与 Red Hat 合作，正式在 VS Code 中推出 Java 调试支持。基于 Java Debug Server 和 Debug Adapter Protocol (DAP)，VS Code 能够提供包括条件断点、数据检查和多线程调试在内的完整 Java 调试体验，并通过 Java Extension Pack 极大地简化了用户的上手流程。
