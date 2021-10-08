/**
 * @param {HTMLElement} root
 * @param {HTMLElement} target
 * @return {HTMLElemnt | null}
 * 给定一个DOM tree和目标节点，请找出其右边的节点。
 * @description
 * bfs解法
 */
function nextRightSibling(root: HTMLElement, target: HTMLElement): Element | null {
  if (!root) return null

  // 每一层的pre
  let pre: Element | null = null
  const queue: Element[] = [root]

  while (queue.length) {
    const len = queue.length
    for (let i = 0; i < len; i++) {
      const node = queue.shift()!
      if (pre === target) return node
      pre = node
      queue.push(...Array.from(node.children))
    }
    pre = null
  }

  return null
}
