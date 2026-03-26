这份 README.md 是 `@mariozechner/pi-web-ui` 这个前端库的官方说明文档。对于新手来说，直接看代码和概念可能会有点懵。

简单来说，**这个库的目的是让你能用最快的速度，在网页上搭建一个类似 ChatGPT 或 Claude 的人工智能聊天界面**。它不仅包含聊天框，还自带了文件上传、代码运行沙盒、历史记录保存等高级功能。

我们来逐段、对照着 README 进行透彻的拆解：

---

### 1. 标题与简介 (Header & Intro)

**原文档**：`# @mariozechner/pi-web-ui` -> 讲明了这是一个基于 `mini-lit`（一个轻量级 Web 组件库）和 `Tailwind CSS v4`（一个流行的 CSS 样式库）构建的 UI 库。
**新手解析**：如果你想做个带 AI 的网页，不用自己从零开始写对话框、气泡、输入框，直接用这个库就行。它底层依赖了作者写的另外两个纯逻辑库：`pi-ai`（负责调用各种 AI 模型）和 `pi-agent-core`（负责管理 AI 的状态和思考过程）。

### 2. 核心特性 (Features)

这里列出了这个库能干什么：

- **Chat UI (聊天界面)**: 现成的对话记录、打字机效果（流式输出）和工具执行状态。
- **Tools (工具)**: 允许 AI 在网页里直接运行 JS 代码 (REPL)，抓取网页正文，或者生成复杂的卡片（Artifacts）。
- **Attachments (附件)**: 支持用户上传 PDF、Word、Excel、图片等，自带预览和文字提取功能（把文件变成 AI 能读懂的文字）。
- **Artifacts (人造物/衍生品)**: 类似 Claude 的 Artifacts 功能，AI 生成的代码（HTML/SVG/Markdown）可以在网页右侧独立渲染并交互。
- **Storage (存储)**: 用浏览器自带的数据库 (IndexedDB) 保存你的聊天记录、API 密钥和设置，刷新页面内容不会丢。
- **CORS Proxy (跨域代理)**: 在浏览器直接调大模型 API 经常会被浏览器的安全机制（跨域拦截）挡住，它自带了代理配置来解决这事。
- **Custom Providers (自定义模型)**: 支持连上 Ollama（本地跑模型）、vLLM 或者其他兼容 OpenAI 格式的接口。

### 3. 安装 (Installation)

**新手解析**：在你的项目里运行这行代码，把 UI 库、核心代理库、AI 模型库这“三剑客”一次性装好。

### 4. 快速上手 (Quick Start)

文档给出了一长串代码。我们把它拆成三步看：

1.  **设置存储 (Set up storage)**：创建了设置、API Key、会话这三个数据仓库，并把它们绑定到 `IndexedDBStorageBackend`（这就像是给你的应用建了一个本地的 MySQL 数据库，存放在用户浏览器里）。
2.  **创建 AI 大脑 (Create agent)**：实例化了一个 `Agent`，给了它一个基础设定（“你是一个得力助手”），指定了它要借用哪家的大模型（比如 Claude 3.5 Sonnet）。
3.  **创建界面并挂载 (Create chat panel)**：通过 `new ChatPanel()` 召唤出聊天界面主体，把前面的 Agent 塞进去，并把它放到网页的 `<body>` 里显示出来。

### 5. 架构图 (Architecture)

原文档画了一个从上到下的关系图：

- **最上层 UI (ChatPanel)**：用户看到的聊天区 (AgentInterface) 和 代码渲染区 (ArtifactsPanel)。
- **中间大脑 (Agent)**：处理业务逻辑，记得你发了什么，AI 回了什么，何时该调用工具。
- **底层存储 (AppStorage)**：保存所有的配置和数据到浏览器的 IndexedDB。

### 6. 组件详情 (Components)

- **ChatPanel**：这是个“傻瓜式”全功能集成的面板。包括了对话和右侧的 Artifact（沙盒展示区）。它还提供了一些“钩子件（Hooks）”，比如在发送消息前做点什么（`onBeforeSend`），或者缺 API Key 时弹窗让你输入（`onApiKeyRequired`）。
- **AgentInterface**：这是比 ChatPanel 更底层一点的纯对话框。如果你只想把聊天框嵌入到你自己的奇特网页布局中，可以用它。它有开关能控制是否显示“附件按钮”、“模型选择器”等。
- **Agent (来自 pi-agent-core)**：这是核心驱动。通过 `.subscribe` 你可以监听到 AI 的一举一动（它开始想了、开始说话了、说完了）。你可以用 `agent.prompt()` 发消息，用 `agent.abort()` 强制中断它说话。

### 7. 消息类型 (Message Types)

这个库把聊天中的消息分了类：

- **UserMessageWithAttachments**：带着附件（比如图片、文档）的用户消息。
- **ArtifactMessage**：AI 生成的需要在右边独占一块区域显示的“复杂产物”（比如一个图表网页）。
- **Custom Message Types**：高级玩法，教你如何自己定义新消息类型（比如定义一种红色的“系统警告消息”）。

### 8. 消息转换工具 (Message Transformer)

**新手解析**：大模型（如 OpenAI、Claude）其实是不懂你的 UI 上的奇奇怪怪的数据结构（比如 `ArtifactMessage`）的。`convertToLlm` 这个函数相当于一个“翻译官”，在发给大模型前，把你在网页上的这些特殊消息洗成了大模型能听懂的纯净格式。

### 9. 工具 (Tools)

赋予 AI 更多的“手脚”：

- **JavaScript REPL**：让 AI 可以在浏览器的安全沙盒里写 JS 然后运行。
- **Extract Document**：给 AI 一个网址，它就能把那个网址的文章剥离出来阅读。
- **Artifacts Tool**：让 AI 创建网页、SVG 画图、Markdown 等内容。
- **Custom Tool Renderers**：如果你自己给 AI 写了个新能力（比如“查天气”），你可以在这定义这个能力返回的结果在页面上长什么样。

### 10. 本地存储管理 (Storage)

非常详细地讲述了怎么管理本地数据。

- **SettingsStore**：存诸如“是否开启代理”这样的开关。
- **ProviderKeysStore**：存你的 OpenAI Key、Anthropic Key，加密保存在浏览器。
- **SessionsStore**：管理聊天列表，可以储存、读取、修改标题、删除特定的历史对话。
- **CustomProvidersStore**：存你自己添加的本地大模型地址（如本地的 Ollama 地址：http://localhost:11434）。

### 11. 附件处理 (Attachments)

教你如何用代码调用它的 API 来读取文件。无论是从 `<input>` 标签上传的文件，还是给一个文件网址，它都能转成叫做 `Attachment` 的结构（里面包含了 base64 内存数据和提取出来的纯文本）。

### 12. 跨域代理 (CORS Proxy)

对于纯前端的 AI 应用，这是一个痛点。浏览器不让你在 A 网站直接发请求给 B 网站（大模型 API）。这里介绍了如何搭配像 `corsproxy.io` 这样的转发服务，绕过浏览器的紧箍咒。

### 13. 对话框/弹窗组件 (Dialogs)

预置了一堆造好的轮子：

- `SettingsDialog`：系统设置弹窗。
- `SessionListDialog`：抽屉式的历史聊天记录列表。
- `ApiKeyPromptDialog`：讨要 API 密钥的弹窗。
- `ModelSelector`：切换大模型的下拉弹窗。

### 14. 样式 (Styling) & 国际化 (Internationalization)

- **Styling**：教你引入写好的 CSS 把界面变漂亮，也支持通过 Tailwind 覆盖默认样式。
- **Internationalization (i18n)**：这玩意能支持多国语言。文档例子演示了怎么把系统的自带英文（比如 Loading...）翻译成德语（Laden...）。

### 15. 杂项 (Examples, Known Issues, License)

- 提供了一个示例路径和作者自己开发的实战产品 `sitegeist` 供参考。
- 目前已知 Bug：`PersistentStorageDialog` 这个组件目前有点问题。
- 授权是 MIT，意味着你可以随意免费使用和修改。

---

**💡 总结**：
读完这篇 README，你只需要知道：**要想做一个高级的网页 AI 助手，不用自己写样式和逻辑了。导进来它的 `ChatPanel`、接好 `AppStorage`（用作本地漫游记录）、设置好大模型的 Key，你的 AI 前端应用就跑起来了，连文件上传和沙盒渲染都送你了。** 对着 Quick Start 看一遍照着配置，是最好的入门方式。这就是这份文档传递的全部核心信息。满分！有什么想具体深入问的配置代码吗？
