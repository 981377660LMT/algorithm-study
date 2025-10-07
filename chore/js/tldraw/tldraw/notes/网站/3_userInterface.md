好的，我们来详细讲解 `tldraw` 中关于 **User Interface (用户界面)** 的定制。这部分文档解释了如何控制和修改 `tldraw` 的菜单、工具栏、快捷键等视觉和交互元素。

---

### **`tldraw` 的用户界面 (User Interface)**

在 `tldraw` 中，用户界面（User Interface）涵盖了你在编辑器中看到和交互的所有元素，包括菜单、工具栏、键盘快捷键以及 UI 事件。

---

### **1. 隐藏整个 UI (`hideUi`)**

如果你想完全移除 `tldraw` 默认的所有 UI 元素，只保留一个纯净的画布，可以使用 `hideUi` prop。

```tsx
function Example() {
  // 这会隐藏工具栏、菜单、快捷键等所有默认 UI
  return <Tldraw hideUi />
}
```

**这个功能有什么用？**

这个功能非常适合当你想要完全替换 `tldraw` 的默认界面，构建自己独特的 UI 风格时使用。

**核心原理**：`tldraw` 的所有默认 UI 都是通过调用核心的 `Editor` API 来工作的。即使你隐藏了 UI，底层的 `editor` 实例依然功能完备。你可以通过 `onMount` 回调或 `useEditor` Hook 获取 `editor` 实例，然后用自己的按钮、菜单去调用 `editor` 的方法（如 `editor.createShapes`, `editor.undo` 等），从而实现自定义的交互。

例如，即使 UI 被隐藏，你仍然可以在浏览器控制台中通过 `editor` 实例切换工具：

```javascript
// 在浏览器控制台中执行
editor.setCurrentTool('draw')
```

之后你就可以在画布上进行绘制了。

---

### **2. 监听 UI 事件 (`onUiEvent`)**

如果你想知道用户在默认 UI 上执行了哪些操作（例如，点击了哪个菜单项），可以使用 `onUiEvent` 回调。

```tsx
function Example() {
  function handleEvent(eventName, data) {
    // 在这里处理 UI 事件
    console.log('UI Event:', eventName, data)
    // 例如: UI Event: 'align-shapes', { source: 'menu', direction: 'left' }
  }

  return <Tldraw onUiEvent={handleEvent} />
}
```

- `eventName`: 事件名称的字符串，如 `'align-shapes'`。
- `data`: 一个包含事件详细信息的对象，通常包括事件来源 `source`（如 `'menu'` 或 `'context-menu'`）以及其他事件特定的数据。

**请注意**：`onUiEvent` **仅在用户与默认 UI 交互时**被调用。如果你通过代码直接调用 `editor.alignShapes()`，这个回调**不会**被触发。

---

### **3. 自定义 UI 内容 (`overrides`)**

`overrides` prop 是定制 `tldraw` UI 的最强大工具。它允许你重写（修改、删除或添加）菜单项、工具栏按钮和快捷键等。

`overrides` 接收一个 `TLUiOverrides` 对象，该对象包含一系列方法，分别对应 UI 的不同部分。

#### **a. 自定义操作 (Actions)**

“操作”是在菜单和快捷键中使用的共享命令。你可以通过 `overrides.actions` 方法来修改它们。

这个方法接收 `editor` 实例和默认的 `actions` 对象，你必须返回一个修改后的 `actions` 对象。

```tsx
const myOverrides: TLUiOverrides = {
  actions(editor, actions) {
    // 示例 1: 删除一个操作 (记得也要删除引用它的菜单项)
    delete actions['insert-embed']

    // 示例 2: 创建一个新操作或覆盖一个现有操作
    actions['my-new-action'] = {
      id: 'my-new-action',
      label: 'My new action', // 显示在菜单中的标签
      readonlyOk: true, // 在只读模式下是否可用
      kbd: 'cmd+u,ctrl+u', // 绑定的快捷键 (macOS, Windows/Linux)
      onSelect(source: any) {
        // 当操作被触发时执行的逻辑
        window.alert('我的新操作被触发了!')
      }
    }
    return actions
  }
  // ... 可能还有其他重写，比如在菜单中添加这个 action
}
```

#### **b. 自定义工具 (Tools)**

与 `actions` 类似，你可以通过 `overrides.tools` 方法来定义工具栏上按钮的行为和外观。这通常用于为你创建的**自定义工具**添加 UI 入口。

```tsx
const myOverrides: TLUiOverrides = {
  tools(editor, tools) {
    // 为我们之前创建的 'card' 自定义工具添加一个工具栏按钮
    tools.card = {
      id: 'card',
      icon: 'color', // 使用一个内置图标 (也可以是自定义图标组件)
      label: 'tools.card', // 标签，用于国际化翻译
      kbd: 'c', // 快捷键
      onSelect: () => {
        // 当点击工具栏按钮时，切换到 'card' 工具
        editor.setCurrentTool('card')
      }
    }
    return tools
  }
}
```

#### **c. 自定义翻译 (Translations)**

如果你在 `actions` 或 `tools` 中使用了自定义的 `label`（如 `'tools.card'`)，你需要提供对应的翻译文本。

```tsx
const myOverrides: TLUiOverrides = {
  // ... 其他重写
  translations: {
    // 至少提供英文版本
    en: {
      'tools.card': 'Card'
    },
    // 你也可以提供其他语言的翻译
    'zh-cn': {
      'tools.card': '卡片'
    }
  }
}
```

通过组合使用 `hideUi` 和 `overrides`，你可以实现从微调到完全重构 `tldraw` 用户界面的任何需求。
