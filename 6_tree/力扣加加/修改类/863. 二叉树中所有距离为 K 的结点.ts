import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'
import { AddPropRecursively } from './117. 填充每个节点的下一个右侧节点指针'

type NodeWithParent = AddPropRecursively<BinaryTree, { parent?: BinaryTree | null }>
/**
 * @param {BinaryTree} root
 * @param {number} target
 * @param {number} k
 * @return {number[]}
 * @description 返回到目标结点 target 距离为 K 的所有结点的值的列表。 答案可以以任何顺序返回。
 * @description 为了标记节点是否访问过，节点需要添加isVisited 但是为了方便让所有节点val不同
 * @description 先dfs找到这个点 并且增加parent节点(当成图) 然后从这个target出发bfs
 */

const distanceK = function (root: NodeWithParent, target: number, k: number): number[] {
  if (!root) return []
  const res: number[] = []

  const findTargetNode = (
    root: NodeWithParent | null,
    parent: NodeWithParent | null
  ): NodeWithParent | null => {
    if (!root) return null
    root.parent = parent
    if (root.val === target) return root
    return findTargetNode(root.left, root) || findTargetNode(root.right, root)
  }

  const targetNode = findTargetNode(root, null)
  if (!targetNode) return []

  const bfs = (root: NodeWithParent | null, steps: number, visited: Set<number>): void => {
    if (!root) return
    if (steps === k) {
      res.push(root.val)
      return
    }

    const next = [root.parent, root.left, root.right]
    for (const nextNode of next) {
      if (nextNode && !visited.has(nextNode.val)) {
        visited.add(nextNode.val)
        bfs(nextNode, steps + 1, visited)
      }
    }
  }
  bfs(targetNode, 0, new Set([targetNode.val]))

  return res
}

console.dir(distanceK(deserializeNode([3, 5, 1, 6, 2, 0, 8, null, null, 7, 4])!, 5, 2), { depth: null })
// [ 1, 7, 4 ]
