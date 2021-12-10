class Node {
  val: number
  left: Node | null
  right: Node | null
  random: Node | null
  constructor(val?: number, left?: Node | null, right?: Node | null, random?: Node | null) {
    this.val = val === undefined ? 0 : val
    this.left = left === undefined ? null : left
    this.right = right === undefined ? null : right
    this.random = random === undefined ? null : random
  }
}

// 给你一个二叉树，树中每个节点都含有一个附加的随机指针，
// 该指针可以指向树中的任何节点或者指向空（null）。
// 请返回该树的 深拷贝 。

// 总结：
// 两次dfs，第一次复制树并建立两个树的节点地址map，第二次填充random
function copyRandomBinaryTree(root: Node | null): NodeCopy | null {
  if (!root) return null
  const record = new WeakMap<Node, Node | undefined>()
  setupNode(root)
  setupPointer(root)
  return record.get(root)!

  function setupNode(root: Node | null) {
    if (!root) return
    //@ts-ignore
    record.set(root, new NodeCopy(root.val))
    setupNode(root.left)
    setupNode(root.right)
  }

  function setupPointer(root: Node | null) {
    if (!root) return
    const newRoot = record.get(root)!
    root.left && (newRoot.left = record.get(root.left)!)
    root.right && (newRoot.right = record.get(root.right)!)
    root.random && (newRoot.random = record.get(root.random)!)
    setupPointer(root.left)
    setupPointer(root.right)
  }
}

export {}
