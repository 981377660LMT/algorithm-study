import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

class Node {
  constructor(
    public val: number = 0,
    public left: Node | null = null,
    public right: Node | null = null
  ) {}
}

/**
 * @param {Node} root
 * @return {Node}
 * @description 将这个二叉搜索树(BST)转化为双向循环链表
 */
const treeToDoublyList = function (root: Node): Node {
  const queue: Node[] = []
  const inorder = (root: Node) => {
    root.left && inorder(root.left)
    queue.push(root)
    root.right && inorder(root.right)
  }
  inorder(root)

  const len = queue.length
  for (let i = 0; i < len; i++) {
    queue[i].right = queue[(i + 1) % len]
    queue[i].left = queue[(i - 1 + len) % len]
  }
  console.dir(queue, { depth: null })
  return root
}

console.log(treeToDoublyList(deserializeNode([4, 2, 5, 1, 3, null, null])!))
