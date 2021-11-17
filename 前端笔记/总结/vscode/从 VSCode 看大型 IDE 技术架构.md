https://zhuanlan.zhihu.com/p/96041706
两种代码编辑器—— Vim/VSCode 类和 IntelliJ IDEA/Visual Studio 类。
代码编辑器的两种解耦模式，分别是『编辑器只负责编辑，后台进程才看得懂代码语义』的可以不读完文件的 LSP 式，和『编辑器搞定一切』的必须读完文件的 JB 式。
根据上面给的 LSP 的 overview，可以看出 LSP 是通过给 language server 发送编辑造成的`增量改动`（range 或者 offset, text，分别对应插入和删除嘛）配合 language server 的`增量 parsing` 来实现更新的。 至于高亮（tokenize 同），要么从 language server 获取，要么用正则表达式或者 lexer。 可以看出，编辑器只知道插入和删除的 range ，代码对它来说只是一个一个的 token。
而 JB 式的编辑器会使用自己内置的 `parser 对 AST 进行更新`（可以做成增量的），代码对它来说是 AST ，可以拿到元素的父子节点，知道代码的层级嵌套关系以及详细的语义。 如果能拿到完整的 code base （大部分情况下都能拿到，但也有特例，比如根据 vczh 的说法 office 的开发就拿不到，因为代码太多了。这种代码也不适合使用对代码进行完整分析的 IDE 开发，一般都是编辑器家族或者大型 IDE 关掉一些分析功能来写），就可以对代码进行完整的静态分析，比如 IDEA 可以根据函数里对参数的处理，推导库函数参数的 @NotNull 和 @Nullable （只要直接调用了上面的函数，或者传给了需要 @NotNull 的函数，那么就是 @NotNull；如果有什么 if (参数 == null) return 默认值的语句，就会推断为 @Nullable，从而对错误的操作（比如传一个 null 给一个直接在这个参数上调用方法的函数，很明显的 NullPointerException）进行报警），寻找命令行报错信息中的文件名和行号进行并提供跳转到出错地点等。

LSP 眼中的代码，是一个`线性的 token 序列`。JB 眼中的代码，是一个`树形的语法结构`。 一个把全部代码分析交给插件，一个只把语法定义交给插件。

1. 谈起 Web IDE，没人能绕开 VSCode，它非常流行，同时又完全开源，总共 350000 行 TypeScript 代码的巨大工程，使用了 142 个开源库。
   大型复杂 GUI 软件（如 IDE 类）如何组织功能模块代码
   如何使用 Electron 技术将 Web 软件桌面化
   如何在打造插件化开放生态的同时保证软件整体质量与性能
   如何打造一款好用的、流行的工具软件
2. 核心功能：
   IntelliSense（代码提示）、 debugging（代码调试）、 git（代码管理）
   extensions （插件）则肩负着打造开放生态的责任
   对于工具软件而言，需要内心能想清楚边界。
   哪些是自己应该专注去做的，哪些可以外溢到交给第三方扩展来满足。
3. 生产力工具类的软件一定要守住主线，否则很可能会变成不收门票的游乐园
   牛逼的产品背后一定有牛逼的团队，比如微软挖到 Anders Hejlsberg，接连创造了 C# 和 TypeScript
   挖到 Erich Gamma，接连诞生了 monaco 和 vscode 这些明珠
4. Electron 是什么
   VSCode 有一个特性是跨平台，它的跨平台实质是通过 electron 实现的
   使用 Web 技术来编写 UI，用 chrome 浏览器内核来运行
   使用 NodeJS 来操作文件系统和发起网络请求
   使用 NodeJS C++ Addon 去调用操作系统的 native API

   **1 个主进程**：一个 Electron App 只会启动一个主进程，它会运行 package.json 的 main 字段指定的脚本
   **N 个渲染进程**：主进程代码可以调用 Chromium API 创建任意多个 web 页面，而 Chromium 本身是多进程架构，每个 web 页面都运行在属于它自己的渲染进程中

   `进程间通讯`：
   Render 进程之间的通讯本质上和多个 Web 页面之间通讯没有差别，可以使用各种浏览器能力如 localStorage
   **Render** 进程与 **Main** 进程之间也可以通过 API 互相通讯 (ipcRenderer/ipcMain)
   electron 的 web 页面所处的 Render 进程可以将任务转发至运行在 NodeJS 环境的 Main 进程，从而实现 native API
   这套架构大大扩展了 electron app 相比 web app 的能力丰富度
   但同时又保留了 web 快捷流畅的开发体验，再加上 web 本身的跨平台优势
   结合起来让 electron 成为性价比非常高的方案

5. 多进程架构
   主进程：VSCode 的入口进程，负责一些类似窗口管理、进程间通信、自动更新等全局任务
   渲染进程：负责一个 Web 页面的渲染
   插件宿主进程：每个插件的代码都会运行在一个独属于自己的 NodeJS 环境的宿主进程中，插件不允许访问 UI
   Debug 进程：Debugger 相比普通插件做了特殊化
   Search 进程：搜索是一类计算密集型的任务，单开进程保证软件整体体验与性能
6. 隔离内核 (src) 与插件 (extensions)，内核分层模块化
   /src/vs：分层和模块化的 core
   /src/vs/base: 通用的公共方法和公共视图组件
   /src/vs/code: VSCode 应用主入口
   /src/vs/platform：可被依赖注入的各种纯服务
   /src/vs/editor: 文本编辑器
   /src/vs/workbench：整体视图框架
   /src/typings: 公共基础类型
   /extensions：内置插件
7. 每层按环境隔离
   内核里面每一层代码都会遵守 `electron 规范`，按不同环境细分文件夹:
   common: 公共的 js 方法，在哪里都可以运行的
   browser: 只使用浏览器 API 的代码，可以调用 common
   node: 只使用 NodeJS API 的代码，可以调用 common
   electron-browser: 使用 electron 渲染线程和浏览器 API 的代码，可以调用 common，browser，node
   electron-main: 使用 electron 主线程和 NodeJS API 的代码，可以调用 common， node
   test: 测试代码

   因此云凤蝶最终也决定采用 (editor/runtime/common) 类似的隔离架构

8. ` 内核代码本身也采用扩展机制`: Contrib
   可以看到 /src/vs/workbench/contrib 这个目录下存放着非常多的 VSCode 的小的功能单元
   ├── backup
   ├── callHierarchy
   ├── cli
   ├── codeActions
   ├── codeEditor
   ├── comments
   ├── configExporter
   ├── customEditor
   ├── debug
   ├── emmet
   ├──....中间省略无数....
   ├── watermark
   ├── webview
   └── welcome
   Contrib 目录下的所有代码不允许依赖任何本文件夹之外的文件
   Contrib 主要是使用 Core 暴露的一些扩展点来做事情
9. VSCode 的代码大量使用了依赖注入
   VSCode 的依赖注入，它没有使用 reflect-metadata 这一套，而是基于 decorator 去标注元信息，整个实现了一套自己的依赖注入方式
10. `绝对路径 import`+tsconfig 里的`paths属性`
    绝对路径 import 是一个非常值得学习的技巧，具体的方式是配置 TypeScript compilerOptions.paths
    相对路径 import 对阅读者的大脑负担高，依赖当前文件位置上下文信息才能理解
    假设修改代码的时候移动文件位置，相对路径需要修改本文件的所有 import，绝对路径不需要
11. 命令系统
    VSCode 和 monaco-editor 都有自己的命令系统
    一个功能的触发方式是多种多样的，比如大纲树右键菜单触发，顶部工具栏触发，画布右键菜单触发，键盘快捷键触发等
    这就要求该功能的实现函数能跨区域在多个位置调用，这很容易导致代码依赖关系异常混乱。
    而命令系统就是一种解决这个问题的很好思路
12. 代码编辑器技术
    monaco-editor:
    Text Buffer 性能优化
    MVVM 架构
    language server protocol:
    `不再关注 AST 和 Parser，转而关注 Document 和 Position，从而实现语言无关。`
    `将语言提示变成 CS 架构，核心抽象成当点击了文档的第几行第几列位置需要 server 作出什么响应的一个简单模型，基于 JSON RPC 协议传输，每个语言都可以基于协议实现通用后端`
    monaco-editor 可以看出专家级人物的领域积累，属于 VSCode 的核心竞争力
    language server protocol 和 Debug Adaptor Prototal 的设计就属于高屋建瓴
    可以看出思想层次，把自己的东西做成标准，做出生态，开放共赢
13. VSCode 插件系统
    对比几大 IDE：

    Visual Studio / IntelliJ：不需要插件，all in one （不够开放）
    Eclipse: 一切皆插件 （臃肿、慢、不稳定、体验差）
    VSCode：中庸之道

    **VSCode 插件的强隔离**
    **Workbench 视图结构**
    标题栏: Title Bar
    活动栏: Activity Bar
    侧边栏: Side Bar
    面板: Panal
    编辑器: Editor
    状态栏: Status Bar
    **插件 API 注入**
    vscode plugin 的 require 全部被 Microsoft/vscode-loader 劫持了，通过对 require 的 hack 将插件 API 注入到了运行环境。
    **插件开发与配置**
    一个插件核心就是一个配置文件：Extension Manifest JSON (package.json 里面的一个字段)
    Activation Events 激活时机
    Contribution Points 扩展点

14. 感想
    健壮的结构远比精巧的设计来得重要。换句话说，结构是第一位的，功能是第二位的。
    保持项目的简单性。设计达到完美的时候，不是无法再增加东西了，而是无法再减少东西了。
    VSCode 的一些工程上的优秀设计，比如依`赖注入、绝对路径引用、命令系统`对于云凤蝶来说是可以马上学以致用的，而 contrib 与 extension 的扩展系统，则非一日之功，也并不宜盲目下手。
    而事实上在尝试打造每一个开发者都梦想的万物皆 plugin 式的工具软件之前，有一些通用的问题需要先冷静下来思考：

    用户核心在操作的资源是什么？
    用户的关键路径是什么？
    这个软件的整体功能形态，交互与视觉设计已经稳定了吗？
    内核功能区和第三方扩展的功能域之间的界限在哪里？
    哪些环节可能会出现外溢需求需要第三方扩展才能被满足，不适宜官方动手做吗？
    对 VSCode 而言：

    核心操作的资源是`文件`
    关键路径是：`打开文件 - 编辑文件 - 保存文件`
    整体功能设计，交互与视觉设计非常稳定
    内核是文件管理与代码编辑，多样性的编程语言生态，CICD 等衍生研发链路等可能会出现扩展需求
