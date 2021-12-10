class Node {
  val: number
  children: Node[]

  constructor(val = 0, children: Node[] = []) {
    this.val = val
    this.children = children
  }
}

// 每个节点都有唯一的值。
// 你可以使用 O(1) 额外内存空间且 O(n) 时间复杂度的算法来找到该树的根节点吗？
// 1.解法1：set存储，看看有没有父亲
// 2.解法2：除了根节点，其他值都出现了两次
function findRoot(tree: Node[]): Node | null {
  const childNode = new WeakSet<Node>()
  for (const node of tree) {
    for (const child of node.children) {
      childNode.add(child)
    }
  }

  for (const node of tree) {
    if (!childNode.has(node)) return node
  }

  throw new Error('invalid tree')
}

function findRoot2(tree: Node[]): Node | null {
  let xor = 0
  for (const node of tree) {
    xor ^= node.val
    for (const child of node.children) {
      xor ^= child.val
    }
  }

  for (const node of tree) {
    if (node.val === xor) return node
  }

  throw new Error('invalid tree')
}

export {}
