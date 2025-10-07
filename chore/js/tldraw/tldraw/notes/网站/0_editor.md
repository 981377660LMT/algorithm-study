https://tldraw.dev/docs/editor

好的，我们来详细讲解 `tldraw` 的核心——`Editor` 类。这篇文档是理解如何以编程方式与 `tldraw` 画布交互的关键。

### **`Editor` 类详解**

`Editor` 类是控制 `tldraw` 编辑器的主要入口。你可以用它来管理编辑器的内部状态、对文档内容进行修改，或者响应已经发生的变化。

从设计上讲，`Editor` 的 API 接口非常广泛，几乎所有功能都可以通过它来访问。例如，需要创建图形吗？使用 `Editor.createShapes`。需要删除它们？使用 `Editor.deleteShapes`。需要获取当前页面上所有已排序的图形？使用 `Editor.getCurrentPageShapesSorted`。

本文档将宏观介绍 `Editor` 类的组织结构和一些相关的架构概念。完整的 API 参考可以在 `Editor` API 文档中找到。

---

### **如何使用 `Editor`**

有两种方式可以访问 `editor` 实例：

1.  **通过 `Tldraw` 组件的 `onMount` 回调函数**
    `editor` 实例会作为回调函数的第一个参数传入。

    ```tsx
    function App() {
      return (
        <Tldraw
          onMount={editor => {
            // 在这里编写你的 editor 代码
          }}
        />
      )
    }
    ```

2.  **通过 `useEditor` Hook**
    这个 Hook 必须在 `<Tldraw>` 组件的 JSX 子节点内部调用。

    ```tsx
    function InsideOfContext() {
      const editor = useEditor()
      // 在这里编写你的 editor 代码
      return null // 或者其他任何组件
    }

    function App() {
      return (
        <Tldraw>
          <InsideOfContext />
        </Tldraw>
      )
    }
    ```

---

### **核心概念**

#### **Store (数据仓库)**

编辑器的所有原始状态都保存在 `Editor.store` 属性中。数据以可序列化为 JSON 的记录（records）表的形式存储。

例如，`store` 中为每个页面包含一个 `TLPage` 记录，为每个页面存储编辑器状态的 `TLInstancePageState` 记录，以及为每个编辑器实例存储当前页面 ID 的单个 `TLInstance` 记录。

`Editor` 还暴露了许多从 `store` 中派生出来的**计算值**。例如，`Editor.getSelectedShapeIds` 方法返回当前页面上被选中图形的 ID 数组。你可以直接使用这些属性，也可以在响应式信号（signals）中使用它们。

```tsx
import { track, useEditor } from 'tldraw'

// 使用 track 创建一个响应式组件
export const SelectedShapeIdsCount = track(() => {
  const editor = useEditor()
  // 当选中图形变化时，这个组件会自动重新渲染
  return <div>{editor.getSelectedShapeIds().length}</div>
})
```

#### **改变状态**

`Editor` 类有许多方法可以更新其状态。例如，使用 `Editor.setSelectedShapes` 更改当前选中的图形。还有一些便捷方法，如 `Editor.select`、`Editor.selectAll` 或 `Editor.selectNone`。

```javascript
editor.selectNone() // 取消所有选择
editor.select(myShapeId, myOtherShapeId) // 选中指定的图形
editor.getSelectedShapeIds() // 返回: [myShapeId, myOtherShapeId]
```

每一次状态变更都发生在一个**事务（transaction）** 中。你可以使用 `Editor.batch` 方法将多个变更合并到单个事务中。尽可能地使用批量处理是个好习惯，因为这可以减少持久化或分发这些变更的开销。

#### **撤销与重做 (Undo and Redo)**

`tldraw` 的历史记录栈包含两种数据：

- **"diffs"**: 对 `store` 的变更记录。
- **"marks"**: 撤销/重做的停止点，通过调用 `Editor.markHistoryStoppingPoint` 创建。

当你调用 `Editor.undo` 时，编辑器会撤销每个 diff，直到遇到一个 mark 或栈的起点。调用 `Editor.redo` 则会重做每个 diff，直到遇到一个 mark 或栈的终点。

```javascript
editor.createShapes(...)
// 状态 A
editor.markHistoryStoppingPoint() // 创建一个停止点
editor.selectAll()
editor.duplicateShapes(editor.getSelectedShapeIds())
// 状态 B

editor.undo() // 将会返回到状态 A
editor.redo() // 将会返回到状态 B
```

- `Editor.bail()`: 撤销到最近的 mark 并删除这些 diffs，使其无法被重做。
- `Editor.bailToMark(markId)`: 撤销到指定的 mark。

#### **在上下文中运行代码 (`Editor.run`)**

使用 `Editor.run` 方法可以在一个事务中运行一个函数。事务期间的所有变更将一次性提交。这可以提高性能并避免不必要的 UI 重新渲染。

```javascript
editor.run(() => {
  editor.createShapes(myShapes)
  editor.sendToBack(myShapes)
  editor.selectNone()
})
```

你还可以为 `run` 方法提供上下文选项：

- **忽略历史记录**:
  ```javascript
  editor.run(
    () => {
      editor.createShapes(myShapes)
    },
    { history: 'ignore' } // 这个操作不会被记录到撤销/重做栈中
  )
  ```
- **操作锁定的图形**:
  ```javascript
  editor.run(
    () => {
      editor.updateShapes(myLockedShapes)
    },
    { ignoreShapeLock: true } // 忽略图形锁定，强制更新
  )
  ```

#### **事件 (Events) 和状态图 (State Chart)**

`Editor` 通过其 `dispatch` 方法接收事件。事件首先在内部处理以更新 `Editor.inputs` 等状态，然后被发送到编辑器的**状态图**中。

状态图是一个由 `StateNode` 实例组成的树，它包含了编辑器工具（如选择工具、绘制工具）的逻辑。用户的交互（如移动光标）会根据当前激活的节点产生不同的状态变更。

- **路径 (Path)**: 你可以通过 `editor.root.path` 获取当前激活的状态路径，例如 `"root.select.idle"`。
- `Editor.isIn(path)`: 检查某个路径是否处于激活状态。
- `Editor.getCurrentToolId()`: 获取当前工具的 ID，例如 `'select'` 或 `'draw'`。

#### **副作用 (Side Effects)**

`Editor.sideEffects` 对象允许你为 `Store` 中记录的生命周期关键部分注册回调。你可以在记录被创建、更改或删除之前或之后注册回调。这对于应用约束、维护关系或检查数据完整性非常有用。

```javascript
// 示例：每次创建箭头图形时打印日志
editor.sideEffects.registerAfterCreateHandler('shape', newShape => {
  if (newShape.type === 'arrow') {
    console.log('一个新的箭头图形被创建了', newShape)
  }
})
```

#### **输入 (Inputs)**

`Editor.inputs` 对象保存了用户的当前输入状态，包括光标位置（页面坐标和屏幕坐标）、按下的键、多点点击状态等。**注意：此属性不是响应式的。**

#### **相机和坐标 (Camera and Coordinates)**

- **视口 (Viewport)**: 编辑器组件所包含的矩形区域。
- **屏幕坐标 (Screen Coordinates)** vs. **页面坐标 (Page Coordinates)**:
  - **屏幕坐标**: 相对于组件左上角的像素距离。
  - **页面坐标**: 相对于无限画布原点 `(0,0)` 的距离。
- **坐标转换**:
  - `Editor.screenToPage(point)`: 屏幕坐标转页面坐标。
  - `Editor.pageToScreen(point)`: 页面坐标转屏幕坐标。
- **控制相机**:
  - `Editor.setCamera({ x, y, z })`: 移动相机到指定位置和缩放级别。
  - `Editor.zoomIn() / zoomOut()`: 缩放。
  - `Editor.zoomToFit()`: 缩放以适应所有内容。
  - `Editor.zoomToBounds(box)`: 缩放以适应指定的边界框。
  - `Editor.resetZoom()`: 重置缩放。

#### **绑定 (Bindings)**

绑定是图形之间的一种关系，用于连接图形，使它们可以一起更新。例如，`tldraw` 的箭头通过绑定将其端点连接到其他图形上。你可以定义自己的绑定类型来实现自定义行为。

#### **图像导出 (Image Exports)**

你可以使用 `Editor.toImage` 将画布内容导出为图像（PNG, JPG, WEBP）。

- `Editor.getSvgElement()` 或 `Editor.getSvgString()`: 如果你想直接处理 SVG。
- **自定义图形导出**: 默认情况下，自定义图形使用其 `component` 方法渲染并嵌入到 `<foreignObject>` 中。为了更好的兼容性，你可以实现 `ShapeUtil.toSvg` 方法来提供一个纯粹的 SVG 表示。

---

### **常见 `Editor` 操作示例**

- **创建图形**:

  ```javascript
  editor.createShapes([
    {
      type: 'geo',
      x: 100,
      y: 100,
      props: { geo: 'rectangle', w: 100, h: 100, color: 'blue' }
    }
  ])
  ```

- **更新图形**:

  ```javascript
  editor.updateShapes([
    {
      id: shape.id, // 必须
      type: shape.type, // 必须
      props: { w: 200 }
    }
  ])
  ```

- **删除图形**:

  ```javascript
  editor.deleteShapes([shape.id])
  ```

- **开启只读模式**:

  ```javascript
  editor.updateInstanceState({ isReadonly: true })
  ```

- **开启暗黑模式**:

  ```javascript
  editor.user.updateUserPreferences({ colorScheme: 'dark' })
  ```

- **隐藏图形**:
  在 `<Tldraw>` 组件上提供 `getShapeVisibility` 函数。
  ```tsx
  <Tldraw getShapeVisibility={shape => (shape.meta.hidden ? 'hidden' : 'inherit')} />
  ```
  这只会隐藏图形的渲染，图形数据仍然存在于 `store` 中。
