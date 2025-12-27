这段 `core/index.tsx` 文件是整个 `react-codemirror` 组件的**入口核心**。它扮演了“桥梁”的角色：将 React 的声明式（Declarative）编程模型转换为 CodeMirror 6 的命令式（Imperative）配置。

以下是对该文件核心设计的深度拆解：

---

### 1. 类型定义与属性设计 (`ReactCodeMirrorProps`)

该接口定义了 React 组件如何与 CodeMirror 交互。

- **继承与排除**：
  - 继承了 `EditorStateConfig` 但排除了 `doc` 和 `extensions`。这是因为 React 中习惯使用 `value` 而非 `doc`，且 `extensions` 需要在内部进行特殊处理（如 `Compartment` 动态重构）。
  - 继承了 `React.HTMLAttributes<HTMLDivElement>`，允许用户像操作普通 `div` 一样传递 `style`、`className` 或 `onClick`。
- **核心配置项**：
  - **`basicSetup`**: 允许用户一键开启行号、括号匹配等基础功能，或者通过 `BasicSetupOptions` 进行微调。
  - **`theme`**: 支持字符串（'light'/'dark'）或直接传入 CM6 的 `Extension` 对象，体现了极高的灵活性。
  - **`initialState`**: 这是一个高级功能，允许从序列化的 JSON 中恢复编辑器状态（包括撤销历史和插件状态）。

---

### 2. 引用暴露机制 (`ReactCodeMirrorRef`)

由于 CodeMirror 是命令式的，开发者经常需要直接操作 `EditorView`（例如手动滚动、派发事务）。

- **`useImperativeHandle`**:
  ```typescript
  useImperativeHandle(ref, () => ({ editor: editor.current, state: state, view: view }), [
    editor,
    container,
    state,
    view
  ])
  ```
  - 通过 `forwardRef`，父组件可以拿到 `view`（视图实例）和 `state`（状态实例）。
  - 这使得 React 组件能够调用 `ref.current.view.dispatch(...)` 来执行非受控的操作。

---

### 3. 容器挂载逻辑 (`setEditorRef`)

这是组件初始化的关键点。

- **Callback Ref 模式**：
  ```typescript
  const setEditorRef = useCallback(
    (el: HTMLDivElement) => {
      editor.current = el
      setContainer(el)
    },
    [setContainer]
  )
  ```
  - 没有使用传统的 `useEffect` + `useRef` 绑定，而是使用了 **Callback Ref**。
  - **优点**：当 `div` 真正挂载到 DOM 时，立即调用 `setContainer`。这会触发 `useCodeMirror` 内部的初始化逻辑。这种方式比 `useEffect` 更能保证在容器可用时第一时间创建编辑器。

---

### 4. 状态驱动中心 (`useCodeMirror`)

虽然 `index.tsx` 是入口，但它把所有的重活都外包给了 `useCodeMirror` 这个自定义 Hook。

- **参数透传**：它将所有的 Props（`value`, `theme`, `extensions` 等）传给 Hook。
- **解构返回**：从 Hook 中拿到 `state` 和 `view`。这两个值是响应式的，当编辑器内部发生变化时，它们会更新，从而触发 `useImperativeHandle` 的更新。

---

### 5. 样式与主题处理

- **动态类名**：
  ```typescript
  const defaultClassNames = typeof theme === 'string' ? `cm-theme-${theme}` : 'cm-theme'
  ```
  - 如果 `theme` 是字符串，它会添加一个类名（如 `cm-theme-dark`），方便外部通过 CSS 进行全局样式覆盖。
- **尺寸控制**：
  - Props 中的 `height`, `width`, `minHeight` 等被传给了 Hook。Hook 内部会将这些值转换成 `EditorView.theme` 扩展，直接作用于编辑器的内部滚动容器。

---

### 6. 健壮性检查

- **类型校验**：
  ```typescript
  if (typeof value !== 'string') {
    throw new Error(`value must be typeof string but got ${typeof value}`)
  }
  ```
  - 在渲染前强制检查 `value` 类型。因为 CM6 的 `Text` 模型只接受字符串，传入 `undefined` 或 `null` 会导致底层崩溃。

---

### 总结：核心流程图

1.  **Render**: React 渲染 `div` 容器。
2.  **Mount**: `setEditorRef` 被调用，获取 DOM 节点。
3.  **Hook Init**: `useCodeMirror` 接收到 DOM 节点，创建 `EditorState` 和 `EditorView`。
4.  **Callback**: `onCreateEditor` 被触发，通知父组件。
5.  **Update**: 当 Props（如 `value`）改变，Hook 内部通过 `view.dispatch` 更新内容，而不是重新创建编辑器。
6.  **Unmount**: Hook 内部的 `useEffect` 清理函数调用 `view.destroy()`。

这个 `core/index.tsx` 的设计非常精简，它只负责**声明属性**和**管理引用**，将复杂的生命周期和事务逻辑封装在 Hook 中，是 React 集成复杂命令式库的标准范式。
