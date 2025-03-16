# React Fiber 遍历算法详解

这段代码实现了 React Fiber 架构中的深度优先迭代遍历算法，让我详细解析其工作原理。

## Fiber 节点结构

```typescript
type Fiber = {
  child: Fiber | undefined // 指向第一个子节点
  sibling: Fiber | undefined // 指向下一个兄弟节点
  return: Fiber | undefined // 指向父节点
}
```

这是简化版的 Fiber 节点，通过三个指针形成了一个可遍历的树结构：

-

child

: 连接到第一个子节点

-

sibling

: 连接到同级的下一个节点

-

return

: 连接到父节点(在 React 中命名为 return 是因为它表示执行完当前节点后应返回的节点)

## 遍历算法分析

```typescript
function enumerateFiber(root: Fiber, f: (f: Fiber) => void): void {
  let node = root
  while (true) {
    f(node) // 1. 处理当前节点
    if (node.child) {
      // 2. 优先遍历子节点
      node = node.child
      continue
    }
    if (node === root) return // 3. 如果回到根节点，结束遍历
    while (!node.sibling) {
      // 4. 如果没有兄弟节点，向上回溯
      if (!node.return || node.return === root) return
      node = node.return
    }
    node = node.sibling // 5. 处理兄弟节点
  }
}
```

### 算法步骤详解

1. **处理当前节点**：对当前访问的节点执行回调函数

f

2. **深度优先原则**：

   - 如果有子节点，立即深入到子节点
   - 没有子节点才考虑兄弟节点或回溯

3. **特殊情况处理**：

   - `if (node === root) return` 处理根节点没有子节点的边缘情况

4. **回溯过程**：

   - 当一个节点没有子节点也没有兄弟节点时，需要向上回溯
   - `while (!node.sibling)` 循环持续向上回溯，直到找到一个有兄弟节点的祖先
   - `if (!node.return || node.return === root) return` 终止条件：回溯到根节点或遍历结束

5. **横向遍历**：
   - 找到有兄弟节点的祖先后，转向其兄弟节点继续遍历

### 遍历顺序示例

对于这样的一个树结构：

```
    A
   / \
  B   C
 / \   \
D   E   F
```

遍历顺序为：A → B → D → E → C → F

### 算法特点

1. **迭代实现**：避免了递归实现可能导致的调用栈溢出问题

2. **易于中断恢复**：可以在任何节点暂停并记录当前状态，这是 React Fiber 实现时间切片的基础

3. **空间复杂度低**：O(1) 额外空间，仅使用指针导航

4. **通用性强**：可以用于任何树形结构的遍历，只要实现了相应的节点关系

这种迭代遍历方式是 React Fiber 架构能够实现增量渲染和优先级调度的基石，允许 React 在渲染大型组件树时保持应用的响应性。
