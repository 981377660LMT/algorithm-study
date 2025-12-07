LangChain 和 LangGraph 是目前构建 LLM（大语言模型）应用最流行的两个框架。简单来说，**LangChain 是基础构建块，而 LangGraph 是用于构建复杂、有状态、多智能体（Multi-Agent）系统的编排引擎。**

以下是对两者的深入剖析：

### 1. LangChain：LLM 应用的“胶水”与“脚手架”

LangChain 的核心理念是将 LLM 与外部数据源、计算能力连接起来。它解决的是“如何让模型不仅仅是聊天，而是去执行任务”的问题。

#### 核心组件剖析：

- **Chains (链):** 最基本的单元。它将多个操作（如 Prompt -> LLM -> Output Parser）串联起来。
  - _局限性:_ 传统的 Chain 是有向无环图（DAG），流程是线性的，很难处理循环逻辑（例如：如果不满意结果，重试直到满意）。
- **Prompts (提示词模板):** 管理和优化发送给模型的指令。
- **Retrieval (RAG):** 包含 Document Loaders（加载文档）、Text Splitters（切分文本）、Vector Stores（向量数据库）和 Retrievers（检索器）。这是构建知识库问答的核心。
- **Agents (智能体 - 旧版):** LangChain 早期的 Agent 使用 `AgentExecutor`。它通过让 LLM 决定下一步调用哪个 Tool（工具）。
  - _痛点:_ `AgentExecutor` 是一个黑盒，很难定制内部的循环逻辑，调试困难，且状态管理较弱。

#### 适用场景：

- 简单的问答系统 (RAG)。
- 一次性的任务执行（输入 -> 处理 -> 输出）。
- 数据提取和结构化。

---

### 2. LangGraph：构建有状态、循环的智能体系统

State、Node、Edge、Graph

LangGraph 是 LangChain 团队推出的一个扩展库，旨在解决 LangChain 在构建复杂 Agent 时的痛点。**它的核心是将 Agent 的流程建模为一个图（Graph）。**

#### 核心概念剖析：

1.  **State (状态):**

    - LangGraph 的核心是 `State` 对象。
    - 图中的每个节点（Node）都可以读取状态，并写入新的信息来更新状态。
    - 这解决了“记忆”问题，使得多轮对话和多步推理的数据流转变得清晰。

2.  **Nodes (节点):**

    - 节点就是 Python 函数。
    - 通常包含：`Agent` 节点（调用 LLM 决策）、`Tool` 节点（执行具体工具）、`Human` 节点（人工介入）。

3.  **Edges (边):**

    - 定义控制流。
    - **Conditional Edges (条件边):** 这是 LangGraph 的精髓。例如：LLM 输出结果后，通过条件边判断是“结束”还是“调用工具”还是“返回修改”。这允许了**循环（Cycles）**的存在。

4.  **Checkpointer (检查点/持久化):**
    - LangGraph 内置了持久化机制。
    - 它可以保存图在任意步骤的状态。这意味着你可以实现“时光倒流”（Time Travel），查看历史状态，甚至修改状态后从中间步骤重新运行。这对于 Human-in-the-loop（人机交互）至关重要。

#### 为什么需要 LangGraph？(对比 LangChain AgentExecutor)

- **循环能力:** 能够处理 "Plan -> Execute -> Reflect -> Re-plan" 这种迭代优化的工作流。
- **细粒度控制:** 你不再依赖黑盒的 `run()` 方法，而是显式定义每一步的逻辑。
- **多智能体协作 (Multi-Agent):** 可以轻松构建多个 Agent（如一个负责写代码，一个负责测试），它们作为图中的不同节点互相传递消息。

---

### 3. 代码对比示例

#### LangChain (线性链式思维)

这是一个简单的 LCEL (LangChain Expression Language) 链：

```python
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI
from langchain_core.output_parsers import StrOutputParser

model = ChatOpenAI(model="gpt-4")
prompt = ChatPromptTemplate.from_template("给我讲一个关于 {topic} 的笑话")
output_parser = StrOutputParser()

# 线性流程：Prompt -> Model -> Parser
chain = prompt | model | output_parser

result = chain.invoke({"topic": "程序员"})
print(result)
```

#### LangGraph (图式/循环思维)

这是一个简单的 Agent 循环：如果模型决定调用工具，它会循环回到模型节点。

```python
from typing import TypedDict, Annotated, Sequence
import operator
from langchain_core.messages import BaseMessage, HumanMessage
from langchain_openai import ChatOpenAI
from langgraph.graph import StateGraph, END
from langgraph.prebuilt import ToolNode

# 1. 定义状态
class AgentState(TypedDict):
    # messages 列表会被追加而不是覆盖
    messages: Annotated[Sequence[BaseMessage], operator.add]

# 2. 定义工具和模型
def multiply(a: int, b: int) -> int:
    """Multiply two numbers."""
    return a * b

tools = [multiply]
tool_node = ToolNode(tools)
model = ChatOpenAI(model="gpt-4").bind_tools(tools)

# 3. 定义节点函数
def call_model(state):
    messages = state['messages']
    response = model.invoke(messages)
    return {"messages": [response]}

def should_continue(state):
    last_message = state['messages'][-1]
    # 如果模型返回了 tool_calls，则进入工具节点，否则结束
    if last_message.tool_calls:
        return "tools"
    return END

# 4. 构建图
workflow = StateGraph(AgentState)

workflow.add_node("agent", call_model)
workflow.add_node("tools", tool_node)

workflow.set_entry_point("agent")

# 添加条件边：agent -> (判断) -> tools 或 END
workflow.add_conditional_edges(
    "agent",
    should_continue,
    {
        "tools": "tools",
        END: END
    }
)

# 添加普通边：工具执行完后，必须回到 agent 重新思考
workflow.add_edge("tools", "agent")

# 5. 编译并运行
app = workflow.compile()

inputs = {"messages": [HumanMessage(content="3 乘以 4 是多少？")]}
result = app.invoke(inputs)
```

### 4. 总结与选型建议

| 特性         | LangChain (LCEL)                 | LangGraph                                       |
| :----------- | :------------------------------- | :---------------------------------------------- |
| **核心结构** | DAG (有向无环图)，流水线         | Graph (图)，包含循环                            |
| **状态管理** | 较弱，通常在内存中传递           | 强，显式的 State Schema，支持持久化             |
| **控制流**   | 主要是线性的                     | 支持条件分支、循环、回退                        |
| **调试难度** | 简单链容易，复杂 Agent 难        | 可视化图结构，每一步状态可查                    |
| **适用场景** | 简单 RAG，单次调用，数据处理管道 | 复杂 Agent，多 Agent 协作，人机交互，长流程任务 |

**一句话建议：** 如果你只是想做一个简单的文档问答机器人，用 **LangChain**；如果你要构建一个能够自我修正代码、进行多步推理或需要人工审批流程的智能助手，必须使用 **LangGraph**。
