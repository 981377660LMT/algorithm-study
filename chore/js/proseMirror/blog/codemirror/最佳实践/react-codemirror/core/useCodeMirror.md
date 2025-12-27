这段 `useCodeMirror.ts` 是整个 `react-codemirror` 组件的**灵魂**。它不仅负责初始化编辑器，还处理了 React 开发中最头疼的**受控组件同步、输入防抖、以及性能优化**。

以下是该 Hook 的深度细节讲解：

---

### 1. 核心矛盾：React 的“声明式”与 CodeMirror 的“命令式”

React 期望通过 `value` 属性控制一切，而 CodeMirror 是一个持久的实例，通过 `Transaction`（事务）改变状态。如果直接在 `useEffect` 里监听 `value` 变了就 `dispatch`，会导致：

- **光标跳动**：用户打字时，React 状态更新慢于编辑器，导致旧值覆盖新值，光标回到行首。
- **死循环**：编辑器改变 -> 触发 `onChange` -> React 更新 `value` -> `useEffect` 发现 `value` 变了 -> 再次更新编辑器。

---

### 2. 绝招一：Typing Latch（输入锁存器）

这是代码中最精妙的设计，用于解决**受控组件的延迟冲突**。

- **`TYPING_TIMEOUT = 200ms`**：定义了用户“正在打字”的判定阈值。
- **`typingLatch`**：一个自定义的计时器（`TimeoutLatch`）。
- **逻辑流程**：
  1.  **用户打字**：触发 `updateListener`。
  2.  **启动/重置计时器**：`typingLatch.current.reset()`。此时 `isTyping` 为 `true`。
  3.  **外部 Value 变化**：在最后一个 `useEffect` 中，如果发现 `value` 变了，但 `isTyping` 为真，**不立即更新编辑器**，而是把更新逻辑存入 `pendingUpdate.current`。
  4.  **停笔冲刷**：当用户停止打字 200ms 后，计时器触发回调，执行 `pendingUpdate.current()`，将外部最新的值同步进来。

**细节：** 这种机制保证了用户在连续输入时，`编辑器绝对不会被外部传入的旧数据“闪回”`，同时保证了最终数据的一致性。

---

### 3. 绝招二：`ExternalChange` 抑制反馈环

为了彻底切断死循环，代码定义了一个注解：

```typescript
export const ExternalChange = Annotation.define<boolean>()
```

- **发送端**：当 Hook 手动同步外部 `value` 到编辑器时，会打上这个标记：`annotations: [ExternalChange.of(true)]`。
- **接收端**：在 `updateListener` 中检查：
  ```typescript
  !vu.transactions.some(tr => tr.annotation(ExternalChange))
  ```
  如果发现这次变更带有 `ExternalChange` 标记，说明是“自产自销”的更新，**不再调用 `onChange`**。

---

### 4. 动态重构：`StateEffect.reconfigure`

CodeMirror 6 的扩展树是不可变的。如果你想在不销毁编辑器的情况下更改插件（比如切换主题、改变只读状态），必须使用 `reconfigure`。

```typescript
useEffect(() => {
  if (view) {
    view.dispatch({ effects: StateEffect.reconfigure.of(getExtensions) })
  }
}, [theme, extensions, editable, ...])
```

**细节：** 这里的依赖数组非常长。这意味着当任何配置项改变时，Hook 都会生成一套新的扩展数组并派发给编辑器。CM6 内部会高效地对比新旧扩展，只更新变化的部分（例如只切换 CSS 类名而不重新解析语法树）。

---

### 5. 初始化细节：`useLayoutEffect`

代码使用 `useLayoutEffect` 来创建 `EditorView`：

- **原因**：`EditorView` 的创建涉及 DOM 挂载和尺寸计算。使用 `useLayoutEffect` 可以确保在浏览器重绘之前完成初始化，避免用户看到“先白屏后出现编辑器”的闪烁。
- **`initialState` 支持**：
  ```typescript
  initialState ? EditorState.fromJSON(...) : EditorState.create(...)
  ```
  这允许从序列化的 JSON 中恢复编辑器，包括撤销历史、插件状态等，非常适合“恢复上次编辑进度”的功能。

---

### 6. 样式注入：`defaultThemeOption`

代码将 React 的 Props（`height`, `width` 等）转换成了 CodeMirror 的内部主题：

```typescript
const defaultThemeOption = EditorView.theme({
  '&': { height, minHeight, ... },
  '& .cm-scroller': { height: '100% !important' }
})
```

**细节：** CM6 的高度控制比较特殊，通常需要作用于 `.cm-scroller`。这里通过注入 `!important` 确保了 React 传入的高度属性具有最高优先级，解决了常见的 CSS 布局失效问题。

---

### 7. 内存管理：清理函数

```typescript
useEffect(
  () => () => {
    if (view) {
      view.destroy() // 销毁编辑器实例
    }
    if (typingLatch.current) {
      typingLatch.current.cancel() // 取消计时器
    }
  },
  [view]
)
```

**细节：** 必须显式调用 `view.destroy()`。否则，CodeMirror 挂载在 `window` 上的全局监听器（如 `resize`）和后台解析线程将永远运行，导致严重的内存泄漏。

### 总结

这个 Hook 的设计体现了对 CodeMirror 6 架构的深度理解：

1.  **异步性**：利用 `TimeoutLatch` 处理 React 状态延迟。
2.  **原子性**：利用 `Annotation` 标记事务来源。
3.  **高效性**：利用 `reconfigure` 实现无损配置更新。
4.  **安全性**：严格的生命周期管理和类型检查。
