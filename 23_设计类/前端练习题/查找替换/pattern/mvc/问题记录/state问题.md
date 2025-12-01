在 MVC 模式中，`SearchStore`（Model 层）的状态边界定义直接决定了代码的**可维护性**和**职责清晰度**。

关于“`state` 需不需要当前的查询字符串（keyword）”，我的结论是：**需要，但取决于你的架构目标**。

### 1. 结论：为什么建议放入 State

在大多数现代前端应用（尤其是像 VS Code 搜索面板这种复杂交互）中，将 `keyword` 放入 `SearchStore` 是更好的选择。

#### ✅ 理由一：单一数据源 (Single Source of Truth)

如果 `keyword` 不在 Store 里，它通常会散落在：

1.  View 层的 React `useState`（输入框的值）。
2.  Controller 层的某个私有变量（用于发起搜索）。

这就导致了**状态同步问题**：

- 如果用户点击了“清除”按钮，View 层清空了输入框，但 Controller 层可能还记着旧的 keyword。
- 如果 Controller 想要恢复上一次的搜索（比如切换 Tab 回来），它需要通知 View 层把字填回去。

将 `keyword` 放入 Store，View 层直接订阅 Store，Controller 直接读取 Store，大家看的是同一份数据。

#### ✅ 理由二：派生状态与联动

搜索结果 `results` 是 `keyword` 的**派生数据**。
`results` 只有在特定的 `keyword` 下才有意义。如果 Store 里只有 `results` 而没有 `keyword`，那么 Store 的状态是不完整的——你不知道这一堆结果是搜什么得出来的。

#### ✅ 理由三：方便实现“历史记录”和“状态恢复”

如果你想做“搜索历史”功能，或者用户关闭面板再打开时保持原样，你只需要序列化/反序列化整个 `SearchStore.state` 即可。如果 `keyword` 在 View 层，这就很难做。

---

### 2. 详细讲解：State 的边界划分原则

在设计 Model 层时，我们通常把状态分为三类。你需要根据这三类来决定哪些进 Store。

#### 第一类：领域数据 (Domain Data) —— **必须进 Store**

这是业务的核心数据，不依赖于 UI 如何展示。

- **例子**：`results` (搜索结果列表), `matchCount` (匹配总数)。
- **判定标准**：如果我换一个 UI（比如从 React 换成 Vue，或者换成命令行界面），这些数据还需要吗？如果需要，就进 Store。

#### 第二类：应用状态 (App State) —— **建议进 Store**

这是描述应用当前“行为”的数据，通常涉及多个组件共享。

- **例子**：`keyword` (查询词), `isSearching` (加载中), `currentIndex` (当前高亮项), `searchOptions` (大小写敏感/正则模式)。
- **判定标准**：
  - 这个状态是否影响了“领域数据”的生成？（`keyword` 影响 `results` -> 进 Store）
  - Controller 逻辑是否依赖它？（Controller 需要知道 `isRegex` 才能发起正则搜索 -> 进 Store）

#### 第三类：UI 状态 (Ephemeral UI State) —— **不要进 Store**

这是纯粹为了视觉展示的临时状态，只属于某个具体组件。

- **例子**：输入框是否聚焦 (`isFocused`)、某个折叠面板是否展开 (`isExpanded`)、滚动条位置 (`scrollTop`)。
- **判定标准**：如果页面刷新后这个状态丢了，用户会生气吗？如果不生气，或者这个状态只影响像素级的渲染而不影响业务逻辑，就留在 View 组件内部 (`useState`)。

---

### 3. 修正后的代码建议

基于上述分析，你的 `SearchState` 应该包含 `keyword`。

```typescript
export interface SearchState {
  // 1. 输入 (App State)
  keyword: string
  options: {
    isCaseSensitive: boolean
    isRegex: boolean
  }

  // 2. 输出 (Domain Data)
  results: ISearchResult[]

  // 3. 交互状态 (App State)
  currentIndex: number // 当前选中的结果索引
  isSearching: boolean // Loading 状态
}
```

```typescript
// ...existing code...
const INITIAL_STATE: SearchState = {
  keyword: '',
  options: { isCaseSensitive: false, isRegex: false },
  results: [],
  currentIndex: -1,
  isSearching: false
}
// ...existing code...
```

### 4. 边界划分的“灰色地带”

有些状态比较暧昧，比如 **“输入框的临时值”** vs **“确认搜索的关键词”**。

- **场景**：用户想搜 "Hello World"。
- 用户输入 "H" -> "He" -> "Hel"...
- 此时 Store 里的 `keyword` 应该跟着变吗？

**方案 A：实时同步 (Real-time)**
用户每打一个字，Store 更新，触发搜索。

- **优点**：即时响应，体验好。
- **缺点**：性能压力大（需要防抖）。
- **适用**：本地快速搜索。

**方案 B：提交时同步 (On Commit)**
View 层用 `useState` 存临时输入，用户按回车后，才调用 `store.setKeyword()`。

- **优点**：性能好，逻辑简单。
- **缺点**：UI 和 Store 在用户输入期间不一致。
- **适用**：全局耗时搜索。

**最佳实践**：
通常 Store 存的是 **“生效的查询词” (Effective Query)**。
View 层可以维护一个临时的 `inputValue`，只有在防抖结束或回车后，才同步给 Store 的 `keyword`。这样 Store 里的 `keyword` 永远代表“当前展示结果所对应的查询词”，保持了 `keyword` 和 `results` 的强一致性。
