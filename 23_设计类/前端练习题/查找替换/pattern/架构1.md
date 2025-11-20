这段代码实现了一个**可扩展的、基于插件化架构的全局搜索替换系统**（类似于 VS Code 的搜索功能）。目前的 `modal.ts` 是一个包含所有逻辑的“巨石”原型。

根据你的 TODO (`Orchestrator 太重了`) 和提供的文件结构 (`search-replace/`)，我将首先深入讲解现有代码的机制，然后给出一个分层解耦的架构设计方案。

---

### 第一部分：现有代码 (`modal.ts`) 深度解析

当前代码采用了 **观察者模式 (Observer)** 和 **策略模式 (Strategy)** 的混合体。

#### 1. 核心角色

1.  **`SearchProvider` (策略/插件)**

    - **职责**：它是“具体怎么搜”和“具体怎么高亮”的实现者。
    - **解耦**：编排器不需要知道是搜表格还是搜评论，只管调用 `search` 接口。
    - **能力**：
      - `search`: 返回标准化的 `SearchResult`。
      - `highlight/focus`: 控制具体的 UI（比如滚动表格、展开侧边栏）。
      - `replace`: 处理数据更新。
      - `onRegister`: **关键点**，它接收 `ProviderContext`，从而获得了反向通知编排器的能力（`requestSearch`）。

2.  **`SearchOrchestrator` (编排器/中枢)**

    - **职责**：它是“大脑”，管理生命周期、状态和并发。
    - **核心机制**：
      - **并发控制**：使用 `Promise.all` 并行执行所有 Provider 的搜索。
      - **竞态处理 (Cancellation)**：使用 `cancellationToken` (`{ isCancelled: boolean }`)。当新搜索开始时，旧的 token 被标记为取消，Provider 在 `await` 之后检查此标记，防止旧结果覆盖新结果。
      - **防抖 (Debounce)**：`triggerDebouncedSearch` 防止用户打字或数据频繁变化导致的高频搜索。
      - **焦点保持 (Focus Preservation)**：在 `refreshSearch` 中，数据更新导致重搜后，它会尝试通过 ID 找回之前选中的那个结果，避免用户丢失上下文。

3.  **数据流向**

    - **搜索流 (User Action)**:
      用户输入 -> Orchestrator.search -> Provider.search -> 聚合结果 -> Provider.highlight
    - **更新流 (Data Change)**:
      外部数据变动 -> Provider.updateData -> 调用 `context.requestSearch` -> Orchestrator 防抖 -> Orchestrator.refreshSearch -> 局部/全量更新 UI

#### 2. 现有代码的痛点 (为什么要重构)

正如 TODO 所述，`SearchOrchestrator` 太重了，它承担了过多的职责：

1.  **状态管理**：持有 `currentResults`, `focusIndex`。
2.  **业务逻辑**：防抖、取消逻辑、焦点计算。
3.  **UI 协调**：直接指挥 Provider 进行 highlight。
4.  **生命周期**：管理 Provider 的注册销毁。

---

### 第二部分：架构设计与重构方案

为了解决“重”的问题，我们需要将**数据(Entity)**、**控制(Controller)**、**视图(View)** 和 **执行单元(Task)** 分离。

结合你提供的 `search-replace` 目录结构，建议采用以下架构：

#### 1. 架构分层图

```mermaid
graph TD
    subgraph View Layer [UI 视图层]
        Modal[SearchReplaceModal] -->|用户输入/按键| Controller
        Modal -->|订阅状态| Store[状态仓库/Reactive State]
    end

    subgraph Controller Layer [控制层]
        Controller[SearchController] -->|创建| Task[SearchTask]
        Controller -->|管理| Registry[ModuleRegistry]
        Controller -->|更新| Store
    end

    subgraph Entity Layer [领域实体层]
        Task -->|执行搜索| Module[SearchModule (Provider)]
        Task -->|持有| Results[SearchResults]
        Task -->|管理| CancellationToken
    end

    subgraph Infrastructure [基础设施层]
        Module -->|访问| GridData[表格数据]
        Module -->|访问| CommentData[评论数据]
    end
```

#### 2. 详细模块设计

根据你的文件结构，重构建议如下：

##### A. `entity/SearchTask/index.ts` (抽离搜索执行)

**目的**：将“一次搜索行为”封装成一个对象。Orchestrator 不再维护 `isCancelled`，而是由 Task 实例自己管理。

```typescript
// 每一个 Task 代表一次具体的搜索请求
export class SearchTask {
  private isCancelled = false
  public results: ISearchResult[] = []

  constructor(
    private modules: ISearchModule[],
    private query: string,
    private options: SearchOptions
  ) {}

  // 执行搜索
  async run(): Promise<ISearchResult[]> {
    const promises = this.modules.map(async module => {
      if (this.isCancelled) return []
      const res = await module.search(this.query, this.options)
      return this.isCancelled ? [] : res.results
    })

    const allResults = await Promise.all(promises)
    if (this.isCancelled) return []

    this.results = allResults.flat()
    return this.results
  }

  cancel() {
    this.isCancelled = true
  }
}
```

##### B. `controller/index.ts` (抽离业务逻辑)

**目的**：作为 UI 和逻辑的胶水层。管理 Task 的创建和销毁，处理防抖。

```typescript
export class SearchController {
  private currentTask: SearchTask | null = null
  private modules: Map<string, ISearchModule> = new Map()

  // 状态可以放在这里，或者使用 MobX/Redux/Vue Ref
  public state = {
    results: [],
    activeIndex: -1,
    isSearching: false
  }

  // 注册模块
  registerModule(module: ISearchModule) {
    /* ... */
  }

  // 响应用户输入
  async onSearch(query: string) {
    // 1. 取消上一个任务
    this.currentTask?.cancel()

    // 2. 创建新任务
    const task = new SearchTask(Array.from(this.modules.values()), query, {})
    this.currentTask = task

    // 3. 执行
    this.state.isSearching = true
    const results = await task.run()

    // 4. 只有当任务没被取消时才更新状态
    if (task === this.currentTask) {
      this.state.results = results
      this.state.isSearching = false
      this.autoFocusFirst()
    }
  }

  // 导航逻辑
  next() {
    /* 修改 activeIndex, 调用 module.scroll/focus */
  }
}
```

##### C. index.d.ts (标准化接口)

你提供的 index.d.ts 已经很好了，建议做以下微调以适应新架构：

- **`ISearchModule`**: 保持不变，这是插件的契约。
- **`ISearchResult`**: 增加 `invalidated` 字段，用于处理替换后的失效状态。
- **`rich<T>`**: 这是一个很好的设计，将原始结果与 Module 绑定，方便 Controller 在点击结果时直接找到对应的 Module 进行渲染或跳转。

##### D. `components/SearchReplaceModal/index.tsx` (纯 UI)

**目的**：只负责渲染。不包含任何搜索逻辑。

- 接收 `results` 和 `activeIndex` 作为 Props（或从 Store 读取）。
- 渲染搜索框和结果列表。
- 将 `onChange`, `onNext`, `onReplace` 事件转发给 `Controller`。

#### 3. 核心改进点总结

1.  **Task 的独立性**：

    - 旧代码：`Orchestrator` 里混杂着 `cancellationToken` 对象，容易出错。
    - 新架构：`SearchTask` 是瞬态的。每次搜索 `new SearchTask()`。取消搜索只需调用 `oldTask.cancel()`，非常清晰。

2.  **状态分离**：

    - 旧代码：`Orchestrator` 既管流程又管状态数据。
    - 新架构：`Controller` 管流程，`Task` 管执行，`Store/State` 管数据。UI 只需要监听数据变化。

3.  **可测试性**：

    - 你可以单独为 `SearchTask` 编写单元测试，验证并发和取消逻辑，而不需要启动整个 UI。

4.  **增量更新 (针对之前的优化)**：
    - 在 `Controller` 中，你可以维护一个 `DirtyModuleSet`。
    - 当 Module 请求更新时，Controller 可以创建一个 `PartialSearchTask`（只包含脏 Module），执行完后将结果 merge 回主状态。

这个架构将原本复杂的 `modal.ts` 拆解为了：**Controller (大脑)** + **Task (执行者)** + **Module (插件)** + **View (皮肤)**，完全符合你文件结构中的规划。
