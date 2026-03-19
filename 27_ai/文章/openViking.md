## OpenViking 代码架构全面讲解

OpenViking 的灵魂是 "用文件系统的目录结构组织上下文 + 用优先队列递归钻取替代扁平向量检索 + 用 6 类记忆自动提取让 Agent 自进化"。

OpenViking 是字节跳动（火山引擎）开源的 **AI Agent 上下文数据库**，核心思想是将 Agent 所需的所有上下文（Memory、Resource、Skill）统一抽象为**虚拟文件系统**，支持语义检索和渐进式内容加载。

---

### 1. 项目整体技术栈

| 语言           | 占比  | 用途                                     |
| -------------- | ----- | ---------------------------------------- |
| **Python**     | 82.7% | 核心 SDK、服务端、解析器、检索、会话管理 |
| **C++**        | 8.4%  | 高性能扩展（pybind11 绑定）              |
| **Rust**       | 3.6%  | CLI 客户端 (`ov` 命令)                   |
| **JavaScript** | 2.0%  | VikingBot bridge (WhatsApp 等)           |

---

### 2. 顶层目录结构

```
openviking/
├── pyproject.toml          # Python 项目配置
├── Cargo.toml              # Rust workspace 配置
├── Makefile                # 构建入口
├── Dockerfile              # Docker 部署
│
├── openviking/             # 👈 Python 核心包（最重要）
│   ├── async_client.py     # 异步客户端
│   ├── sync_client.py      # 同步客户端
│   ├── core/               # 核心数据模型
│   ├── parse/              # 资源解析
│   ├── retrieve/           # 检索系统
│   ├── session/            # 会话管理
│   ├── server/             # HTTP 服务端 (FastAPI)
│   ├── console/            # Web Console
│   ├── storage/            # 存储层
│   ├── utils/              # 工具类
│   └── prompts/            # LLM 提示词模板
│
├── openviking_cli/         # Python CLI 工具和配置
├── crates/ov_cli/          # Rust CLI（高性能命令行）
├── src/                    # C++ 扩展 (pybind11)
├── third_party/agfs/       # AGFS 文件系统（核心依赖）
├── bot/                    # VikingBot（AI Agent 框架）
├── tests/                  # 测试套件
├── docs/                   # 文档（en/zh）
└── examples/               # 使用示例
```

---

### 3. 系统架构 — 分层设计

```
┌──────────────────────────────────────────────────────────┐
│                     Client Layer                          │
│              OpenViking / SyncHTTPClient                   │
│              (统一入口，委托给 Service 层)                   │
└──────────────────┬───────────────────────────────────────┘
                   │ delegates
┌──────────────────▼───────────────────────────────────────┐
│                    Service Layer                          │
│  FSService | SearchService | SessionService               │
│  ResourceService | RelationService | PackService          │
│  DebugService (ObserverService)                           │
└───┬──────────────┬──────────────────┬────────────────────┘
    │              │                  │
    ▼              ▼                  ▼
┌────────┐  ┌──────────┐      ┌──────────┐
│Retrieve│  │ Session  │      │  Parse   │
│上下文检索│  │ 会话管理  │      │ 上下文提取│
│        │  │          │      │          │
│search  │  │add/commit│      │文档解析   │
│意图分析  │  │  used    │      │L0/L1/L2  │
│Rerank  │  │          │      │树构建     │
└───┬────┘  └────┬─────┘      └────┬─────┘
    │            │                  │
    │       ┌────▼─────┐           │
    │       │Compressor│           │
    │       │记忆压缩/去重│           │
    │       └────┬─────┘           │
    └────────────┼─────────────────┘
                 ▼
┌──────────────────────────────────────────────────────────┐
│                    Storage Layer                          │
│         AGFS (文件内容) + Vector Index (索引)               │
└──────────────────────────────────────────────────────────┘
```

---

### 4. 核心模块详解

#### 4.1 Client（客户端入口）

- `openviking/sync_client.py` — 同步客户端 `OpenViking`
- `openviking/async_client.py` — 异步客户端 `AsyncOpenViking`

两种部署模式：

- **嵌入式模式**：`client = OpenViking(path="./data")`，自动启动 AGFS 子进程，单例模式
- **HTTP 模式**：`client = SyncHTTPClient(url="http://localhost:1933")`，连接独立部署的 Server

#### 4.2 Service Layer（业务逻辑层）

Service 层解耦了业务逻辑和传输层，HTTP Server 和 CLI 复用同一套逻辑：

| Service             | 职责             | 关键方法                                                                                  |
| ------------------- | ---------------- | ----------------------------------------------------------------------------------------- |
| **FSService**       | 虚拟文件系统操作 | `ls`, `mkdir`, `rm`, `mv`, `tree`, `stat`, `read`, `abstract`, `overview`, `grep`, `glob` |
| **SearchService**   | 语义搜索         | `search`, `find`                                                                          |
| **SessionService**  | 会话管理         | `session`, `sessions`, `commit`, `delete`                                                 |
| **ResourceService** | 资源导入         | `add_resource`, `add_skill`, `wait_processed`                                             |
| **RelationService** | 关联管理         | `relations`, `link`, `unlink`                                                             |
| **PackService**     | 导入导出打包     | `export_ovpack`, `import_ovpack`                                                          |
| **DebugService**    | 调试/可观测性    | `observer`                                                                                |

#### 4.3 Parse（上下文提取）

```
openviking/parse/
├── parsers/          # 各类解析器实现（PDF/MD/HTML/Code/Image/Audio/Video）
│   └── code/ast/     # 代码 AST 骨架提取
│       ├── extractor.py    # 语言检测+分发
│       ├── skeleton.py     # 数据结构定义
│       └── languages/      # 各语言专属提取器
├── tree_builder.py   # 树构建器：将解析结果组织为目录树
└── registry.py       # 解析器注册表
```

**解析流程**：输入文档 → Parser 解析为结构化内容 → TreeBuilder 构建为虚拟目录树 → 入队异步语义处理

#### 4.4 Core（核心数据模型）

```python
# openviking/core/
├── context.py        # Context 基类、ResourceContentType
├── building_tree.py  # BuildingTree 构建中的树
├── directories.py    # 预设目录定义（resources/user/agent）
└── skill_loader.py   # 技能加载器
```

核心概念是 `viking://` URI 协议下的三大上下文类型：

```
viking://
├── resources/        # 资源：项目文档、代码仓库、网页等
├── user/             # 用户：个人偏好、习惯等记忆
│   └── memories/
└── agent/            # Agent：技能、指令、任务记忆
    ├── skills/
    ├── memories/
    └── instructions/
```

#### 4.5 L0/L1/L2 三层上下文（核心创新）

写入时自动对每个资源生成三层信息：

| 层级   | 名称             | Token 量    | 用途           |
| ------ | ---------------- | ----------- | -------------- |
| **L0** | Abstract（摘要） | ~100 tokens | 快速筛选相关性 |
| **L1** | Overview（概览） | ~2k tokens  | Agent 规划决策 |
| **L2** | Details（详情）  | 完整原文    | 按需深度阅读   |

这避免了向 LLM 一次性灌入所有上下文，实现**按需加载**，大幅节省 Token 开销。

#### 4.6 Retrieve（检索系统）

```
openviking/retrieve/
├── retriever.py          # HierarchicalRetriever 主检索器
├── intent_analyzer.py    # 意图分析器
└── reranker.py           # 重排序器
```

**检索四步流程**：

1. **意图分析**：LLM 分析查询意图，生成 0-5 个类型化子查询
2. **层级检索**：目录级递归搜索，用优先队列定位高得分目录，再向下钻取
3. **Rerank**：标量过滤 + 模型重排序精选最相关结果
4. **返回**：按相关性排序的上下文

这就是所谓的"**目录递归检索**"策略：先锁定高分目录，再精细探索内容。

#### 4.7 Session（会话管理）

```
openviking/session/
├── session.py       # 会话核心逻辑
└── compressor.py    # 记忆压缩器
```

Session 支持 Agent 的多轮对话管理，核心流程：

- **消息记录**：累积 `user`/`assistant` 消息，支持 TextPart、ContextPart、ToolPart
- **压缩归档**：保留最近 N 轮，旧消息自动归档并生成结构化摘要
- **记忆提取**：Compressor 从对话中提取 **6 种分类记忆**（偏好、事件、实体、决策等），自动去重后写入 `viking://user/memories/` 和 `viking://agent/memories/`

这实现了 Agent **越用越聪明**的自进化能力。

#### 4.8 Storage（双层存储）

```
openviking/storage/
├── viking_fs.py      # VikingFS 虚拟文件系统接口
└── vectordb/
    ├── project/      # 项目管理（ProjectGroup/IProject）
    │   ├── project.py          # IProject 抽象接口
    │   ├── project_group.py    # 多项目管理
    │   ├── local_project.py    # 本地项目实现
    │   └── vikingdb_project.py # 远程 VikingDB 项目
    └── service/
        ├── api_fastapi.py       # VectorDB REST API
        └── server_fastapi.py    # FastAPI 服务启动
```

**双层分离设计**：

| 层               | 技术                                 | 存储内容                                    |
| ---------------- | ------------------------------------ | ------------------------------------------- |
| **AGFS**         | 第三方文件系统 (`third_party/agfs/`) | L0/L1/L2 完整内容、多媒体文件、关联关系     |
| **Vector Index** | 本地/VikingDB 向量库                 | URI、embedding 向量、元数据（**不存内容**） |

设计原则：**单一数据源**，所有内容从 AGFS 读取，向量库只存引用和索引。

#### 4.9 Server（HTTP 服务端）

```
openviking/server/
├── app.py           # FastAPI 应用工厂 create_app()
├── bootstrap.py     # 启动入口（openviking-server 命令）
└── routers/         # API 路由

# API 端点概览：
GET  /api/v1/fs/ls          # 列目录
GET  /api/v1/fs/tree        # 目录树
GET  /api/v1/content/read   # 读 L2 全文
GET  /api/v1/content/abstract  # 读 L0 摘要
GET  /api/v1/content/overview  # 读 L1 概览
POST /api/v1/search/find    # 语义搜索
POST /api/v1/resource/add   # 添加资源
...
```

---

### 5. 三大数据流

#### 5.1 添加上下文

```
输入(URL/文件/目录) → Parser(解析) → TreeBuilder(树构建) → AGFS(持久化)
                                                            → SemanticQueue(异步) → 生成 L0/L1
                                                            → Vector Index(建索引)
```

#### 5.2 检索上下文

```
用户查询 → IntentAnalyzer(意图分析) → HierarchicalRetriever(层级递归检索)
         → Reranker(重排序) → 返回排序后的上下文
```

#### 5.3 会话提交

```
累积消息 → Compressor(压缩) → 归档旧消息(生成L0/L1)
        → 记忆提取(6类分类) → LLM 去重 → 写入 AGFS + Vector Index
```

---

### 6. CLI 双实现

- **Python CLI** (`openviking_cli/`): 配置管理、服务器引导
- **Rust CLI** (`crates/ov_cli/`): 高性能命令行客户端 `ov`，支持 `ls`、`tree`、`find`、`add-resource`、`chat`、`tui` 等子命令

---

### 7. VikingBot（扩展 Agent 框架）

```
bot/
├── vikingbot/        # VikingBot 核心
│   └── cli/commands  # CLI 命令入口
├── deploy/           # Docker 部署配置
└── eval/             # 评测（SkillsBench 等）
```

VikingBot 是基于 OpenViking 构建的 AI Agent 框架，通过 `openviking-server --with-bot` 启用，支持交互式 `ov chat`。

---

### 8. 设计原则总结

| 原则             | 说明                                     |
| ---------------- | ---------------------------------------- |
| **文件系统范式** | 用目录结构统一管理所有上下文，告别碎片化 |
| **三层信息**     | L0/L1/L2 渐进加载，按需消费 Token        |
| **两阶段检索**   | 向量召回 + Rerank 精排                   |
| **存储层纯粹**   | AGFS 管内容，向量库管索引，职责分离      |
| **单一数据源**   | 所有内容从 AGFS 读取，保证一致性         |

整个项目结构清晰、分层明确，核心创新在于**用文件系统范式统一上下文 + L0/L1/L2 三层按需加载 + 目录递归检索策略**。如果你想深入学习某个模块，建议从 `openviking/core/` 和 `openviking/retrieve/` 开始看起。

---

## OpenViking 的灵魂：三大核心算法的源码深度剖析

OpenViking 的灵魂可以归纳为三个核心机制，它们合在一起让 Agent 实现了 **"精准找到所需上下文 → 越用越聪明"** 的闭环。下面按重要程度逐一深入源码逻辑。

---

## 灵魂一：HierarchicalRetriever — 目录递归检索（最核心）

这是 OpenViking 区别于传统 RAG 的**最关键创新**。传统 RAG 是扁平 topK 向量检索，OpenViking 把上下文组织成目录树，用类似操作系统文件系统的方式**递归钻取**。

### 核心数据结构与常量

```python
# openviking/retrieve/hierarchical_retriever.py

class HierarchicalRetriever:
    MAX_CONVERGENCE_ROUNDS = 3    # Top-K 不变超过 3 轮就停止
    SCORE_PROPAGATION_ALPHA = 0.5 # 50% 自身 embedding 分 + 50% 父目录分
    GLOBAL_SEARCH_TOPK = 5        # 全局搜索只粗召回 5 个
    HOTNESS_ALPHA = 0.2           # 热度加权（最近使用/频繁访问的上下文排名提升）
```

### retrieve() 主流程 — 6 步

```python
async def retrieve(self, query: TypedQuery, ctx, limit=5, mode="thinking", ...):
    # Step 1: 确定搜索根目录
    if target_dirs:
        root_uris = target_dirs
    else:
        root_uris = self._get_root_uris_for_type(query.context_type, ctx)
        # MEMORY → ["viking://user/{user}/memories", "viking://agent/{agent}/memories"]
        # RESOURCE → ["viking://resources"]
        # SKILL → ["viking://agent/{agent}/skills"]

    # Step 2: 全局向量搜索 — 粗召回起始目录
    global_results = await self._global_vector_search(...)

    # Step 3: 合并起始点（全局搜索结果 + 根目录）
    starting_points = self._merge_starting_points(query, root_uris, global_results)

    # Step 4: 🔥 递归搜索（灵魂核心）
    candidates = await self._recursive_search(starting_points, ...)

    # Step 5: 转换为 MatchedContext（加入热度加权）
    matched = await self._convert_to_matched_contexts(candidates, ctx)

    # Step 6: 记录统计指标
    return QueryResult(query=query, matched_contexts=matched[:limit])
```

### \_recursive_search() — 灵魂中的灵魂

这个方法用 **堆（优先队列）+ BFS** 实现目录树的递归向下钻取，同时有**分数传播**和**收敛检测**：

```python
async def _recursive_search(self, starting_points, query, query_vector, ...):
    collected_by_uri = {}   # URI → 候选结果（去重保最高分）
    dir_queue = []          # 最大堆: (-score, uri)
    visited = set()
    prev_topk_uris = set()
    convergence_rounds = 0
    alpha = 0.5  # SCORE_PROPAGATION_ALPHA

    # 初始化：把起始目录推入堆
    for uri, score in starting_points:
        heapq.heappush(dir_queue, (-score, uri))

    while dir_queue:
        temp_score, current_uri = heapq.heappop(dir_queue)
        current_score = -temp_score
        if current_uri in visited:
            continue
        visited.add(current_uri)

        # 🔍 向量搜索当前目录的子节点
        results = await vector_proxy.search_children_in_tenant(
            parent_uri=current_uri,
            query_vector=query_vector,
            ...
        )

        # Rerank 精排（仅 THINKING 模式）
        if self._rerank_client and mode == "thinking":
            documents = [r.get("abstract", "") for r in results]
            query_scores = self._rerank_scores(query, documents, default_scores)

        for r, score in zip(results, query_scores):
            uri = r.get("uri", "")

            # 🔥 核心：分数传播公式
            # final_score = α × 自身embedding分 + (1-α) × 父目录传播分
            final_score = alpha * score + (1 - alpha) * current_score

            if not passes_threshold(final_score):
                continue

            # 去重：保留最高分
            if uri not in collected_by_uri or final_score > collected_by_uri[uri]["_final_score"]:
                r["_final_score"] = final_score
                collected_by_uri[uri] = r

            # 只有目录（L0/L1）才继续递归，文件（L2）是终端节点
            if uri not in visited and r.get("level", 2) != 2:
                heapq.heappush(dir_queue, (-final_score, uri))

        # 🛑 收敛检测：如果 Top-K 连续 3 轮不变，提前停止
        current_topk_uris = {c["uri"] for c in sorted(
            collected_by_uri.values(), key=lambda x: x["_final_score"], reverse=True
        )[:limit]}

        if current_topk_uris == prev_topk_uris and len(current_topk_uris) >= limit:
            convergence_rounds += 1
            if convergence_rounds >= 3:  # MAX_CONVERGENCE_ROUNDS
                break
        else:
            convergence_rounds = 0
            prev_topk_uris = current_topk_uris

    return sorted(collected_by_uri.values(), key=lambda x: x["_final_score"], reverse=True)
```

**为什么这是灵魂？** 因为它解决了传统 RAG 的致命问题：

| 传统 RAG                               | OpenViking                                                       |
| -------------------------------------- | ---------------------------------------------------------------- |
| 扁平向量 topK，语义相似但上下文碎片    | 先定位高分目录，再钻取完整上下文                                 |
| 查询"OAuth 认证"可能找到不相关的 token | 先定位到 `viking://resources/api-docs/`，再精细找 OAuth 相关文件 |
| 无法感知目录层级关系                   | 父目录分数通过 α=0.5 传播给子节点，全局语境感                    |
| 无停止条件，固定 topK                  | 收敛检测，Top-K 稳定后自动终止                                   |

**分数传播公式的直觉**：

$$\text{final\_score} = \alpha \cdot s_{\text{embedding}} + (1 - \alpha) \cdot s_{\text{parent}}$$

一个文件即使自身 embedding 匹配度不高，如果它在一个高分目录下，也会得到较高分数 — 这模拟了人类浏览文件系统时 "进对了文件夹，里面的东西大概率相关" 的直觉。

### 热度加权 — 让结果有时间感知

```python
# _convert_to_matched_contexts 中
semantic_score = c.get("_final_score", 0.0)
h_score = hotness_score(active_count, updated_at)  # 基于访问频率和更新时间
final_score = (1 - 0.2) * semantic_score + 0.2 * h_score
```

这让最近频繁使用和更新的上下文排名更靠前。

---

## 灵魂二：IntentAnalyzer — 意图分析与查询分解

用户一句话可能同时需要 Memory + Resource + Skill，IntentAnalyzer 通过 LLM 将用户查询**分解为 0-5 个类型化子查询**。

### 核心流程

```python
# openviking/retrieve/intent_analyzer.py

class IntentAnalyzer:
    MAX_COMPRESSION_SUMMARY_CHARS = 30000

    async def analyze(self, compression_summary, messages, current_message, ...):
        # 1. 构建 prompt（融合会话摘要 + 最近5条消息 + 当前查询）
        prompt = self._build_context_prompt(
            compression_summary, messages, current_message, context_type, target_abstract
        )

        # 2. 调用 LLM
        response = await get_openviking_config().vlm.get_completion_async(prompt)

        # 3. 解析 JSON 响应
        parsed = parse_json_from_response(response)

        # 4. 构建 TypedQuery 列表
        queries = []
        for q in parsed.get("queries", []):
            queries.append(TypedQuery(
                query=q["query"],         # 重写后的查询
                context_type=ContextType(q["context_type"]),  # memory/resource/skill
                intent=q["intent"],       # 查询意图说明
                priority=q["priority"],   # 1-5
            ))

        return QueryPlan(queries=queries, reasoning=parsed["reasoning"])
```

**关键设计**：

```python
@dataclass
class TypedQuery:
    query: str              # LLM 重写的查询（比原始查询更精准）
    context_type: ContextType  # 决定搜索哪个根目录
    intent: str             # 可追踪的检索理由
    priority: int           # 并发执行时的优先级
    target_directories: List[str]  # LLM 可以直接定位到具体目录
```

**为什么重要？**

- 用户说 "帮我创建一个 RFC 文档" → LLM 分解为：
  - `TypedQuery(query="RFC 文档模板", type=RESOURCE)` — 找模板
  - `TypedQuery(query="创建 RFC 文档", type=SKILL)` — 找技能
  - `TypedQuery(query="用户的文档格式偏好", type=MEMORY)` — 找记忆

- 多个 TypedQuery **并发执行** `asyncio.gather`，每个走独立的 HierarchicalRetriever，最后聚合。

### find() vs search() 的区别

```python
# VikingFS.search() — 复杂查询，有意图分析
async def search(self, query, session_info=None, ...):
    if session_summary or recent_messages:
        # 有会话上下文 → 用 IntentAnalyzer
        analyzer = IntentAnalyzer()
        query_plan = await analyzer.analyze(
            compression_summary=session_summary,
            messages=recent_messages,
            current_message=query,
        )
        typed_queries = query_plan.queries
    else:
        # 无会话上下文 → 直接创建单个 TypedQuery
        typed_queries = [TypedQuery(query=query, context_type=None, ...)]

    # 并发执行所有子查询
    results = await asyncio.gather(*[retriever.retrieve(tq) for tq in typed_queries])

# VikingFS.find() — 简单查询，跳过意图分析
async def find(self, query, ...):
    typed_query = TypedQuery(query=query, context_type=ContextType.RESOURCE, ...)
    result = await retriever.retrieve(typed_query)
```

---

## 灵魂三：SessionCompressor + MemoryExtractor — 记忆自进化

这是让 Agent **越用越聪明**的关键。每次会话结束 `commit()` 时，从对话中**自动提取长期记忆**。

### 6 种记忆分类

```python
class MemoryCategory(str, Enum):
    PROFILE = "profile"           # 用户身份属性（归 user/）
    PREFERENCES = "preferences"   # 用户偏好（归 user/）
    ENTITIES = "entities"         # 人物/项目实体（归 user/）
    EVENTS = "events"             # 事件/决策（归 user/）
    TOOLS = "tools"               # 工具使用经验（归 agent/）
    SKILLS = "skills"             # 技能使用经验（归 agent/）
    CASES = "cases"               # 问题+解法（归 agent/）  [8类]
    PATTERNS = "patterns"         # 可复用模式（归 agent/）
```

### extract_long_term_memories() — 记忆提取主流程

```python
class SessionCompressor:
    async def extract_long_term_memories(self, messages, user, session_id, ctx):
        # Phase 1: LLM 提取候选记忆
        candidates = await self.extractor.extract(
            {"messages": messages}, user, session_id
        )

        for candidate in candidates:
            # --- Profile 类直接创建，不去重 ---
            if candidate.category == MemoryCategory.PROFILE:
                memory = await self.extractor.create_memory(candidate, ...)
                await self._index_memory(memory, ctx)
                continue

            # --- Tool/Skill 类有特殊合并逻辑 ---
            if candidate.category in {TOOLS, SKILLS}:
                # 校准名称：用 ToolPart 的 ground truth 修正 LLM 的猜测
                tool_name, skill_name, status = self._get_tool_skill_info(candidate, tool_parts)
                # 合并到已有的 tool/skill memory 中（累加统计数据）
                memory = await self.extractor._merge_tool_memory(tool_name, candidate, ctx)
                continue

            # --- 其他类别走去重流程 ---
            # Phase 2: 向量搜索找相似记忆
            result = await self.deduplicator.deduplicate(candidate, ctx)
            # result.decision 是：SKIP / CREATE / NONE

            if result.decision == DedupDecision.SKIP:
                # 完全重复，跳过
                continue

            if result.decision == DedupDecision.NONE:
                # 不创建新记忆，但对已有记忆执行操作
                for action in result.actions:
                    if action.decision == MemoryActionDecision.MERGE:
                        # 🔥 合并到已有记忆（通过 LLM 生成合并后的 L0/L1/L2）
                        await self._merge_into_existing(candidate, action.memory, ...)
                    elif action.decision == MemoryActionDecision.DELETE:
                        # 删除过时的冲突记忆
                        await self._delete_existing_memory(action.memory, ...)

            if result.decision == DedupDecision.CREATE:
                # 创建新记忆（可能先删除冲突的旧记忆）
                for action in result.actions:
                    if action.decision == MemoryActionDecision.DELETE:
                        await self._delete_existing_memory(action.memory, ...)
                memory = await self.extractor.create_memory(candidate, ...)
                await self._index_memory(memory, ctx)

        # Phase 3: 建立双向关联
        # 记忆 ↔ 引用的资源/技能
        used_uris = self._extract_used_uris(messages)
        await self._create_relations(memories, used_uris, ctx)
```

### 记忆合并的 LLM 调用

```python
async def _merge_into_existing(self, candidate, target_memory, viking_fs, ctx):
    existing_content = await viking_fs.read_file(target_memory.uri, ctx)
    # 一次 LLM 调用生成合并后的 L0 + L1 + L2
    payload = await self.extractor._merge_memory_bundle(
        existing_abstract=target_memory.abstract,
        existing_content=existing_content,
        new_abstract=candidate.abstract,
        new_content=candidate.content,
        category=candidate.category.value,
    )
    # 写回 AGFS
    await viking_fs.write_file(target_memory.uri, payload.content, ctx)
    target_memory.abstract = payload.abstract
    # 重新向量化
    await self._index_memory(target_memory, ctx)
```

### 去重决策矩阵

| 场景               | LLM 决策            | 操作                         |
| ------------------ | ------------------- | ---------------------------- |
| 完全重复           | `skip`              | 什么都不做                   |
| 新信息，需要新记忆 | `create`            | 创建新 Context，可能先删旧的 |
| 已有记忆需要更新   | `none` + `merge`    | 把新内容融合到已有记忆中     |
| 已有记忆互相冲突   | `none` + `delete`   | 删除过时记忆                 |
| 新记忆 + 删冲突    | `create` + `delete` | 先删后建                     |

---

## 三者如何闭环

```
                    ① Agent 对话
                         │
                         ▼
              ┌──────────────────────┐
              │  IntentAnalyzer      │ ← 理解意图，拆分查询
              │  (0-5 TypedQueries)  │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │ HierarchicalRetriever│ ← 从目录树中精准找上下文
              │ (递归 + 分数传播     │
              │  + 收敛检测)         │
              └──────────┬───────────┘
                         │
                         ▼
              Agent 使用上下文完成任务
                         │
                         ▼
              ┌──────────────────────┐
              │ SessionCompressor    │ ← commit 时提取长期记忆
              │ (6分类 + 去重 + 合并) │
              └──────────┬───────────┘
                         │
                         ▼
              写入 viking://user/memories/
              写入 viking://agent/memories/
                         │
                         ▼
              ┌──────────────────────┐
              │ 下次检索时，记忆被    │ ← 记忆参与未来检索
              │ HierarchicalRetriever│
              │ 检索到并送给 Agent    │
              └──────────────────────┘
```

**一句话总结**：OpenViking 的灵魂是 **"用文件系统的目录结构组织上下文 + 用优先队列递归钻取替代扁平向量检索 + 用 6 类记忆自动提取让 Agent 自进化"**。

传统 RAG 是 "把文本切片灌进向量库，查询时 topK 召回"。OpenViking 是 "把上下文组织成有层级的文件系统，先定位到正确的文件夹，再逐层深入精准找到所需内容，并且每次对话都能学到新记忆"。
