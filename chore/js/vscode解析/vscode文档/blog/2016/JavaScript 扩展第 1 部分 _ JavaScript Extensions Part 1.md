# JavaScript 扩展第 1 部分

链接：https://code.visualstudio.com/blogs/2016/09/14/js_roundup_1

## 一针见血的分析

这篇文章反映了 VS Code 早期对 JavaScript 生态的“精准爆破”策略：

1.  **“Salsa” 项目的余威**：2016年是 VS Code 将 JavaScript 核心解析引擎替换为 TypeScript 服务（即 Salsa 项目）后的第一年。通过“JS 只是没有类型的 TS”这一哲学，VS Code 瞬间让 JS 拥有了超越时代的 IntelliSense（跳转、补全、重构）。
2.  **补齐生态短板**：作为当时的“新秀”，VS Code 面临着来自 IDE（如 WebStorm）和编辑器（如 Sublime）的双重挤压。这篇文章推介的 **ESLint**、**Path IntelliSense**、**npm IntelliSense** 插件，标志着 VS Code 成功将原本繁琐的配置转化为“开箱即用”的体验。
3.  **桥接 IDE 与轻量编辑器**：文章明确提出了 VS Code 的定位——不是全能 IDE，但通过扩展性，它能够填补轻量编辑器在语言支持上的鸿沟。这种“渐进式功能”策略是其后来成功的关键。
4.  **Linters 的核心地位**：特意强调 ESLint 反映了当时前端工程化的趋势。VS Code 团队成员亲自维护 ESLint 扩展，体现了官方对前端代码质量工具的高度重视。

## 摘要

2016年9月，VS Code 团队推出了 JavaScript 扩展系列第一部分，重点推介了 ESLint、ES6 代码片段、路径智能提示以及 npm 智能提示等核心扩展。这些扩展结合了底层 Salsa 引擎带来的强大 JS 支持，确立了 VS Code 在构建现代化 JavaScript 开发环境方面的领先地位。
