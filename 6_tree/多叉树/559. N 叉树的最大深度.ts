class Node {
  val: number
  children: Node[]
  constructor(val = 0) {
    this.val = val
    this.children = []
  }
}

function maxDepth(root: Node | null): number {
  if (!root) return 0
  if (!root.children) return 1
  return Math.max(...root.children.map(child => maxDepth(child)), 0) + 1
}

export {}
