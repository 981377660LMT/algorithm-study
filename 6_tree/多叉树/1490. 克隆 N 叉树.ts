class Node {
  val: number
  children: Node[]

  constructor(val = 0, children: Node[] = []) {
    this.val = val
    this.children = children
  }
}

function cloneTree(root: Node): Node {
  if (!root) return root
  const newNode = new Node(root.val)

  for (const child of root.children) {
    newNode.children.push(cloneTree(child))
  }

  return newNode
}

export {}
