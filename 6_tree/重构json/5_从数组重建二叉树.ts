class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val: number) {
    this.val = val
    this.left = null
    this.right = null
  }
}

// bfs
const arrayToTree = (arr: (number | null)[]) => {
  const toNode = (item: number | null) => (item == null ? null : new TreeNode(item))

  const root = toNode(arr.shift()!)!
  const queue: TreeNode[] = [root]
  while (queue.length) {
    const head = queue.shift()!
    head.left = toNode(arr.shift()!)
    head.right = toNode(arr.shift()!)
    head.left && queue.push(head.left)
    head.right && queue.push(head.right)
  }

  return root
}

console.dir(arrayToTree([1, 2, 3, null, 5, 6, 7]), { depth: null })
export {}
