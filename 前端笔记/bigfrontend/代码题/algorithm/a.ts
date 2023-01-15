// 给定一个DOM tree和目标节点，请找出其左边的节点。

function previousLeftSibling(root: Element, target: Element): Element | null {
  let queue: Element[] = [root]
  while (queue.length) {
    const nextQueue: Element[] = []
    const n = queue.length
    let pre: Element | null = null
    for (let i = 0; i < n; i++) {
      const cur = queue[i]
      if (cur === target) return pre
      pre = cur
      nextQueue.push(...Array.from(cur.children))
    }

    queue = nextQueue
  }

  return null
}
