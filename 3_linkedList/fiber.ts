// 单链表树遍历算法-fiber
// https://github.com/facebook/react/issues/7942
// React Fiber 架构中的深度优先迭代遍历算法

/**
 * 多叉树链表.
 */
type Fiber = {
  child: Fiber | undefined // 指向第一个子节点
  sibling: Fiber | undefined // 指向下一个兄弟节点
  return: Fiber | undefined // 指向父节点
}

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

export {}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const root: Fiber = {
    child: {
      child: {
        child: undefined,
        sibling: undefined,
        return: undefined
      },
      sibling: {
        child: undefined,
        sibling: undefined,
        return: undefined
      },
      return: undefined
    },
    sibling: {
      child: {
        child: undefined,
        sibling: undefined,
        return: undefined
      },
      sibling: {
        child: undefined,
        sibling: undefined,
        return: undefined
      },
      return: undefined
    },
    return: undefined
  }

  enumerateFiber(root, f => console.log(f))
}
