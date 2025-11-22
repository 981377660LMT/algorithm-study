这是一个极具挑战性的任务。要不依赖浏览器实现 DOM 事件机制，我们需要从零构建一套符合 **W3C DOM Level 2 Events** 标准的模拟引擎。

这不仅仅是一个简单的发布订阅模式，它必须包含：

1.  **树形结构**：节点必须有父子关系。
2.  **事件对象**：包含 target, `currentTarget`, `eventPhase`, `stopPropagation`, `preventDefault` 等标准属性。
3.  **完整事件流**：必须严格遵循 **捕获 (Capturing) -> 目标 (Target) -> 冒泡 (Bubbling)** 的传播路径。
4.  **监听器选项**：支持 `capture: true` (捕获阶段监听) 和 `once` 等选项。

以下是完整的 TypeScript 实现：

### 1. 核心类型定义与枚举

首先定义事件流的三个阶段和基础接口。

```typescript
// W3C 标准事件阶段常量
export enum EventPhase {
  NONE = 0,
  CAPTURING_PHASE = 1,
  AT_TARGET = 2,
  BUBBLING_PHASE = 3
}

// 监听器选项接口
export interface AddEventListenerOptions {
  capture?: boolean
  once?: boolean
}

// 统一处理监听器参数（支持 boolean 或 对象）
type EventListenerOptionsOrUseCapture = boolean | AddEventListenerOptions

// 事件回调函数类型
export interface EventListener {
  (evt: SimulatorEvent): void
}

// 内部存储的监听器结构
interface ListenerEntry {
  callback: EventListener
  capture: boolean
  once: boolean
}
```

### 2. 模拟 Event 对象

我们需要模拟 `Event` 类，实现 `stopPropagation` 和 `stopImmediatePropagation` 的状态管理。

```typescript
// ...existing code...

export class SimulatorEvent {
  public type: string
  public target: EventTargetNode | null = null
  public currentTarget: EventTargetNode | null = null
  public eventPhase: EventPhase = EventPhase.NONE
  public bubbles: boolean
  public cancelable: boolean

  // 默认行为是否被阻止
  public defaultPrevented: boolean = false

  // 内部标志位
  private _propagationStopped: boolean = false
  private _immediatePropagationStopped: boolean = false

  constructor(type: string, eventInitDict?: { bubbles?: boolean; cancelable?: boolean }) {
    this.type = type
    this.bubbles = eventInitDict?.bubbles ?? false
    this.cancelable = eventInitDict?.cancelable ?? false
  }

  public stopPropagation(): void {
    this._propagationStopped = true
  }

  public stopImmediatePropagation(): void {
    this._propagationStopped = true
    this._immediatePropagationStopped = true
  }

  public preventDefault(): void {
    if (this.cancelable) {
      this.defaultPrevented = true
    }
  }

  // 供内部引擎检查状态
  public get isPropagationStopped() {
    return this._propagationStopped
  }
  public get isImmediatePropagationStopped() {
    return this._immediatePropagationStopped
  }
}
```

### 3. 模拟 DOM 节点 (EventTargetNode)

这是核心部分。我们需要实现树结构管理 (`parentNode`) 和 `dispatchEvent` 的完整逻辑。

```typescript
// ...existing code...

export class EventTargetNode {
  public parentNode: EventTargetNode | null = null
  public name: string // 用于调试识别

  // 存储监听器：Map<事件名, 监听器列表>
  private _listeners: Map<string, ListenerEntry[]> = new Map()

  constructor(name: string) {
    this.name = name
  }

  // 模拟 appendChild 以建立树形关系
  public appendChild(child: EventTargetNode): void {
    child.parentNode = this
  }

  public addEventListener(
    type: string,
    callback: EventListener,
    options?: EventListenerOptionsOrUseCapture
  ): void {
    if (!callback) return

    // 规范化参数
    let capture = false
    let once = false

    if (typeof options === 'boolean') {
      capture = options
    } else if (typeof options === 'object') {
      capture = !!options.capture
      once = !!options.once
    }

    if (!this._listeners.has(type)) {
      this._listeners.set(type, [])
    }

    const list = this._listeners.get(type)!

    // 避免重复添加相同的 (callback, capture) 对
    const exists = list.some(l => l.callback === callback && l.capture === capture)
    if (!exists) {
      list.push({ callback, capture, once })
    }
  }

  public removeEventListener(
    type: string,
    callback: EventListener,
    options?: EventListenerOptionsOrUseCapture
  ): void {
    const list = this._listeners.get(type)
    if (!list) return

    const capture = typeof options === 'boolean' ? options : !!options?.capture

    const index = list.findIndex(l => l.callback === callback && l.capture === capture)
    if (index !== -1) {
      list.splice(index, 1)
    }
  }

  /**
   * 核心：派发事件
   * 严格遵循 W3C 事件流：捕获 -> 目标 -> 冒泡
   */
  public dispatchEvent(event: SimulatorEvent): boolean {
    // 1. 初始化事件状态
    event.target = this
    event.eventPhase = EventPhase.NONE

    // 重置停止标志（如果是复用的事件对象）
    // 但通常 dispatchEvent 应该接收一个新的事件对象

    // 2. 构建传播路径 (Propagation Path)
    // 从当前节点一直向上找到根节点
    const propagationPath: EventTargetNode[] = []
    let current: EventTargetNode | null = this
    while (current) {
      propagationPath.push(current)
      current = current.parentNode
    }

    // 路径是 [Target, Parent, GrandParent, Root]
    // 捕获需要反向遍历，冒泡需要正向遍历

    // ==========================================
    // 阶段 1: 捕获阶段 (Capturing Phase)
    // ==========================================
    event.eventPhase = EventPhase.CAPTURING_PHASE
    // 从 Root 到 Parent (不包含 Target)
    for (let i = propagationPath.length - 1; i > 0; i--) {
      if (event.isPropagationStopped) break
      const node = propagationPath[i]
      event.currentTarget = node
      node._invokeListeners(event)
    }

    // ==========================================
    // 阶段 2: 目标阶段 (Target Phase)
    // ==========================================
    if (!event.isPropagationStopped) {
      event.eventPhase = EventPhase.AT_TARGET
      event.currentTarget = this
      this._invokeListeners(event)
    }

    // ==========================================
    // 阶段 3: 冒泡阶段 (Bubbling Phase)
    // ==========================================
    if (event.bubbles && !event.isPropagationStopped) {
      event.eventPhase = EventPhase.BUBBLING_PHASE
      // 从 Parent 到 Root
      for (let i = 1; i < propagationPath.length; i++) {
        if (event.isPropagationStopped) break
        const node = propagationPath[i]
        event.currentTarget = node
        node._invokeListeners(event)
      }
    }

    // 3. 清理
    event.currentTarget = null
    event.eventPhase = EventPhase.NONE

    return !event.defaultPrevented
  }

  /**
   * 内部方法：在当前节点触发监听器
   */
  private _invokeListeners(event: SimulatorEvent): void {
    const list = this._listeners.get(event.type)
    if (!list || list.length === 0) return

    // 创建副本以应对回调中移除监听器的情况
    const listenersCopy = [...list]

    for (const entry of listenersCopy) {
      // 如果 stopImmediatePropagation 被调用，停止当前节点剩余监听器的执行
      if (event.isImmediatePropagationStopped) return

      // 匹配阶段：
      // 1. 捕获阶段：只执行 capture 为 true 的监听器
      // 2. 冒泡阶段：只执行 capture 为 false 的监听器
      // 3. 目标阶段：两者都执行 (但在浏览器实现中，目标阶段通常按注册顺序执行，不分 capture)

      const isCapturingPhase = event.eventPhase === EventPhase.CAPTURING_PHASE
      const isBubblingPhase = event.eventPhase === EventPhase.BUBBLING_PHASE

      // 逻辑判断：
      // 如果是捕获阶段，必须是 capture 监听器
      // 如果是冒泡阶段，必须是 非capture 监听器
      // 如果是目标阶段，通吃
      if (isCapturingPhase && !entry.capture) continue
      if (isBubblingPhase && entry.capture) continue

      // 处理 once
      if (entry.once) {
        this.removeEventListener(event.type, entry.callback, { capture: entry.capture })
      }

      try {
        entry.callback.call(this, event)
      } catch (err) {
        console.error(`Error in event handler on node ${this.name}:`, err)
      }
    }
  }
}
```

### 4. 验证与演示

现在我们编写一段代码来验证这套机制是否符合预期。

```typescript
import { EventTargetNode, SimulatorEvent, EventPhase } from './dom-event-simulator'

// 1. 构建 DOM 树
// root -> parent -> child
const root = new EventTargetNode('Root')
const parent = new EventTargetNode('Parent')
const child = new EventTargetNode('Child')

root.appendChild(parent)
parent.appendChild(child)

// 2. 辅助日志函数
const log = (msg: string, e: SimulatorEvent) => {
  const phaseMap = {
    [EventPhase.NONE]: 'NONE',
    [EventPhase.CAPTURING_PHASE]: 'CAPTURE',
    [EventPhase.AT_TARGET]: 'TARGET',
    [EventPhase.BUBBLING_PHASE]: 'BUBBLE'
  }
  console.log(`[${phaseMap[e.eventPhase]}] ${e.currentTarget?.name}: ${msg}`)
}

// 3. 绑定事件

// Root: 捕获
root.addEventListener('click', e => log('Root Capturing', e), true)

// Parent: 捕获 (演示 stopPropagation)
parent.addEventListener(
  'click',
  e => {
    log('Parent Capturing', e)
    // Uncomment next line to stop event from reaching Child
    // e.stopPropagation();
  },
  true
)

// Parent: 冒泡
parent.addEventListener('click', e => log('Parent Bubbling', e), false)

// Child: 目标阶段 (注册两个，演示 stopImmediatePropagation)
child.addEventListener(
  'click',
  e => {
    log('Child Handler 1', e)
    // e.stopImmediatePropagation(); // 如果取消注释，Handler 2 和 Parent Bubbling 都不会执行
  },
  false
)

child.addEventListener(
  'click',
  e => {
    log('Child Handler 2', e)
  },
  false
)

// 4. 触发事件
console.log('--- Dispatching Click Event on Child ---')
const event = new SimulatorEvent('click', { bubbles: true }) // 必须开启 bubbles 才能冒泡
child.dispatchEvent(event)

/**
 * 预期输出:
 * --- Dispatching Click Event on Child ---
 * [CAPTURE] Root: Root Capturing
 * [CAPTURE] Parent: Parent Capturing
 * [TARGET] Child: Child Handler 1
 * [TARGET] Child: Child Handler 2
 * [BUBBLE] Parent: Parent Bubbling
 * (Root 没有冒泡监听器，所以不输出)
 */
```

### 设计要点解析

1.  **路径计算 (`propagationPath`)**：
    在 `dispatchEvent` 开始时，我们必须先计算出完整的路径。这是因为 DOM 标准规定，如果在事件传播过程中 DOM 结构发生了变化（例如在捕获阶段把目标节点移除了），事件传播路径**不应该**改变。预先计算路径符合这一标准。

2.  **阶段过滤**：
    在 `_invokeListeners` 中，最关键的逻辑是：

    ```typescript
    if (isCapturingPhase && !entry.capture) continue
    if (isBubblingPhase && entry.capture) continue
    ```

    这确保了捕获监听器只在捕获阶段触发，冒泡监听器只在冒泡阶段触发。而在 `AT_TARGET` 阶段，两者都会触发（虽然现代浏览器通常按注册顺序，但标准允许实现差异，这里简化为都触发）。

3.  **`stopImmediatePropagation` vs `stopPropagation`**：

    - `stopPropagation` 设置 `_propagationStopped`，这会在 `dispatchEvent` 的大循环（父子节点遍历）中被检测到，从而停止向上传播。
    - `stopImmediatePropagation` 设置 `_immediatePropagationStopped`，这会在 `_invokeListeners` 的小循环（当前节点监听器列表）中被检测到，从而停止当前节点后续监听器的执行。

4.  **`this` 上下文**：
    使用 `entry.callback.call(this, event)` 确保回调函数中的 `this` 指向当前的 `currentTarget`，这与原生 DOM 行为一致。
