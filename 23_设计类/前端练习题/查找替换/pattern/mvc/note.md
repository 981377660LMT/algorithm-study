# 查找替换模型抽象

1. State
   - State: SearchState 管理了所有搜索结果，包括当前搜索关键词和每个插件的结果。
   - 状态管理：SearchState 是一个普通的类，不依赖于任何状态管理框架，便于测试和迁移。
   - 响应式：插件在数据变化时通知 Controller。
2. Controller
   - SearchController 极其轻量，只维护 currentTask 和 currentIndex。它不知道如何搜索，只知道如何调度。
   - Contoller 可被继承扩展，以支持更多功能。
3. Task
   Task: SearchTask 承载了并发逻辑。通过 searchSingleModule 方法，既支持了初始的全量搜索，也支持了 reSearch 触发的增量更新。
4. Plugin(Module)
   Plugin: 接口强制要求处理 CancellationToken，并在注册时通过闭包/Context 获得了反向调用的能力。
5. View
   View: 视图层（React 组件）只需要 observer(controller)，渲染 controller.results 即可，逻辑完全从 UI 中剥离。

要求：

- module 的查找可以中断(cancelToken 机制)，可以分片
- 每个模块在注册时获得 reSearch，当本模块数据变化时主动重新搜索
- SearchResult 支持泛型
- 数据不要依赖于状态管理框架如 mobx 等，逻辑要纯净
- 查找结果支持排序，由模块的 priority 定义
- module 支持查找、滚动、高亮(匹配项、命中项)、替换、全部替换，可同步可异步
- 仅在 View 层防抖
- 严格分析，防止并发问题

---

- 纯净无依赖：仅使用原生 TS，通过 EventEmitter 实现响应式，易于移植和测试。
- 泛型 Payload：ISearchResult<T> 使得 Controller 可以管理异构结果，而 Plugin 内部又能享受强类型检查。
- 反向控制 (IoC)：通过 reSearch 闭包，插件可以在不知道 Controller 存在的情况下触发更新，完美解耦。
- 并发安全：SearchTask + CancellationToken 确保了在快速输入或网络延迟下，结果始终与当前关键词匹配。
- 增量更新：SearchTask 内部使用 Map 缓存，当某个插件 reSearch 时，只更新该插件的部分，而不影响其他插件的结果，且保持了整体顺序。

---

优化点：

- React 渲染性能（避免无意义的重订阅）。
- 系统负载（减少高频事件下的对象创建压力）。
- 内存安全（防止异步回调在组件卸载后操作已销毁的对象）。=> dispose 模式需要注意这个问题，React 里也有经典报错 `Can't perform a React state update on an unmounted component`。
- 并发控制：双 Token 机制完美解决了全局/局部竞态。
- 性能：一次性内存分配、View 层订阅优化、空删除优化。
- 动态性：支持动态插件注册和动态优先级。
- 健壮性：ID 碰撞防护、生命周期检查、错误隔离。
- API 设计：清晰的同步/异步边界。
- 逻辑：无懈可击（并发、动态性、生命周期均已覆盖）。
- 性能：关键路径（内存分配、渲染防抖）已优化。
- 健壮性：边界情况（空删除、ID 冲突）已处理。

---

## 总结

通过对这套查找替换（Search & Replace）模块的深度分析与重构，我们可以提炼出一套通用的**复杂异步任务管理与插件化架构**的设计思路。这套思路不仅适用于查找替换，也适用于文件扫描、代码诊断（Linting）、全局搜索等场景。

以下是核心设计思路与关键关注点的总结：

### 一、 核心设计思路 (Design Philosophy)

#### 1. 关注点分离 (Separation of Concerns)

- **Model (Task)**: 负责**生产数据**。只关注“如何高效、准确地执行搜索”，不关心数据如何展示，也不关心 UI 的渲染频率。
- **Model (Controller)**: 负责**协调与状态管理**。它是 Task 与 View 的桥梁，管理插件生命周期、处理用户意图（Next/Prev/Replace），并维护全局状态（Current Index, Keyword）。
- **View**: 负责**消费数据与交互**。负责防抖（Debounce）、渲染列表、处理用户输入。
- **Plugin**: 负责**具体实现**。每个插件只关心自己的领域（如编辑器、终端、文件树），实现统一的接口。

#### 2. 任务即对象 (Task as Object)

- 将一次完整的操作（如“搜索 'abc'”）封装为一个独立的 `SearchTask` 对象。
- **优势**：
  - **生命周期管理**：新任务开始 = 旧任务对象销毁。天然解决了“如何取消上一次搜索”的问题。
  - **状态隔离**：每个任务持有自己的 `CancellationToken` 和结果缓存，互不干扰。

#### 3. 双层并发控制 (Dual-Layer Concurrency Control)

- **全局层 (Global)**: 用户输入新关键词 -> 取消整个 `SearchTask`。
- **局部层 (Local)**: 某个插件数据变动 -> 仅取消该插件在当前 Task 中的子任务 (`searchSinglePlugin`)，不影响其他插件。
- **实现**：`MainToken` + `Map<PluginId, PluginToken>`。

#### 4. 插件化与动态性 (Plugin Architecture)

- **统一接口**：定义清晰的 `ISearchPlugin` 接口（Search, Render, ScrollTo）。
- **动态注册**：支持运行时 `register/unregister`，适应懒加载场景。
- **动态优先级**：通过 `priority()` 函数支持运行时排序，适应拖拽排序等需求。

---

### 二、 关键关注点 (Critical Concerns)

在设计此类模块时，必须时刻关注以下 5 个维度：

#### 1. 竞态处理 (Race Conditions)

这是异步系统中最容易出 Bug 的地方。

- **搜索竞态**：先发的搜索请求后返回，覆盖了新请求的结果。 -> **解法**：`CancellationToken` + 版本号/对象引用检查。
- **渲染竞态**：快速点击 Next，导致多次异步 `scrollTo` 冲突。 -> **解法**：DOM 覆盖特性（轻量）或 渲染锁/版本号（严格）。
- **生命周期竞态**：异步回调回来时，组件已卸载。 -> **解法**：`isDisposed` 检查。

#### 2. 性能优化 (Performance)

- **内存分配**：在大数据量聚合时，避免 `push(...arr)`，使用预计算长度 + `new Array(len)` 一次性分配。
- **渲染频率**：
  - **Task 层**：只在数据实质变化（非空）时触发聚合。
  - **View 层**：分离订阅逻辑，避免无意义的重渲染；使用防抖处理高频输入。
- **插件侧**：插件内部也应有防抖，避免高频事件（如打字）频繁触发 Controller 的 `reSearch`。

#### 3. 健壮性 (Robustness)

- **错误隔离**：插件是不可信的。Controller 调用插件方法（Search, Render）时必须包裹 `try-catch`，防止一个插件崩溃搞挂整个应用。
- **ID 冲突**：不同插件可能生成相同的 Result ID。 -> **解法**：使用 `PluginId + ResultId` 复合主键。
- **空值处理**：搜索结果为空、插件列表为空、卸载不存在的插件等边界情况。

#### 4. API 设计 (API Design)

- **同步 vs 异步**：明确哪些操作必须异步（Search），哪些允许异步但通常不等待（Render/ScrollTo）。
- **职责边界**：
  - `render`: 负责状态展示（高亮），并确保元素处于“可被滚动”的状态（如展开折叠）。
  - `scrollTo`: 仅负责视口移动。
- **不可变性**：对外暴露的状态（State）应该是 Immutable 的，防止外部意外修改内部数据。

#### 5. 可测试性与调试 (Testability)

- **纯净逻辑**：核心算法（Task, Controller）不应依赖 UI 框架（React/Vue）或全局单例，便于编写单元测试。
- **状态可见**：通过 `EventEmitter` 暴露状态，方便外部监控和调试。

### 三、 总结代码结构

最终形成的架构图谱如下：

```mermaid
graph TD
    User[User Input] --> View[View Layer (React)]
    View -- Debounce --> Controller[SearchController]

    subgraph Model Layer
        Controller -- Create --> Task[SearchTask]
        Controller -- Register --> Plugins[Plugin Map]

        Task -- 1. Global Cancel --> Token[MainToken]
        Task -- 2. Local Cancel --> TokenMap[PluginTokenMap]

        Task -- Async Search --> PluginA[Plugin A]
        Task -- Async Search --> PluginB[Plugin B]

        PluginA -- Result --> Task
        PluginB -- Result --> Task

        Task -- Aggregate (Sort & Flat) --> State[SearchState]
    end

    State -- Emit --> View
    Controller -- Render/Scroll --> PluginA
    Controller -- Render/Scroll --> PluginB
```

这套设计模式是处理前端复杂异步业务逻辑的标准范式。

---

## 对比 MVC 架构

- Model (SearchStore):

职责: 仅存储数据 (keyword, results 等) 和通知变更。
特点: 它是“哑”的，不知道搜索是怎么发生的，也不知道视图长什么样。
维护性: 如果你想加一个字段（比如 searchTime），只需要改 Store，不会破坏 Controller 的逻辑。
Controller (SearchController):

- 职责: 处理业务逻辑 (search, next, replace)，协调 SearchTask 和 Plugins，最后调用 store.setState。
  特点: 它是“大脑”，负责决策。
  维护性: 复杂的并发控制、插件管理都在这里，与数据存储分离。

- View (SearchPanel):

职责: 渲染 UI，管理 UI 状态 (searchMode, isExpanded)，调用 Controller 方法。
特点: 它是“皮”，只负责展示。
