/**
 * @param {Element | null} tree
 * @return {number}
 */
function getHeight(tree: Element | null): number {
  // your code here
  if (!tree) return 0
  return 1 + Math.max(...Array.from(tree.children).map(getHeight), 0)
}

function getHeight2(tree: Element | null): number {
  // your code here
  let height = 0
  if (!tree) return height

  const queue: [Element, number][] = [[tree, 1]]
  while (queue.length) {
    const [node, h] = queue.shift()!
    height = Math.max(h, height)
    for (const child of Array.from(node.children)) {
      queue.push([child, h + 1])
    }
  }

  return height
}
