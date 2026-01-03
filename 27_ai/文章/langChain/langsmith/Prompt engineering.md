# LangSmith Prompt Engineering 快速上手

Prompt Engineering 是设计、测试和优化 LLM 指令的过程。LangSmith 提供了创建、版本控制、测试和协作 Prompt 的工具。

### 1. 环境准备

- **LangSmith 账号**: [smith.langchain.com](https://smith.langchain.com)。
- **API Keys**: 需要 LangSmith API Key 和 OpenAI API Key。

### 2. UI 工作流

1.  **设置密钥**: 在 **Settings** -> **Secrets** 中添加 `OPENAI_API_KEY`。
2.  **创建 Prompt**: 在 **Prompts** 页面点击 **+ Prompt**。
3.  **测试 Prompt**:
    - 在 Playground 中点击齿轮图标配置模型（Provider, Model, Temperature 等）。
    - 在 **Inputs** 框中输入变量值，点击 **Start** 运行。
4.  **版本管理**: 每次保存都会创建一个新的 **Commit**。可以为特定的 Commit 添加 **Tag**（如 `prod`, `staging`）以便在代码中引用。

### 3. SDK 工作流 (TypeScript)

#### A. 安装依赖

```bash
npm install langsmith openai @langchain/core
export LANGSMITH_API_KEY='ls...'
export OPENAI_API_KEY='sk...'
```

#### B. 创建并推送 Prompt

```typescript
import { Client } from 'langsmith'
import { ChatPromptTemplate } from '@langchain/core/prompts'

const client = new Client()

const prompt = ChatPromptTemplate.fromMessages([
  ['system', '你是一个专业的助手。'],
  ['user', '{question}']
])

// 推送到 LangSmith Hub
await client.pushPrompt('my-first-prompt', {
  object: prompt
})
```

#### C. 拉取并测试 Prompt

```typescript
import { OpenAI } from 'openai'
import { pull } from 'langchain/hub'
import { convertPromptToOpenAI } from '@langchain/openai'

const oaiClient = new OpenAI()

// 拉取最新版本或指定版本 "my-first-prompt:tag_name"
const prompt = await pull('my-first-prompt')

// 填充变量
const formattedPrompt = await prompt.invoke({
  question: '天空是什么颜色的？'
})

const response = await oaiClient.chat.completions.create({
  model: 'gpt-4o',
  messages: convertPromptToOpenAI(formattedPrompt).messages
})
```

### 4. 核心概念

- **Prompts vs. Prompt Templates**:
  - **Prompt**: 实际发送给 LLM 的最终消息。
  - **Prompt Template**: 包含变量（如 `{question}`）的模板，用于动态生成 Prompt。
- **消息类型**:
  - **Chat Style**: 消息列表（System, User, AI），是目前主流模型推荐的格式。
  - **Completion Style**: 纯文本字符串，主要用于旧版模型。
- **模板格式**:
  - **f-string**: 默认格式，使用 `{variable}`。
  - **mustache**: 支持更复杂的逻辑，如条件判断 `{{#is_logged_in}}...{{/is_logged_in}}` 和循环。
- **工具 (Tools) 与 结构化输出 (Structured Output)**: 可以在 Prompt 中定义工具 Schema 或要求模型按特定 JSON 格式响应。

### 5. 版本管理 (Versioning)

- **Commits**: 每次保存都会生成唯一的 Commit Hash。
  - 支持查看历史记录、回滚版本。
  - 代码中可通过 `client.pull_prompt("name:hash")` 锁定特定版本。
- **Tags**: 人类可读的标签（如 `production`, `staging`, `v1`）。
  - 标签可以移动（指向不同的 Commit）。
  - **优势**: 代码引用标签 `pull("name:production")`，更新 Prompt 时只需在 UI 移动标签，无需修改代码。

### 6. Playground 进阶功能

- **多 Prompt 对比**: 在 Playground 中可以添加多个 Prompt 窗口，使用相同的输入对比不同指令的效果。
- **数据集测试**: 点击 Playground 右上角的 **Dataset**，可以针对整个测试集运行 Prompt，批量观察模型表现。
- **Polly (AI 辅助)**: 使用内置的 AI 助手 Polly 来优化 Prompt、生成工具定义或创建输出 Schema。
