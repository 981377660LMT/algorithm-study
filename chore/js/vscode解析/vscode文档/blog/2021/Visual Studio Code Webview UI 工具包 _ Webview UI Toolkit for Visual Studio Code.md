# Visual Studio Code Webview UI 工具包

链接：https://code.visualstudio.com/blogs/2021/10/11/ui-toolkit

## 深入分析

VS Code的扩展开发者经常需要创建自定义UI（使用Webview），但之前缺乏标准的UI组件库。

### 问题
- 扩展开发者要么使用Web框架（React、Vue）自己写UI
- 要么复制粘贴代码片段，导致扩展间的UI风格不一致
- 没有官方的、与VS Code界面风格匹配的UI组件

### 解决方案
- Webview UI Toolkit提供了一套官方的UI组件（按钮、输入框、下拉菜单等）
- 这些组件与VS Code的设计系统一致，使用VS Code的颜色、字体等
- 组件是Web标准（Web Components），可与任何框架配合使用

### 技术特点
- 使用Web Components标准，而非某个特定框架
- 基于FAST设计系统（由微软开源）
- 轻量级，加载快，性能好

### 生态影响
- 扩展开发变得更容易，门槛降低
- 用户体验更一致，整体生态的质量提升
- 许多优质扩展开始采用这套工具包
