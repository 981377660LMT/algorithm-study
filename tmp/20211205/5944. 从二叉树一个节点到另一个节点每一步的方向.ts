import { BinaryTree } from '../../6_tree/力扣加加/Tree'
import { deserializeNode } from '../../6_tree/力扣加加/构建类/297.二叉树的序列化与反序列化'
type Direction = 'L' | 'R' | 'U'

// 请找到从节点 s 到节点 t 的 最短路径 ，并以字符串的形式返回每一步的方向。每一步用 大写 字母 'L' ，'R' 和 'U' 分别表示一种方向：
function getDirections(root: BinaryTree | null, startValue: number, destValue: number): string {
  const adjMap = new Map<number, [next: number, direction: Direction][]>()

  buildGraph(root, -1)

  const queue: [cur: number, path: string][] = [[startValue, '']]
  const visited = new Set<number>()

  while (queue.length > 0) {
    const [cur, path] = queue.shift()!
    if (visited.has(cur)) continue
    visited.add(cur)
    if (cur === destValue) return path
    for (const [next, direction] of adjMap.get(cur) || []) {
      queue.push([next, path + direction])
    }
  }

  return ''

  function buildGraph(root: BinaryTree | null, parent: number): void {
    if (!root) return

    !adjMap.has(root.val) && adjMap.set(root.val, [])
    adjMap.get(root.val)!.push([parent, 'U'])

    if (root.left) {
      adjMap.get(root.val)!.push([root.left.val, 'L'])
      buildGraph(root.left, root.val)
    }

    if (root.right) {
      adjMap.get(root.val)!.push([root.right.val, 'R'])
      buildGraph(root.right, root.val)
    }
  }
}

console.log(getDirections(deserializeNode([5, 1, 2, 3, null, 6, 4]), 3, 6))
