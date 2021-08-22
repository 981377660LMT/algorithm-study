import { BinaryTree } from '../Tree'

type AddPropRecursively<O extends object | null, P> = {
  [K in keyof O]: O[K] extends O | null ? AddPropRecursively<O[K], P> | null : O[K]
} &
  P

type NodeWithNext = AddPropRecursively<BinaryTree, { next: BinaryTree | null }>

const bt: NodeWithNext = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
      next: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
      next: null,
    },
    next: null,
  },
  right: {
    val: 3,
    left: null,
    right: {
      val: 7,
      left: null,
      right: null,
      next: null,
    },
    next: null,
  },
  next: null,
}

/**
 * @param {NodeWithNext} root
 * @return {NodeWithNext}
 * @description bfs记录上一个即可
 */
const connect = (root: NodeWithNext): NodeWithNext | null => {
  if (!root) return null

  const dfs = (root: NodeWithNext | null) => {
    if (!root) return
    dfs(root.left)
    dfs(root.right)
    // 执行完dfs后，这个子树就完成连接了，此时只需要将两个子树连接起来。
    // 把左子树每一层最右边的节点，连接到右子树每一层最左边的节点。
    let left: NodeWithNext | null = root.left
    let right: NodeWithNext | null = root.right
    while (left && right) {
      left.next = right
      left = left.right
      right = right.left
    }
  }
  dfs(root)

  return root
}

console.dir(connect(bt), { depth: null })
// 输出：[1,#,2,3,#,4,5,7,#]
export { AddPropRecursively }
