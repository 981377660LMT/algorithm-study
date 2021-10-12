import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {BinaryTree}
 * @description
 * 对位于 (row, col) 的每个结点而言，
 * 其左右子结点分别位于 (row + 1, col - 1) 和 (row + 1, col + 1) 。
 * 树的根结点位于 (0, 0) 。
 * @description
 */
const verticalTraversal = (root: BinaryTree | null): number[][] => {
  if (!root) return []

  const tmp: [number, number, number][] = []
  const inOrder = (root: BinaryTree, x: number, y: number) => {
    tmp.push([x, y, root.val])
    root.left && inOrder(root.left, x - 1, y + 1)
    root.right && inOrder(root.right, x + 1, y + 1)
  }
  inOrder(root, 0, 0)
  tmp.sort((a, b) => a[0] - b[0] || a[1] - b[1] || a[2] - b[2])

  // x值为key的map
  const map = new Map<number, number[]>()
  for (const item of tmp) {
    const key = item[0]
    const val = item[2]
    !map.has(key) && map.set(key, [])
    map.get(key)!.push(val)
  }

  return [...map.values()]
}

console.dir(verticalTraversal(deserializeNode([3, 9, 20, null, null, 15, 7])!), {
  depth: null,
})
// 输出：[[9],[3,15],[20],[7]]
console.dir(
  verticalTraversal(
    deserializeNode([1, 3, 2, null, 4, 5, null, 6, 13, 7, 8, 14, 12, 10, null, null, 11, 9])!
  ),
  {
    depth: null,
  }
)
export {}
