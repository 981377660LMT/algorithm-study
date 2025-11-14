# [低代码仓库列表](https://lowcode-engine.cn/site/docs/guide/appendix/repos)

1. 引擎主包
   包含引擎的 4 大模块，入料、编排、渲染和出码。
   仓库地址：https://github.com/alibaba/lowcode-engine 子包明细：

   designer
   editor-core
   editor-skeleton
   engine
   ignitor
   plugin-designer
   plugin-outline-pane
   react-renderer
   react-simulator-renderer
   renderer-core
   types
   utils
   material-parser
   code-generator

2. 引擎官方扩展包
   包含了常用的设置器（setter）、跟 setter 绑定的插件等
   仓库地址：https://github.com/alibaba/lowcode-engine-ext 子包明细：

   - 设置器
     setter
     array-setter
     bool-setter
     classname-setter
     color-setter
     events-setter
     expression-setter
     function-setter
     i18n-setter
     icon-setter
     json-setter
     mixed-setter
     number-setter
     object-setter
     out.txt
     radiogroup-setter
     select-setter
     slot-setter
     string-setter
     style-setter
     textarea-setter
     variable-setter
   - 插件 plugin
     plugin-event-bind-dialog 事件绑定浮层
     plugin-variable-bind-dialog 变量绑定浮层

3. 低代码插件
   包含了常用的插件等
   仓库地址：https://github.com/alibaba/lowcode-plugins 子包明细：

   base-monaco-editor 基础代码编辑器
   plugin-code-editor 源码编辑面板
   plugin-datasource-pane 数据源面板
   plugin-manual 产品使用手册面板
   plugin-schema 页面数据面板
   plugin-undo-redo 前进/后退功能
   plugin-zh-cn 中英文切换功能

4. 引擎 demo
   展示使用引擎编排和渲染等模块以及相应的依赖资源配置基础 demo
   仓库地址：https://github.com/alibaba/lowcode-demo

5. 工具链包
   包含生成引擎生态元素（setter、物料、插件）的脚手架，启动脚本，调试插件等
   仓库地址：https://github.com/alibaba/lowcode-tools

6. 低代码数据源引擎
   负责在渲染&出码两种运行时实现数据源管理，承担低代码搭建数据请求的能力；
   仓库地址：https://github.com/alibaba/lowcode-datasource

7. 基础物料 & 物料描述
   仓库地址：https://github.com/alibaba/lowcode-materials

8. 出码 demo
   仓库地址：https://github.com/alibaba/lowcode-code-generator-demo

---

好的，根据您提供的仓库列表，我为您梳理了一个推荐的阅读顺序，旨在帮助您由浅入深、循序渐进地理解阿里低代码引擎的整个生态。

### 推荐阅读顺序

1.  **lowcode-demo (仓库 4)**

    - **目的**：快速上手，建立感性认识。
    - **原因**：首先运行官方 Demo，可以最直观地体验低代码编辑器的核心功能和最终产物。这能帮助你带着问题和目标去阅读后续的源码。

2.  **lowcode-engine (仓库 1)**

    - **目的**：理解引擎核心架构与脉络。
    - **原因**：这是所有功能的核心。在体验过 Demo 后，深入此仓库可以理解引擎是如何将物料、编排、渲染、出码等核心模块组织在一起的。建议内部子包阅读顺序：`types` -> `engine` -> `editor-skeleton` -> `designer` -> `renderer-core` -> `material-parser` -> `code-generator`。

3.  **lowcode-materials (仓库 7)**

    - **目的**：理解物料规范。
    - **原因**：物料是低代码平台搭建的“砖块”。了解了核心引擎后，需要理解物料是如何被定义和描述的，这样才能知道引擎是如何消费它们的。

4.  **lowcode-engine-ext (仓库 2)**

    - **目的**：学习属性配置面板（Setter）的实现。
    - **原因**：Setter 是连接“物料”和“编辑器”的桥梁，负责提供组件属性的可视化配置界面。理解了核心引擎和物料后，学习 Setter 是掌握编辑器扩展能力的关键一步。

5.  **lowcode-plugins (仓库 3)**

    - **目的**：学习如何扩展编辑器功能。
    - **原因**：插件系统是引擎扩展性的核心体现。通过阅读官方插件，你可以学会如何为编辑器增加新的面板、工具栏按钮等功能，例如源码编辑、数据源管理等。

6.  **lowcode-datasource (仓库 6)**

    - **目的**：掌握数据流管理。
    - **原因**：现代应用离不开数据。该仓库独立负责数据源的定义、请求和状态管理，是实现动态化页面的关键。通常它会与 lowcode-plugins 中的数据源面板插件配合使用。

7.  **lowcode-code-generator-demo (仓库 8)**

    - **目的**：深入理解出码能力。
    - **原因**：在了解了核心引擎中的 `code-generator` 模块后，通过这个 Demo 可以看到一个完整、可运行的出码示例，帮助你理解从页面 Schema 到项目源码的完整流程。

8.  **lowcode-tools (仓库 5)**
    - **目的**：掌握工程化能力。
    - **原因**：当您需要开发自己的物料、插件或 Setter 时，这个工具链仓库提供了脚手架、调试工具等，能极大提升开发效率。因此，在理解了所有核心概念后，最后来学习它是最合适的。

### 总结

总的来说，这个顺序遵循了 **“体验 -> 核心 -> 生态 -> 工具”** 的路径，希望能帮助您更高效地掌握阿里低代码引擎。
