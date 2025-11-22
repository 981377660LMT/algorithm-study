这是一个极具挑战性的任务。要模拟 React 的合成事件机制（Synthetic Event System），我们需要构建一个独立于浏览器的“逻辑世界”。

这个系统必须包含以下核心要素：

1.  **Fiber 树模拟**：React 事件是沿着组件树（Fiber）传播的，而不是 DOM 树。
2.  **事件插件系统 (Event Plugin System)**：React 通过插件来标准化不同类型的事件。
3.  **合成事件对象 (SyntheticEvent)**：包装原生事件，提供跨浏览器兼容性（这里我们模拟其接口）。
4.  **两阶段传播**：捕获（Capture）和冒泡（Bubble）。
5.  **批量执行**：收集所有监听器后再一次性执行。

以下是完整的 TypeScript 实现：

### 1. 基础结构定义 (Fiber & DOM 模拟)

为了模拟 React 的环境，我们需要先定义简化的 Fiber 节点和 DOM 节点。

```typescript
// 模拟原生 DOM 节点（物理世界）
export class NativeNode {
  public parent: NativeNode | null = null
  public id: string

  constructor(id: string) {
    this.id = id
  }

  appendChild(child: NativeNode) {
    child.parent = this
  }
}

// 模拟 React Fiber 节点（逻辑世界）
export interface Fiber {
  tag: string // 组件名或标签名
  stateNode: NativeNode | null // 指向物理 DOM
  return: Fiber | null // 父 Fiber
  child: Fiber | null
  sibling: Fiber | null

  // 存储 props，包含事件处理函数 (onClick, onClickCapture 等)
  memoizedProps: Record<string, any>
}

// 模拟原生事件（浏览器传来的）
export interface NativeEvent {
  type: string
  target: NativeNode
}
```

### 2. 合成事件对象 (SyntheticEvent)

这是 React 暴露给开发者的事件对象。

```typescript
// ...existing code...

export class SyntheticEvent {
  nativeEvent: NativeEvent
  type: string
  target: NativeNode
  currentTarget: NativeNode | null = null // 当前处理事件的组件对应的 DOM

  _isPropagationStopped = false

  constructor(type: string, nativeEvent: NativeEvent) {
    this.type = type
    this.nativeEvent = nativeEvent
    this.target = nativeEvent.target
  }

  stopPropagation() {
    this._isPropagationStopped = true
  }

  isPropagationStopped() {
    return this._isPropagationStopped
  }
}
```

### 3. 核心引擎：事件传播与收集

这是 React 事件机制的心脏。它不直接遍历 DOM，而是遍历 Fiber 树。

```typescript
// ...existing code...

// 监听器封装
interface DispatchListener {
  instance: Fiber // 哪个组件实例
  listener: Function // 回调函数
  currentTarget: NativeNode // 对应的 DOM
}

export class ReactEventSystem {
  /**
   * 模拟 React 的根容器事件入口
   * 当原生事件触发时，React 会调用这个方法
   */
  public dispatchEvent(nativeEvent: NativeEvent): void {
    // 1. 找到触发事件的 DOM 对应的 Fiber 节点
    // 在真实 React 中，通过 internalInstanceKey 获取
    // 这里我们假设有一个 helper 函数能找到
    const targetFiber = this.getClosestInstanceFromNode(nativeEvent.target)

    if (!targetFiber) {
      console.warn('No fiber found for target node')
      return
    }

    // 2. 创建合成事件对象
    // React 会根据事件类型选择不同的合成类 (SyntheticMouseEvent, SyntheticKeyboardEvent 等)
    // 这里简化为统一类
    const syntheticEvent = new SyntheticEvent(nativeEvent.type, nativeEvent)

    // 3. 收集路径上的所有监听器 (两阶段)
    const dispatchQueue = this.extractEvents(targetFiber, syntheticEvent)

    // 4. 批量执行监听器
    this.processDispatchQueue(dispatchQueue, syntheticEvent)
  }

  /**
   * 核心算法：遍历 Fiber 树收集监听器
   */
  private extractEvents(
    targetFiber: Fiber,
    event: SyntheticEvent
  ): { capture: DispatchListener[]; bubble: DispatchListener[] } {
    const captureListeners: DispatchListener[] = []
    const bubbleListeners: DispatchListener[] = []

    let node: Fiber | null = targetFiber

    // 向上遍历 Fiber 树 (逻辑亲缘关系)
    while (node !== null) {
      const { memoizedProps, stateNode } = node

      if (stateNode) {
        // 只有 HostComponent (DOM 节点) 才有事件
        const eventName = event.type // e.g., 'click'

        // React 约定：onClickCapture
        const captureName = `on${capitalize(eventName)}Capture`
        const captureListener = memoizedProps[captureName]
        if (captureListener) {
          // 捕获阶段：父 -> 子。我们现在是从子 -> 父遍历，所以用 unshift (插到头部)
          captureListeners.unshift({
            instance: node,
            listener: captureListener,
            currentTarget: stateNode
          })
        }

        // React 约定：onClick
        const bubbleName = `on${capitalize(eventName)}`
        const bubbleListener = memoizedProps[bubbleName]
        if (bubbleListener) {
          // 冒泡阶段：子 -> 父。我们是从子 -> 父遍历，所以用 push (插到尾部)
          bubbleListeners.push({
            instance: node,
            listener: bubbleListener,
            currentTarget: stateNode
          })
        }
      }

      node = node.return // 沿着 Fiber 树向上
    }

    return { capture: captureListeners, bubble: bubbleListeners }
  }

  /**
   * 执行阶段
   */
  private processDispatchQueue(
    queue: { capture: DispatchListener[]; bubble: DispatchListener[] },
    event: SyntheticEvent
  ) {
    // 1. 执行捕获阶段
    for (const { listener, currentTarget } of queue.capture) {
      if (event.isPropagationStopped()) return

      event.currentTarget = currentTarget
      listener(event)
    }

    // 2. 执行冒泡阶段
    for (const { listener, currentTarget } of queue.bubble) {
      if (event.isPropagationStopped()) return

      event.currentTarget = currentTarget
      listener(event)
    }

    event.currentTarget = null // 清理
  }

  // 模拟：通过 DOM 节点反查 Fiber (在 React 源码中这是通过 DOM 属性 __reactFiber$xxx 实现的)
  // 为了演示，我们需要外部注入这个映射关系
  private instanceMap = new Map<NativeNode, Fiber>()

  public registerFiber(node: NativeNode, fiber: Fiber) {
    this.instanceMap.set(node, fiber)
  }

  private getClosestInstanceFromNode(node: NativeNode): Fiber | undefined {
    return this.instanceMap.get(node)
  }
}

// Helper
function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}
```

### 4. 验证与演示 (Portal 场景)

现在我们来模拟那个经典的 **Portal 冒泡** 场景，验证这套机制是否符合 React 的行为。

```typescript
import { NativeNode, Fiber, ReactEventSystem, NativeEvent } from './react-event-simulator'

// ==========================================
// 1. 构建物理世界 (DOM Tree)
// ==========================================
// 结构：body -> div#app -> div#portal-container
const body = new NativeNode('body')
const appDiv = new NativeNode('div#app')
const portalContainer = new NativeNode('div#portal-container') // Portal 挂载点
const button = new NativeNode('button#btn') // 实际按钮

body.appendChild(appDiv)
body.appendChild(portalContainer) // 注意：物理上 portalContainer 和 appDiv 是兄弟
portalContainer.appendChild(button) // 按钮在 portal 里

// ==========================================
// 2. 构建逻辑世界 (Fiber Tree)
// ==========================================
// 结构：App -> Parent -> Portal -> Button
// 逻辑上 Button 是 Parent 的后代，尽管物理上它在外面

const appFiber: Fiber = {
  tag: 'App',
  stateNode: appDiv,
  return: null,
  child: null,
  sibling: null,
  memoizedProps: {
    // App 监听冒泡
    onClick: (e: any) => console.log(`[React] App Clicked (target: ${e.target.id})`)
  }
}

const parentFiber: Fiber = {
  tag: 'Parent',
  stateNode: appDiv, // 假设 Parent 渲染在 appDiv 上
  return: appFiber,
  child: null,
  sibling: null,
  memoizedProps: {
    // Parent 监听捕获
    onClickCapture: (e: any) => console.log(`[React] Parent Capture (target: ${e.target.id})`)
  }
}

const buttonFiber: Fiber = {
  tag: 'Button',
  stateNode: button,
  return: parentFiber, // 关键：逻辑父级指向 Parent，而不是 null 或 portalContainer
  child: null,
  sibling: null,
  memoizedProps: {
    onClick: (e: any) => {
      console.log(`[React] Button Clicked`)
      // e.stopPropagation(); // 如果取消注释，App 将收不到事件
    }
  }
}

// ==========================================
// 3. 初始化系统并运行
// ==========================================
const reactSystem = new ReactEventSystem()

// 注册映射关系 (模拟 React 挂载过程)
reactSystem.registerFiber(appDiv, appFiber)
reactSystem.registerFiber(button, buttonFiber)

console.log('--- Simulating Click on Portal Button ---')

// 模拟原生事件触发
// 物理上，用户点击了 button
const nativeClickEvent: NativeEvent = {
  type: 'click',
  target: button
}

// React 介入：从根入口分发
reactSystem.dispatchEvent(nativeClickEvent)

/**
 * 预期输出：
 * --- Simulating Click on Portal Button ---
 * [React] Parent Capture (target: button#btn)  <-- 捕获阶段，逻辑父级收到
 * [React] Button Clicked                       <-- 目标阶段
 * [React] App Clicked (target: button#btn)     <-- 冒泡阶段，逻辑祖先收到
 */
```

### 设计亮点解析

1.  **逻辑优先 (`node = node.return`)**：
    在 `extractEvents` 中，我们通过 `node.return` 向上遍历。

    - 在 DOM 树中，`button` 的父级是 `portalContainer`，再往上是 `body`。
    - 在 Fiber 树中，`buttonFiber` 的 `return` 直接指向了 `parentFiber`。
    - 这就是为什么 **Portal 能冒泡** 的根本原因：React 遍历的是 Fiber 树，完全无视了 DOM 结构。

2.  **两阶段队列 (`captureListeners.unshift` vs `bubbleListeners.push`)**：
    我们只遍历了一次 Fiber 树（从底向上），但通过不同的数组操作方法，同时构建了捕获（从顶向下）和冒泡（从底向上）两个执行队列。这是 React 性能优化的一个细节。

3.  **合成事件封装**：
    `SyntheticEvent` 提供了 `stopPropagation`。在 `processDispatchQueue` 中，每次循环都会检查 `event.isPropagationStopped()`。一旦被调用，后续的监听器（无论是同一个节点的后续监听器，还是父级节点的监听器）都会被跳过。

4.  **环境无关**：
    整个实现没有使用 `window`、`document` 或 `HTMLElement`。所有的节点都是我们自定义的 `NativeNode` 类。这意味着这套逻辑可以运行在 Node.js、Canvas 渲染器甚至嵌入式设备中，只要你能提供对应的“原生节点”映射。
