// 单链表树遍历算法-fiber
// https://github.com/facebook/react/issues/7942

type Fiber = {
  child: Fiber | undefined
  sibling: Fiber | undefined
  return: Fiber | undefined
}

function enumerateFiber(root: Fiber, f: (f: Fiber) => void): void {
  let node = root
  while (true) {
    f(node)
    if (node.child) {
      node = node.child
      continue
    }
    if (node === root) return
    while (!node.sibling) {
      if (!node.return || node.return === root) return
      node = node.return
    }
    node = node.sibling
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
