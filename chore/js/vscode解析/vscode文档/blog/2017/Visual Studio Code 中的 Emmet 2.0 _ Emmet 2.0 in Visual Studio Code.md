# Visual Studio Code 中的 Emmet 2.0

链接：https://code.visualstudio.com/blogs/2017/08/07/emmet

## 一针见血的分析

这篇文章记录了 VS Code 将核心功能“插件化”的一次大规模重构：

1.  **从单一代码库到模块化生态**：Emmet 2.0 放弃了过去臃肿的单体代码，转而采用 `@emmetio` 模块化生态。这意味着解析、缩写扩展、渲染被拆解为独立的 npm 模块。这种做法不仅方便了 VS Code 的集成，也让 Emmet 能够被第三方工具轻松复用。
2.  **“Tab 键之争”的终结**：长期以来，Emmet 的默认扩展键是 Tab，这经常与代码缩进产生冲突。VS Code 1.15 做出了一个大胆的决定：**取消 Tab 键默认绑定**，转而将 Emmet 建议集成进 IDE 的建议列表（Suggestion List）。用户现在可以像选择变量名一样选择 Emmet 展开。这不仅解决了冲突，还让 Emmet 的发现性（Discoverability）大大增强。
3.  **核心精简策略**：通过将 Emmet 从 VS Code 核心代码中剥离，转型为独立的 **Extension**（内置扩展），VS Code 的核心架构变得更加轻巧。这种“内功外放”的策略也成了后来 VS Code 治理大型功能的标准模版。
4.  **多光标的降维打击**：受益于新架构，Emmet 终于在大多数操作中支持了多光标（Multi-cursor）。这使得 HTML 结构的批量修改效率提升了一个数量级。

## 摘要

2017年8月，VS Code 正式发布了 Emmet 2.0 集成。这次更新通过模块化的 `@emmetio` 重写了所有动作，并将 Emmet 建议集成到自动补全列表中，取代了传统的 Tab 键扩展方式。此外，Emmet 正式成为了 VS Code 的一个独立内置扩展，支持多光标操作，显著提升了 Web 开发效率。
