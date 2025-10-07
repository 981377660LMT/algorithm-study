好的，我们来详细讲解 `tldraw` 中关于 **Tools (工具)** 的概念。这是理解 `tldraw` 交互逻辑的核心。

---

### **什么是工具 (Tool)？**

在 `tldraw` 中，一个**工具 (Tool)** 指的是其内部**状态图 (State Chart)** 中的任何一个顶层状态。

可以把 `tldraw` 编辑器想象成一个大型的**状态机**。用户的每一个动作（点击、拖拽、按下键盘）会触发什么行为，完全取决于当前处于哪个“状态”。而“工具”就是这个状态机最外层的状态。

例如，当用户激活**选择工具 (Select Tool)** 时，点击图形会选中它；而当用户激活**箭头工具 (Arrow Tool)** 时，点击并拖拽则会创建一根箭头。

![tldraw 状态图](https://raw.githubusercontent.com/tldraw/tldraw/main/apps/docs/public/state-chart-tools.png)
_图示：状态图的第一层状态（除 Root 外）就是我们所说的“工具”。_

要了解状态图的更多细节，可以参考 `Editor` 页面的文档。下面我们将重点介绍工具以及如何创建你自己的工具。

---

### **工具的类型**

1.  **核心工具 (Core Tools)**

    - 这些是编辑器内置且始终存在的工具，包括：选择工具 (`select`)、缩放工具 (`zoom`) 和文本工具 (`text`)。

2.  **默认工具 (Default Tools)**

    - 这些是 `<Tldraw>` 组件默认提供的常用工具，例如：手绘工具 (`draw`)、抓手工具 (`hand`)、箭头工具 (`arrow`) 等。

3.  **自定义工具 (Custom Tools)**
    - 你可以创建自己的工具，并通过 `<Tldraw>` 组件的 `tools` prop 将它们添加到状态图中。

**注意**：创建工具的**逻辑**（`StateNode`）与在用户界面（如工具栏）中添加对应的**按钮**是分开的。你需要参考 UI 定制部分来为你的自定义工具添加 UI 入口。

---

### **如何创建自定义工具**

工具的本质是一个继承自 `StateNode` 的类。`StateNode` 定义了状态的行为，包括它的子状态以及如何响应事件。

#### **1. 定义工具的基本结构**

一个工具必须有一个静态的 `id`。如果它有子状态，还必须指定一个 `initial` 初始子状态，并用 `children` 方法返回所有子状态的类。

```typescript
import { StateNode, TLPointerEventInfo } from 'tldraw'

// 定义一个 "空闲" 子状态
class MyIdleState extends StateNode {
  static override id = 'idle'
  // ... 稍后添加事件处理
}

// 定义一个 "点击中" 子状态
class MyPointingState extends StateNode {
  static override id = 'pointing'
  // ... 稍后添加事件处理
}

// 定义我们的主工具
class MyTool extends StateNode {
  // 唯一的 ID
  static override id = 'my-tool'

  // 初始子状态的 ID
  static override initial = 'idle'

  // 定义所有子状态
  static override children() {
    return [MyIdleState, MyPointingState]
  }
}
```

在这个例子中，`MyTool` 工具包含两个子状态：`idle` 和 `pointing`。当 `MyTool` 被激活时，它会自动进入 `idle` 状态。

#### **2. 处理事件**

`StateNode` 类提供了一系列 `on...` 回调方法来响应用户的输入事件，例如 `onPointerDown` (指针按下)、`onPointerMove` (指针移动)、`onKeyDown` (键盘按下) 等。

当编辑器通过 `Editor.dispatch` 方法接收到一个事件时，该事件会从状态图的根节点开始，沿着当前激活的状态链（例如 `root` -> `my-tool` -> `idle`）依次传递。

**事件处理有两个重要规则：**

**规则一：父状态先于子状态处理事件**

```typescript
class MyIdleState extends StateNode {
  static override id = 'idle'
  onPointerDown(info: TLPointerEventInfo) {
    console.log('world') // 第二步执行
  }
}

class MyTool extends StateNode {
  static override id = 'my-tool'
  static override initial = 'idle'
  static override children() {
    return [MyIdleState]
  }

  onPointerDown(info: TLPointerEventInfo) {
    console.log('hello') // 第一步执行
  }
}
```

当 `MyTool` 处于激活状态时，如果发生 `pointer_down` 事件，`MyTool` 的 `onPointerDown` 会先被调用，然后事件才会传递给它的激活子状态 `MyIdleState`。

**规则二：状态转换会中断事件传递**

如果一个父状态在处理事件时触发了状态转换（例如切换到另一个工具或另一个子状态），那么事件传递链会立即停止，后续的子状态将不会收到该事件。

```typescript
class MyIdleState extends StateNode {
  static override id = 'idle'
  onPointerDown(info: TLPointerEventInfo) {
    console.log('这段代码不会执行')
  }
}

class MyTool extends StateNode {
  static override id = 'my-tool'
  static override initial = 'idle'
  static override children() {
    return [MyIdleState]
  }

  onPointerDown(info: TLPointerEventInfo) {
    // 切换到 select 工具，这是一个状态转换
    this.editor.setCurrentTool('select')
  }
}
```

在这个例子中，`MyTool` 的 `onPointerDown` 方法将当前工具切换到了 `select`。这个操作导致了状态转换，因此事件不会再传递给 `MyIdleState`。

---

### **切换工具**

你可以使用 `editor.setCurrentTool` 方法来以编程方式改变当前激活的工具。

```javascript
// 切换到选择工具
editor.setCurrentTool('select')
```

你也可以通过提供一个状态路径来进行“深度转换”，直接进入某个工具的特定子状态。

```javascript
// 直接进入选择工具的橡皮擦功能的 pointing 子状态
editor.setCurrentTool('select.eraser.pointing')
```
