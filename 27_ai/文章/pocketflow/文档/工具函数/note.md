## LLM Wrapper

```py
from functools import lru_cache

@lru_cache(maxsize=1000)
def cached_call(prompt):
    pass

def call_llm(prompt, use_cache):
    if use_cache:
        return cached_call(prompt)
    # Call the underlying function directly
    return cached_call.__wrapped__(prompt)

class SummarizeNode(Node):
    def exec(self, text):
        # 仅在不重试时使用缓存结果
        return call_llm(f"Summarize: {text}", self.cur_retry==0)
```

LLM Wrappers (LLM 包装器) 深度解析与内化

### 一、 核心价值：抽象与解耦 (Abstraction & Decoupling)

**LLM 包装器是 Agentic Coding 中最重要的 Utilities 之一。**

1.  **屏蔽差异：** 不同厂商（OpenAI, Anthropic, Google, 本地模型）的 SDK 调用方式各异。如果在每个 Node 里都写特定的 SDK 代码，系统将变得极其脆弱且难以迁移。
2.  **统一入口：** 通过 `call_llm` 这个统一的函数签名，我们可以在一处修改，全局生效（例如全局更换模型、全局增加日志、全局增加异常处理）。
3.  **开发效率：** 开发者无需反复翻阅各大厂的 API 文档，只需调用包装好的函数即可。

### 二、 主流厂商实现分析 (内化速查表)

这部分展示了各大模型的“最小可行实现”。

| 厂商                   | 关键库         | 核心对象                              | 注意事项                                                                               |
| :--------------------- | :------------- | :------------------------------------ | :------------------------------------------------------------------------------------- |
| **OpenAI**             | `openai`       | `client.chat.completions.create`      | 事实上的工业标准接口。DeepSeek 等许多第三方也兼容此格式。                              |
| **Anthropic (Claude)** | `anthropic`    | `client.messages.create`              | 响应结构略有不同 (`r.content[0].text`)。目前在代码能力上可与 GPT-4 媲美。              |
| **Google (Gemini)**    | `google-genai` | `client.models.generate_content`      | 接口风格较独特。模型名称如 `gemini-2.5-pro`。                                          |
| **Azure OpenAI**       | `openai`       | `AzureOpenAI`                         | 需要配置 Endpoint 和 Deployment Name。企业级使用的首选。                               |
| **DeepSeek**           | `openai`       | `base_url="https://api.deepseek.com"` | **高性价比之选**。直接复用 OpenAI 库，只需改 `base_url` 和 `api_key`，无需学习新 SDK。 |
| **Ollama (本地)**      | `ollama`       | `chat`                                | 离线开发神器。零成本，不泄露隐私。                                                     |

### 三、 进阶封装策略 (Improvements)

一个工业级的 `call_llm` 绝不仅仅是 API 转发。

#### 1. 缓存策略 (Caching) —— 省钱加速神器

- **为什么需要？** 开发调试时，同样的 Prompt 会跑无数遍。如果不缓存，就是纯粹的烧钱和浪费时间。
- **痛点 (Trap)：** 简单的 `@lru_cache` 会破坏 Agent 的自我修正机制。
  - _场景：_ Agent 写错了 -> Retry 触发 -> 再次调用 call_llm。
  - _问题：_ 如果全缓存了，再次调用还是返回之前那个“错误的答案”，死循环。
- **解法：** **条件缓存**。
  ```python
  def call_llm(prompt, use_cache=True):
      if use_cache:
          return _cached_impl(prompt)
      else:
          # 强制穿透缓存，重跑（例如在 Retry 时）
          return _cached_impl.__wrapped__(prompt)
  ```
  在 Node 内部，只有 `self.cur_retry == 0` (第一次尝试) 时才开启缓存。

#### 2. 日志记录 (Logging) —— 调试之眼

- **LLM 是个黑盒**。没有详细的日志（Prompt 是什么？Response 是什么？耗时多少？），出了问题（如幻觉、截断）根本没法查。
- **最佳实践：** 记录 I/O 对，最好还能记录 Token 消耗量。

#### 3. 支持消息历史 (History Support)

- 简单的 `call_llm(prompt: str)` 只能处理单轮指令。
- 更强大的 `call_llm(messages: list)` 允许传入 `[{"role": "system", ...}, {"role": "user", ...}]`，这对于构建**有记忆的 ChatBot** 或 **角色扮演 Agent** 至关重要。

### 四、 终极建议：拥抱标准，保留退路

1.  **Litellm 库：** 如果你的项目需要频繁切换几十种模型，别自己手写包装器了，直接用 `litellm`。它把所有厂商的接口都拉平成了 OpenAI 格式。
2.  **环境变量：** 永远不要把 Key 硬编码在代码里。永远使用 `os.getenv("OPENAI_API_KEY")`。
3.  **从简单开始：** 先用 `call_llm(prompt)` 跑通 MVP，等需要多轮对话时再升级成 `call_llm(messages)`。不要过度设计。

## Viz and Debug

- 不提供内置的可视化和调试功能
- 使用 Mermaid 进行可视化
  https://github.com/The-Pocket/PocketFlow/tree/main/cookbook/pocketflow-visualization
  递归遍历嵌套图形，为每个节点分配唯一 ID，并将 Flow 节点视为子图，以生成 Mermaid 语法用于分层可视化。
  ```py
  def build_mermaid(start):
    ids, visited, lines = {}, set(), ["graph LR"]
    ctr = 1
    def get_id(n):
        nonlocal ctr
        return ids[n] if n in ids else (ids.setdefault(n, f"N{ctr}"), (ctr := ctr + 1))[0]
    def link(a, b):
        lines.append(f"    {a} --> {b}")
    def walk(node, parent=None):
        if node in visited:
            return parent and link(parent, get_id(node))
        visited.add(node)
        if isinstance(node, Flow):
            node.start_node and parent and link(parent, get_id(node.start_node))
            lines.append(f"\n    subgraph sub_flow_{get_id(node)}[{type(node).__name__}]")
            node.start_node and walk(node.start_node)
            for nxt in node.successors.values():
                node.start_node and walk(nxt, get_id(node.start_node)) or (parent and link(parent, get_id(nxt))) or walk(nxt)
            lines.append("    end\n")
        else:
            lines.append(f"    {(nid := get_id(node))}['{type(node).__name__}']")
            parent and link(parent, nid)
            [walk(nxt, nid) for nxt in node.successors.values()]
    walk(start)
    return "\n".join(lines)
  ```
- 调用栈调试

```py
import inspect

def get_node_call_stack():
    stack = inspect.stack()
    node_names = []
    seen_ids = set()
    for frame_info in stack[1:]:
        local_vars = frame_info.frame.f_locals
        if 'self' in local_vars:
            caller_self = local_vars['self']
            if isinstance(caller_self, BaseNode) and id(caller_self) not in seen_ids:
                seen_ids.add(id(caller_self))
                node_names.append(type(caller_self).__name__)
    return node_names
```

## Web Search

## Text Chunking

文本分块 (Text Chunking) 深度解析与内化

### 一、 核心定位：微优化 (Micro-Optimization)

**请记住这个优先级：先有好的架构（Flow Design），再有好的分块（Chunking）。**
如果你整体的 RAG 流程或 Map Reduce 逻辑是混乱的，把 Chunking 做得再完美也于事无补。先让系统跑通（Make it work），再让分块精准（Make it better）。

### 二、 常用策略深度横评

分块的本质是**在“上下文完整性”和“检索粒度”之间做权衡**。

#### 1. 朴素（固定大小）分块 (Naive Chunking)

- **原理：** 就像切香肠一样，不管切到哪里，每 100 个字符切一刀。
- **代码逻辑：**
  ```python
  def fixed_size_chunk(text, chunk_size=100):
      # 简单粗暴，不管语义，不管是句号还是逗号
      return [text[i:i+chunk_size] for i in range(0, len(text), chunk_size)]
  ```
- **优点：** 极快，零依赖，即使是很烂的脏数据也能跑。
- **缺点致命：** **断章取义**。如 "I love apple p" 和 "ie." 被切到了两块里，语义完全丢失。这对 Embedding 模型的理解是毁灭性的打击。
- **最佳实践：** 仅用于早期的快速原型验证 (Quick Prototype)。

#### 2. 基于句子的分块 (Sentence-Based Chunking)

- **原理：** 尊重语言自然单位。先把文本拆成句子，再按固定句子数量打包。
- **代码逻辑：** 使用 `nltk` 或 `spacy` 这种 NLP 库。
  ```python
  import nltk # 需要预先下载 punkt 模型
  def sentence_based_chunk(text, max_sentences=2):
      sentences = nltk.sent_tokenize(text)
      # 每 2 句合成一段
      return [" ".join(sentences[i:i+max_sentences]) for i in range(0, len(sentences), max_sentences)]
  ```
- **优点：** 保证了最小语义单元（句子）是不被切断的。
- **缺点：**
  - **长难句灾难：** 有些法律文书一句话就 500 个词，可能撑爆 chunk size。
  - **上下文丢失：** "他同意了。" 这句话如果不和上一句话放在一起，毫无意义。

#### 3. 进阶策略 (Advanced Strategies)

- **基于段落 (Paragraph-Based):**
  - _逻辑：_ 按 `\n\n` 切分。
  - _适用：_ 结构良好的文章。但如果遇到一大段不换行的排版，效果极差。
- **语义分块 (Semantic Chunking):**
  - _逻辑：_ 两句话的 Embedding 相似度高就归为一块，突变了就切开。
  - _评价：_ 效果最好，但计算成本最高（需要大量调用 Embedding 模型）。
- **重叠分块 (Sliding Window):** **这是工业界的默认标准**。
  - _逻辑：_ 切 500 词，向后滑 100 词。
  - _意义：_ 解决了边界断裂问题，保证每句话至少在完整的上下文中出现一次。

### 三、 给开发者的行动指南

1.  **起步：** 用 **Markdown 分块** (按标题层级) 或 **重叠固定大小分块** (Chunk Size=1000, Overlap=200)。
2.  **避坑：** 不要只从字符数考虑。如果要切代码，要按函数切；要切 JSON，要按对象切。
3.  **内化：** 分块不是目的，让 LLM 读懂才是目的。如果检索出来的块 LLM 直呼“看不懂”，那就该调整 Chunk 策略了。

## Embedding

Embedding (嵌入) 深度解析与架构内化

### 一、 Embedding 的本质：计算机的“翻译器”

如果说 LLM 是“大脑”，那么 Embedding 就是**“神经传导液”**。

1.  **它是什么？** 将人类的自然语言（Word, Sentence）“压缩”成高维空间中的一个坐标点（List of Floats）。
2.  **有什么用？**
    - **语义近似（Semantic Proximity）：** 在向量空间中，"Apple" 和 "Fruit" 的距离很近，和 "Car" 的距离很远。
    - **计算而非匹配：** 传统搜索匹配关键词（字面匹配），Embedding 搜索匹配意图（语义匹配）。这使得计算机能通过数学公式（如余弦相似度）计算出“这两句话说的是不是一回事”。

### 二、 架构视野下的“微优化”定位

不要在项目一开始就在 Embedding 的选型上纠结太久。

- **初期原则：** 选最快、最稳定、最通用的（如 OpenAI `text-embedding-ada-002` 或 `text-embedding-3-small`）。
- **后期优化：** 当你发现 RAG 检索不准，或者成本太高时，再考虑更换。
  - _要找多语言支持？_ 换 Cohere 或 Google Vertex。
  - _数据极其敏感？_ 此类情况应考虑私有化部署（HuggingFace 本地模型）。
  - _长文本需求？_ Jina 或 Cohere 等提供更长窗口的模型。

### 三、 厂商生态大比拼 (内化决策表)

| 厂商                 | 特点一句话评价                                                                   | 适用场景                               |
| :------------------- | :------------------------------------------------------------------------------- | :------------------------------------- |
| **OpenAI**           | **行业标杆，万金油** 。即开即用，性能均衡，如果不缺钱也不想折腾，选它。          | 快速原型 (MVP)、通用 RAG 系统。        |
| **Azure OpenAI**     | **企业级安全盾**。底层同 OpenAI，但多了合规性和 SLA 保障。                       | 大企业内部项目，对数据隐私有硬性要求。 |
| **Google Vertex AI** | **多模态与长文本**。背靠 Google 强大的搜索算法积淀，多语言处理能力强。           | 全球化多语言业务。                     |
| **AWS Bedrock**      | **云原生集成**。如果你的架构全泡在 AWS 上，用它最省心（Titan 模型便宜大碗）。    | 深度绑定 AWS 生态的后端服务。          |
| **Hugging Face**     | **开源与自由**。可以调用 Inference API，也可以下载到本地跑（完全免费，只费电）。 | 离线环境、极致成本控制、学术研究。     |

### 四、 代码实现的“统一接口”思想

在 `utils/get_embedding.py` 中，建议**封装**一层。不要把 `openai.create` 直接写死在 Node 里。因为你随时可能想换个便宜厂商试试。

**参考封装模式 (Utility Wrapper):**

```python
# utils/get_embedding.py
import os
import numpy as np
from openai import OpenAI

# 环境变量控制，方便切换
PROVIDER = os.getenv("EMBEDDING_PROVIDER", "openai")

def get_embedding(text: str) -> list[float]:
    """统一的 Embedding 接口，对上层屏蔽具体厂商差异"""

    if PROVIDER == "openai":
        client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))
        # 注意: 这里的 model 也可以做成配置项
        resp = client.embeddings.create(model="text-embedding-3-small", input=text)
        return resp.data[0].embedding

    elif PROVIDER == "vertex":
        # ... Google Implementation ...
        pass

    # ... 其他厂商 ...

    else:
        raise ValueError(f"Unknown provider: {PROVIDER}")

if __name__ == "__main__":
    # 自测代码
    vec = get_embedding("Test embedding wrapper")
    print(f"Provider: {PROVIDER}, Dimension: {len(vec)}")
```

**内化总结：**
Embedding 是 RAG 系统的基石。基石歪了（Embedding 质量差），上层建筑（LLM 回答）一定会塌。但在开发初期，**用现成的 API 远比自己训练或微调模型重要**。

## VectorDB

向量数据库 (Vector Databases) 深度解析与架构内化

### 一、 核心定位：Agent 的“海马体” (长期记忆)

如果说 LLM 是 CPU（处理器），Context Window 是 RAM（内存），那么 **Vector Database 就是硬盘（长期存储）**。

在 Agentic Coding 架构中，向量数据库是 **RAG (Retrieval Augmented Generation)** 模式的物理载体。它的核心价值在于：

1.  **突破 Token 限制：** 存储无限的知识，按需提取。
2.  **语义检索：** 即使关键词不匹配，只要意思相近（向量距离近），也能搜到。

### 二、 选型决策矩阵 (Decision Matrix)

面对这么多选项，不要迷茫。根据项目阶段和需求进行“漏斗式”选择：

| 场景                          | 推荐工具                | 理由 (Why?)                                                                                      |
| :---------------------------- | :---------------------- | :----------------------------------------------------------------------------------------------- |
| **MVP / 快速原型 / 个人项目** | **Chroma** 或 **FAISS** | **零成本，零运维**。直接嵌入在代码里作为库运行，不需要部署服务器。数据存在本地文件即可。         |
| **生产环境 (Managed Cloud)**  | **Pinecone**            | **省心，贵**。完全托管，API 最简单，不需要关心底层索引构建，专注于业务逻辑。                     |
| **生产环境 (复杂元数据过滤)** | **Qdrant**              | **平衡之选**。在 Filter（元数据过滤）和 Vector Search 的结合上做得非常好，且 Rust 内核性能极高。 |
| **超大规模 (十亿级向量)**     | **Milvus**              | **重型武器**。云原生架构，Kubernetes 友好，适合数据量巨大的企业级应用。                          |
| **既要速度又要缓存**          | **Redis**               | **速度魔鬼**。如果你已经在用 Redis 做缓存，直接开启 Vector 模块是最快路径。                      |

### 三、 代码模式内化：统一的 CRUD 接口

无论选哪个库，操作逻辑在 **Utility** 层是一致的。**千万不要把特定的数据库 SDK 代码散落在各个 Node 里。**

请在 `utils/vector_store.py` 中封装如下通用逻辑：

1.  **Connect (连接/初始化):** 建立连接，或者加载本地索引。
2.  **Upsert (更新/插入):** `(ID, Vector, Metadata)` 三元组入库。
    - _注意：_ 大多数 DB 要求向量维度（Dimension）必须固定（如 OpenAI 是 1536，本地模型可能是 768）。
3.  **Search (查询):** `(QueryVector, TopK)` -> 返回 `List[ChunkText]`。

**抽象封装示例：**

```python
# utils/vector_store.py
import os
import chromadb

# 抽象层：如果明天想换 Pinecone，只需要改这个文件的实现，Nodes 不需要动
class VectorStore:
    def __init__(self):
        # 默认使用本地 Chroma，轻量级
        self.client = chromadb.PersistentClient(path="./chroma_db")
        self.collection = self.client.get_or_create_collection("knowledge_base")

    def add_documents(self, documents: list[str], embeddings: list[list[float]], ids: list[str]):
        self.collection.upsert(
            documents=documents,
            embeddings=embeddings,
            ids=ids
        )

    def search(self, query_embedding: list[float], top_k=3):
        results = self.collection.query(
            query_embeddings=[query_embedding],
            n_results=top_k
        )
        # 扁平化返回结果
        return results["documents"][0]
```

### 四、 避坑指南 (Best Practices)

1.  **维度陷阱：** 建库时（Create Index）填写的维度必须和你的 Embedding 模型（如 text-embedding-3-small）输出维度完全一致，否则报错。
2.  **元数据的重要性：** 存向量时，务必把**源文本 (Raw Text)** 存入 Metadata 或 Payload 中。因为向量数据库搜出来的是向量，我们需要反查出文本喂给 LLM。
    - _错误做法：_ DB 里只存向量和 ID -> 搜到 ID -> 去 MySQL 查文本 -> 慢！
    - _正确做法：_ DB 里存向量 + 文本 -> 搜到直接返回文本 -> 快！
3.  **冷启动问题：** FAISS 和 Chroma 在数据量极大时做全量加载会很慢，生产环境请务必使用 Client/Server 模式 (Qdrant/Milvus/Pinecone)。

**一句话总结：**
对于 Agent 开发，**Chroma 是最好的起步板，Pinecone 是最舒服的云服务，Qdrant 是最硬核的自托管方案。** 先跑通，再扩容。

## Text-to-Speech

Text-to-Speech (TTS) 深度解析与架构内化

### 一、 核心定位：Agent 的“声带”

在智能代理系统中，TTS (Text-to-Speech) 是 **输出层 (Output Utility)** 的关键组件。
如果说 LLM 负责“思考”和“生成文本”，那么 TTS 则负责将这些无声的字符转化为**情感化的语音流**。这对于电话机器人、语音助手、游戏 NPC 或有声书生成 Agent 至关重要。

### 二、 选型决策矩阵 (Decision Matrix)

市场目前的格局非常清晰：**“云巨头”主打性价比与稳定性，“垂直独角兽”主打拟真度。**

| 类别           | 推荐服务                     | 核心优势                                                                                                  | 适用场景                               |
| :------------- | :--------------------------- | :-------------------------------------------------------------------------------------------------------- | :------------------------------------- |
| **性价比之王** | **Azure TTS**                | **综合最强**。50 万字符/月的免费额度非常良心，且 Neural 语音质量已经是云厂商中的第一梯队。                | 客服机器人、即时通知、长文本朗读。     |
| **质量天花板** | **ElevenLabs**               | **拟真度极高**。拥有惊人的情感表现力和声音克隆能力。它听起来不像机器，像人。但价格极其昂贵。              | 游戏 NPC、播客生成、虚拟人、情感陪伴。 |
| **生态绑定**   | **AWS Polly / Google Cloud** | **基建成熟**。如果你已经在用 AWS/GCP 全家桶，集成它们是最顺手的。Polly 的标准版非常便宜，适合大批量处理。 | 内部系统通知、海量文档转语音。         |
| **离线/隐私**  | **Local TTS (如 Coqui)**     | **隐私安全**。无需联网，数据不出域。                                                                      | 边缘设备、隐私敏感项目。               |

### 三、 架构设计：统一语音接口

不要在 Node 里直接写特定的 SDK 代码。这会导致系统与特定厂商绑定（Vendor Lock-in）。
**正确的 Agentic Coding 做法是：封装一个与厂商无关的通用工具函数。**

#### 建议文件结构：`utils/tts_service.py`

```python
# utils/tts_service.py
import os
import requests
import boto3
# 其他 SDK 按需导入...

# 环境变量：TTS_PROVIDER = "azure" / "elevenlabs" / "polly"

def generate_speech(text: str, output_file: str = "output.mp3", voice_id: str = None):
    provider = os.getenv("TTS_PROVIDER", "elevenlabs")

    if provider == "elevenlabs":
        _call_elevenlabs(text, output_file, voice_id)
    elif provider == "polly":
        _call_polly(text, output_file, voice_id)
    elif provider == "azure":
        _call_azure(text, output_file, voice_id)
    else:
        raise ValueError(f"Unknown TTS provider: {provider}")

def _call_elevenlabs(text, filepath, voice_id):
    # 默认 Voice ID
    vid = voice_id or "21m00Tcm4TlvDq8ikWAM"
    api_key = os.getenv("ELEVENLABS_API_KEY")

    url = f"https://api.elevenlabs.io/v1/text-to-speech/{vid}"
    headers = {
        "xi-api-key": api_key,
        "Content-Type": "application/json"
    }
    data = {"text": text, "voice_settings": {"stability": 0.5, "similarity_boost": 0.75}}

    resp = requests.post(url, headers=headers, json=data)
    if resp.status_code == 200:
        with open(filepath, "wb") as f:
            f.write(resp.content)
    else:
        print(f"Error: {resp.text}")

def _call_polly(text, filepath, voice_id):
    client = boto3.client("polly", region_name="us-east-1")
    resp = client.synthesize_speech(
        Text=text, OutputFormat="mp3", VoiceId=voice_id or "Joanna"
    )
    with open(filepath, "wb") as f:
        f.write(resp["AudioStream"].read())

# ... 其他实现 ...
```

### 四、 避坑与微优化 (Optimization)

1.  **缓存 (Caching) 是刚需：**
    - TTS 的 API 调用通常比 LLM 更慢，且按字符收费。
    - 对于重复的固定语料（如“欢迎光临”、“请稍候”），**必须**将生成的 MP3 缓存到本地文件系统或 S3。第二次直接播放文件，不要调 API。
2.  **流式传输 (Streaming)：**
    - 在实时对话 Agent 中，生成完整个音频文件再播放会造成巨大的延迟。
    - 应优先选择支持 **Stream** 模式的 API（ElevenLabs 和 OpenAI TTS 都支持），边生成边播放，让用户感觉响应很快。
3.  **文本预处理：**
    - LLM 生成的 Markdown 符号（如 `**加粗**`、`### 标题`）直接喂给 TTS 会被读出来，非常尴尬。
    - **关键步骤：** 在送入 TTS 之前，必须有一个 `clean_text_for_speech()` 函数，去除 Markdown 标记和特殊符号。

**一句话总结：**
日常使用选 Azure（便宜好用），追求惊艳选 ElevenLabs（贵但在理）。**务必封装接口，务必做文件缓存。**
