# 为 VS Code 推介 Chrome 调试器

链接：https://code.visualstudio.com/blogs/2016/02/23/introducing-chrome-debugger-for-vs-code

## 一针见血的分析

这篇文章标志着 VS Code 正式进入“全栈开发”工具阵营：

1.  **打通前后端调试的最后三公里**：在 2016 年之前，前端开发者习惯于在各种浏览器 DevTools 中调试。VS Code 通过 **Chrome Debugger Protocol** 实现了编辑器与浏览器的远程连接，让开发者可以在代码原地设置断点。这消除了“切屏调试”造成的认知负荷。
2.  **Source Maps：虚幻与现实的桥梁**：这篇文章强调了对 **Source Maps** 的深度支持。这在当时 WebPack / Babel 混淆后的前端代码调试中至关重要。通过映射，开发者调试的是“人类可读的代码”，而底层运行的是“机器优化的字节”，这种抽象能力的标准化是 DAP（调试适配协议）早期最成功的实践。
3.  **单点登录的愿景**：虽然当时的 Chrome 调试还有“只能连接一个客户端”的局限（打开浏览器 DevTools 就会断开 VS Code 连接），但它开启了“所有厂商工具协同工作”的愿景。这种开放心态使得 VS Code 迅速赢得了原本属于 Sublime Text 和 WebStorm 的前端市场。
4.  **DAP 协议的实战验证**：这不仅是个扩展，更是对 **Debug Adapter Protocol** 普适性的最好证明——同一个 UI 既能调 Node.js，也能调 Chrome。

## 摘要

2016年2月，微软发布了 Chrome Debugger 扩展，允许前端开发者在 VS Code 中直接调试运行在 Chrome 里的 JavaScript 脚本。通过对 Chrome 调试协议的封装和对 Source Maps 的支持，VS Code 实现了断点、单步执行、变量观察等完整调试体验，极大地简化了 Web 开发流程。
