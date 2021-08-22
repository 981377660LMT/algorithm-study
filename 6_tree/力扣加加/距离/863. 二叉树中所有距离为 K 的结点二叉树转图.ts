import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @param {number} target
 * @param {number} k
 * @return {number[]}
 * @summary 二叉树转无向图 dfs建图 bfs寻找距离为k的节点
 */

const distanceK = function (root: BinaryTree, target: number, k: number): number[] {
  if (!root) return []
  const adjMap: Map<number, Set<number>> = new Map()

  const dfs = (root: BinaryTree) => {
    if (root.left) {
      adjMap.set(root.val, (adjMap.get(root.val) || new Set()).add(root.left.val))
      adjMap.set(root.left.val, (adjMap.get(root.left.val) || new Set()).add(root.val))
      dfs(root.left)
    }
    if (root.right) {
      adjMap.set(root.val, (adjMap.get(root.val) || new Set()).add(root.right.val))
      adjMap.set(root.right.val, (adjMap.get(root.right.val) || new Set()).add(root.val))
      dfs(root.right)
    }
  }
  dfs(root)

  const res: number[] = []
  const queue: number[] = [target]
  const visited = new Set<number>([target])

  while (queue.length) {
    k--
    const len = queue.length
    for (let i = 0; i < len; i++) {
      const head = queue.shift()!
      console.log(head)
      for (const next of adjMap.get(head)!) {
        if (visited.has(next)) continue
        visited.add(next)
        queue.push(next)
        k === 0 && res.push(next)
      }
    }
    if (k === 0) return res
  }
  // return res
  return res
}

console.dir(distanceK(deserializeNode([3, 5, 1, 6, 2, 0, 8, null, null, 7, 4])!, 5, 2), {
  depth: null,
})
// [ 1, 7, 4 ]
