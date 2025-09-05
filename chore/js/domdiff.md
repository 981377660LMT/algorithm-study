好的，我们来详细讲解 DOM Diff 算法的流程，并提供一个简化的 TypeScript 实现。

### 什么是 DOM Diff？

DOM Diff（DOM Difference Algorithm）是一种在现代前端框架（如 React, Vue）中广泛使用的核心算法。它的主要目标是**高效地更新用户界面**。

当应用的状态发生变化时，框架会根据新状态生成一个新的虚拟 DOM (Virtual DOM) 树。DOM Diff 算法会比较这棵新树和上一次渲染时生成的旧树，找出两棵树之间的差异，然后只将这些**最小化的差异**应用到真实的 DOM 上，从而避免了大规模、昂贵的 DOM 重绘和重排，极大地提升了性能。

### Diff 算法的核心思想与策略

直接、完整地比较两棵树的差异，其算法复杂度是 O(n³)，这个开销是无法接受的。因此，主流框架的 Diff 算法都建立在三个核心的启发式策略之上，将复杂度优化到了 O(n)。

1.  **只在同层级进行比较 (Tree Diff)**

    - 算法只会对同一层级的节点进行比较，不会跨层级移动节点。
    - 如果一个节点在旧树的 A 层，在新树的 B 层，那么算法不会尝试去“移动”它，而是会直接在 A 层删除旧节点，在 B 层创建新节点。
    - **实践**：如果一个 `<div>` 变成了 `<p>`，即使它们的子节点完全相同，框架也会销毁整个 `<div>` 及其所有子节点，然后创建一个新的 `<p>` 及其子节点。这大大简化了算法。

2.  **不同类型的组件会生成不同的树 (Component Diff)**

    - 如果两个节点的类型不同（例如，从 `<div>` 变成了 `<ComponentA>`），框架会直接将旧节点及其子树整个销毁，然后创建并挂载新节点。
    - 如果两个节点的类型相同（都是 `ComponentA`），框架会保留该组件实例，仅更新其 `props`，然后让组件实例内部继续执行其自身的 Diff 过程。

3.  **通过 `key` 属性来识别相同节点 (Element Diff)**
    - 对于一层级的多个子节点，`key` 是至关重要的。`key` 帮助 Diff 算法识别哪些节点是“同一个”节点，只是位置发生了变化。
    - **无 `key` 的情况**：算法会按顺序比较，如果发现不一致，就会就地修改。这在列表反转或中间插入/删除元素的场景下效率极低，会导致大量不必要的 DOM 操作。
    - **有 `key` 的情况**：算法会创建一个 `key` 到旧节点位置的映射。然后遍历新节点列表，通过 `key` 快速找到对应的旧节点。这样就可以高效地处理节点的移动、新增和删除。

### Diff 算法流程详解 (以双端比较为例)

双端比较（Double-Ended Comparison）是 Vue 2 中使用的一种高效的子节点比较策略。它通过在旧子节点列表和新子节点列表的两端同时设置指针，并进行比较，来最大化地复用节点。

假设我们有 `oldChildren` 和 `newChildren` 两个数组。

**定义四个指针：**

- `oldStartIdx`: 指向 `oldChildren` 的头部
- `oldEndIdx`: 指向 `oldChildren` 的尾部
- `newStartIdx`: 指向 `newChildren` 的头部
- `newEndIdx`: 指向 `newChildren` 的尾部

**循环比较过程 (while `oldStartIdx <= oldEndIdx` && `newStartIdx <= newEndIdx`)**

在循环中，进行以下四种核心比较：

1.  **`oldStartVnode` vs `newStartVnode`**

    - 如果它们是同一个节点（`key` 和 `tag` 相同），则对它们进行 `patch`（更新属性、递归处理子节点），然后两个 `start` 指针都向后移动一位 (`++`)。

2.  **`oldEndVnode` vs `newEndVnode`**

    - 如果它们是同一个节点，则 `patch` 它们，然后两个 `end` 指针都向前移动一位 (`--`)。

3.  **`oldStartVnode` vs `newEndVnode`**

    - 如果它们是同一个节点，说明旧的头部节点移动到了新的尾部。`patch` 它们，然后将该 DOM 元素移动到当前旧尾部元素的**后面**。`oldStartIdx++`，`newEndIdx--`。

4.  **`oldEndVnode` vs `newStartVnode`**
    - 如果它们是同一个节点，说明旧的尾部节点移动到了新的头部。`patch` 它们，然后将该 DOM 元素移动到当前旧头部元素的**前面**。`oldEndIdx--`，`newStartIdx++`。

**如果以上四种情况都不匹配：**

- 创建一个从 `key` 到 `oldChildren` 索引的映射 `(key -> index)`。
- 用 `newStartVnode` 的 `key` 在映射中查找。
  - **如果找不到**：说明 `newStartVnode` 是一个全新的节点，创建它并插入到 `oldStartVnode` 对应的 DOM 元素**前面**。
  - **如果找到了**：说明这个节点需要被移动。`patch` 这两个节点，然后将找到的旧节点对应的 DOM 元素移动到 `oldStartVnode` 对应的 DOM 元素**前面**。同时，将 `oldChildren` 中该位置的旧节点设为 `undefined`，以标记它已被处理。
- 最后，`newStartIdx++`。

**循环结束后的处理：**

- 如果 `oldStartIdx > oldEndIdx`，说明 `oldChildren` 已经遍历完，但 `newChildren` 还有剩余。这些剩余的节点都是需要新增的，将它们批量插入。
- 如果 `newStartIdx > newEndIdx`，说明 `newChildren` 已经遍历完，但 `oldChildren` 还有剩余。这些剩余的节点（不为 `undefined` 的）都是需要删除的，将它们批量删除。

### TypeScript 代码实现

下面是一个简化的 DOM Diff 算法的 TS 实现，主要关注 `patch` 和 `updateChildren` 逻辑。

```typescript
// 1. 定义虚拟节点接口
export interface VNode {
  tag?: string // 标签名，如 'div'
  props?: { [key: string]: any } // 属性，如 { id: 'app' }
  children?: string | VNode[] // 子节点
  key?: string | number // key，用于 diff
  el?: Node // 对应的真实 DOM 节点
}

// 2. 创建虚拟节点的辅助函数
export function h(
  tag: string,
  props: { [key: string]: any } | null,
  children: string | VNode[]
): VNode {
  return { tag, props, children, key: props?.key }
}

// 3. 挂载函数：将 VNode 转化为真实 DOM
export function mount(vnode: VNode, container: Node): void {
  const el = (vnode.el = document.createElement(vnode.tag!))

  // 处理 props
  if (vnode.props) {
    for (const key in vnode.props) {
      const value = vnode.props[key]
      if (key.startsWith('on')) {
        el.addEventListener(key.slice(2).toLowerCase(), value)
      } else {
        el.setAttribute(key, value)
      }
    }
  }

  // 处理 children
  if (typeof vnode.children === 'string') {
    el.textContent = vnode.children
  } else if (Array.isArray(vnode.children)) {
    vnode.children.forEach(child => {
      mount(child, el)
    })
  }

  container.appendChild(el)
}

// 4. Diff 核心函数：patch
export function patch(n1: VNode, n2: VNode): void {
  const el = (n2.el = n1.el!) // 复用 DOM 节点

  // a. 如果标签不同，直接替换
  if (n1.tag !== n2.tag) {
    const parent = el.parentNode!
    mount(n2, parent)
    parent.removeChild(el)
    return
  }

  // b. 标签相同，处理 props
  const oldProps = n1.props || {}
  const newProps = n2.props || {}

  // 更新/新增 props
  for (const key in newProps) {
    if (newProps[key] !== oldProps[key]) {
      if (key.startsWith('on')) {
        ;(el as HTMLElement).addEventListener(key.slice(2).toLowerCase(), newProps[key])
      } else {
        ;(el as HTMLElement).setAttribute(key, newProps[key])
      }
    }
  }
  // 删除不存在的 props
  for (const key in oldProps) {
    if (!(key in newProps)) {
      if (key.startsWith('on')) {
        ;(el as HTMLElement).removeEventListener(key.slice(2).toLowerCase(), oldProps[key])
      } else {
        ;(el as HTMLElement).removeAttribute(key)
      }
    }
  }

  // c. 处理 children
  const oldChildren = n1.children
  const newChildren = n2.children

  if (typeof newChildren === 'string') {
    // Case 1: 新子节点是文本
    if (typeof oldChildren === 'string') {
      if (newChildren !== oldChildren) {
        el.textContent = newChildren
      }
    } else {
      el.textContent = newChildren
    }
  } else if (Array.isArray(newChildren)) {
    // Case 2: 新子节点是数组
    if (typeof oldChildren === 'string') {
      el.innerHTML = ''
      newChildren.forEach(child => mount(child, el))
    } else if (Array.isArray(oldChildren)) {
      // *** 这是最复杂的部分：子节点 Diff ***
      updateChildren(el, oldChildren, newChildren)
    }
  }
}

function isSameVNodeType(n1: VNode, n2: VNode): boolean {
  return n1.tag === n2.tag && n1.key === n2.key
}

// 5. 子节点 Diff 核心：双端比较
function updateChildren(parentEl: Node, oldChildren: VNode[], newChildren: VNode[]) {
  let oldStartIdx = 0
  let newStartIdx = 0
  let oldEndIdx = oldChildren.length - 1
  let newEndIdx = newChildren.length - 1

  let oldStartVnode = oldChildren[0]
  let oldEndVnode = oldChildren[oldEndIdx]
  let newStartVnode = newChildren[0]
  let newEndVnode = newChildren[newEndIdx]

  let oldKeyToIdx: { [key: string]: number } | undefined

  while (oldStartIdx <= oldEndIdx && newStartIdx <= newEndIdx) {
    if (oldStartVnode === undefined) {
      oldStartVnode = oldChildren[++oldStartIdx]
    } else if (oldEndVnode === undefined) {
      oldEndVnode = oldChildren[--oldEndIdx]
    } else if (isSameVNodeType(oldStartVnode, newStartVnode)) {
      // 1. oldStart vs newStart
      patch(oldStartVnode, newStartVnode)
      oldStartVnode = oldChildren[++oldStartIdx]
      newStartVnode = newChildren[++newStartIdx]
    } else if (isSameVNodeType(oldEndVnode, newEndVnode)) {
      // 2. oldEnd vs newEnd
      patch(oldEndVnode, newEndVnode)
      oldEndVnode = oldChildren[--oldEndIdx]
      newEndVnode = newChildren[--newEndIdx]
    } else if (isSameVNodeType(oldStartVnode, newEndVnode)) {
      // 3. oldStart vs newEnd
      patch(oldStartVnode, newEndVnode)
      parentEl.insertBefore(oldStartVnode.el!, oldEndVnode.el!.nextSibling)
      oldStartVnode = oldChildren[++oldStartIdx]
      newEndVnode = newChildren[--newEndIdx]
    } else if (isSameVNodeType(oldEndVnode, newStartVnode)) {
      // 4. oldEnd vs newStart
      patch(oldEndVnode, newStartVnode)
      parentEl.insertBefore(oldEndVnode.el!, oldStartVnode.el!)
      oldEndVnode = oldChildren[--oldEndIdx]
      newStartVnode = newChildren[++newStartIdx]
    } else {
      // 5. 四种情况都不匹配，使用 key 查找
      if (!oldKeyToIdx) {
        oldKeyToIdx = {}
        for (let i = oldStartIdx; i <= oldEndIdx; i++) {
          const key = oldChildren[i].key
          if (key != null) oldKeyToIdx[key] = i
        }
      }
      const idxInOld = oldKeyToIdx[newStartVnode.key as string]
      if (idxInOld === undefined) {
        // 在旧节点中找不到，是新节点
        mount(newStartVnode, parentEl)
        parentEl.insertBefore(newStartVnode.el!, oldStartVnode.el!)
      } else {
        // 找到了，移动节点
        const vnodeToMove = oldChildren[idxInOld]
        patch(vnodeToMove, newStartVnode)
        parentEl.insertBefore(vnodeToMove.el!, oldStartVnode.el!)
        oldChildren[idxInOld] = undefined as any // 标记为已处理
      }
      newStartVnode = newChildren[++newStartIdx]
    }
  }

  // 循环结束后，处理剩余节点
  if (oldStartIdx > oldEndIdx) {
    // 新增剩余的新节点
    const anchor = newChildren[newEndIdx + 1] ? newChildren[newEndIdx + 1].el : null
    for (let i = newStartIdx; i <= newEndIdx; i++) {
      mount(newChildren[i], parentEl)
      parentEl.insertBefore(newChildren[i].el!, anchor)
    }
  } else if (newStartIdx > newEndIdx) {
    // 删除多余的旧节点
    for (let i = oldStartIdx; i <= oldEndIdx; i++) {
      if (oldChildren[i]) {
        parentEl.removeChild(oldChildren[i].el!)
      }
    }
  }
}
```

---

好的，这是一个非常深入的问题。标准 Diff 算法**刻意避免**跨层级比较，是因为这样做会导致算法复杂度急剧上升，从 O(n) 退化到 O(n³)，这在实践中是不可接受的。

不过，我们可以从理论上探讨一下，如果**必须**实现跨层级比较和移动（swap），需要怎么做。这通常需要一个完全不同的策略。

### 为什么标准 Diff 不支持跨层级比较？

回顾一下，标准 Diff 的 O(n) 性能依赖于一个核心假设：**节点只会在同一层级的兄弟节点之间移动、新增或删除**。这个假设大大缩小了搜索范围。例如，比较 `<A>` 的子节点时，算法只会在 `<A>` 的旧子节点列表里去寻找匹配项，而不会去 `<B>` 或其他任何地方寻找。

一旦打破这个假设，对于新树中的任意一个节点，你都必须去旧树的**所有层级**中搜索是否存在可复用的节点。这就是一个完全的“树编辑距离”（Tree Edit Distance）问题，计算开销极大。

### 实现跨层级比较的理论方案

要实现跨层级移动，你需要一种方法来在**整个树的范围内**唯一标识一个节点，并跟踪它的移动。基于 `key` 的同级比较策略就不够用了。

以下是一种可能的实现思路，我们称之为**“全局节点映射”**策略：

#### 1. 引入全局唯一 ID (Global Unique ID)

首先，你需要一个能在整个虚拟 DOM 树中唯一标识一个组件或元素的 ID，而不仅仅是在兄弟节点中唯一。我们可以称之为 `globalId`。这个 `globalId` 必须在节点的整个生命周期中保持不变。

```typescript
interface VNode {
  tag?: string
  props?: { [key: string]: any }
  children?: string | VNode[]
  key?: string | number // 兄弟节点中唯一
  globalId?: string // 整棵树中唯一
  el?: Node
}
```

#### 2. 算法流程

算法将分为三个主要阶段：

**阶段一：构建旧节点映射**

1.  在开始 `patch` 之前，先**深度优先遍历（DFS）整棵旧 VDOM 树**。
2.  创建一个 `Map`，将每个带有 `globalId` 的旧节点的 `globalId` 作为 key，节点本身 (`VNode`) 作为 value。
3.  这个 Map (`oldNodesMap: Map<string, VNode>`) 存储了所有可能被复用的节点。

**阶段二：遍历新树，处理节点**

1.  **深度优先遍历新 VDOM 树**。对于新树中的每一个节点 (`newNode`)：
2.  检查它是否有 `globalId`。
3.  在 `oldNodesMap` 中查找这个 `globalId`。

    - **情况 A：找到了匹配的旧节点 (`oldNode`)**

      - 这说明节点被**复用**了，它可能被移动了，也可能在原地更新。
      - **复用 DOM 元素**：`newNode.el = oldNode.el`。
      - **检查父节点是否变化**：比较 `newNode` 的父节点和 `oldNode` 的父节点。如果父节点不同，说明发生了**跨层级移动**。你需要执行真实的 DOM 移动操作：`newParentEl.appendChild(newNode.el)`。_(注意：这里为了简化只用了 appendChild，实际情况需要更精确的 insertBefore 来保证顺序)_。
      - **执行 `patch`**：对 `oldNode` 和 `newNode` 执行属性和样式的 `patch` 操作（但不递归处理子节点，因为子节点也会在这个全局遍历中被独立处理）。
      - **从 Map 中删除**：从 `oldNodesMap` 中删除这个 `globalId`，表示它已经被处理过了。

    - **情况 B：没有找到匹配的旧节点**
      - 这说明 `newNode` 是一个**全新的节点**。
      - **创建新 DOM**：执行 `mount(newNode, newParentEl)` 来创建并挂载新的 DOM 元素。

**阶段三：卸载废弃的节点**

1.  在遍历完新树后，`oldNodesMap` 中剩下的所有节点都是在新树中不存在的、需要被**废弃**的节点。
2.  遍历这个 Map，对其中每一个 `VNode`，从其父节点上执行 `removeChild(vnode.el)` 操作，将其卸载。

### 伪代码实现

```typescript
// 全局映射
let oldNodesMap: Map<string, VNode>

function buildOldNodesMap(node: VNode) {
  if (node.globalId) {
    oldNodesMap.set(node.globalId, node)
  }
  if (Array.isArray(node.children)) {
    node.children.forEach(buildOldNodesMap)
  }
}

function patchWithGlobalMove(oldTree: VNode, newTree: VNode, container: Node) {
  oldNodesMap = new Map()
  // 阶段一：构建旧节点 Map
  buildOldNodesMap(oldTree)

  // 阶段二：遍历新树
  patchNewNode(newTree, container)

  // 阶段三：卸载废弃节点
  for (const oldNode of oldNodesMap.values()) {
    if (oldNode.el && oldNode.el.parentNode) {
      oldNode.el.parentNode.removeChild(oldNode.el)
    }
  }
}

function patchNewNode(newNode: VNode, parentEl: Node) {
  const oldNode = newNode.globalId ? oldNodesMap.get(newNode.globalId) : undefined

  if (oldNode) {
    // 情况 A: 节点复用
    const el = (newNode.el = oldNode.el!)

    // 比较 props (简化版)
    // ... update props ...

    // 检查是否需要移动
    if (el.parentNode !== parentEl) {
      parentEl.appendChild(el) // 移动 DOM 节点
    }

    oldNodesMap.delete(newNode.globalId!) // 标记为已处理

    // 递归处理子节点
    if (Array.isArray(newNode.children)) {
      newNode.children.forEach(child => patchNewNode(child, el))
    }
  } else {
    // 情况 B: 新增节点
    mount(newNode, parentEl) // mount 会创建 DOM 并递归挂载子节点
  }
}
```

### 总结与权衡

| 特性         | 标准 Diff (同层比较)       | 全局移动 Diff                                                                                               |
| :----------- | :------------------------- | :---------------------------------------------------------------------------------------------------------- |
| **核心思想** | 假设结构稳定，只在同级比较 | 假设节点可任意移动，全局搜索                                                                                |
| **复杂度**   | **O(n)**                   | 接近 **O(n²)** 或更差 (取决于 Map 操作和 DOM 操作)                                                          |
| **性能**     | **极高**                   | **较低**，不适用于频繁更新                                                                                  |
| **适用场景** | 几乎所有现代 UI 框架       | 理论探讨，或用于某些特定场景（如动画库 `Framer Motion` 的 `layout` 动画，它会独立跟踪元素位置实现类似效果） |
| **实现难度** | 复杂，尤其双端比较逻辑     | 概念更直接，但工程实现细节多                                                                                |

**结论：** 虽然理论上可以实现跨层级比较和移动，但由于其巨大的性能开销，现代主流前端框架无一例外地选择了**放弃**这个功能，转而采用更具性价比的**同层比较策略**。框架的设计者们认为，为了处理极少数的跨层级移动场景而牺牲掉所有场景下的高性能是不值得的。开发者应该遵循框架的设计理念，通过改变状态和数据结构来避免需要跨层级移动节点的 UI 设计。
