import { BinaryTree } from '../Tree'
import { deserializeNode } from './297二叉树的序列化与反序列化'

const distanceK = function (root: BinaryTree | null) {
  const adjMap: Map<number, number[]> = new Map()

  const dfs = (root: BinaryTree | null, parent: BinaryTree | null) => {
    if (!root) return
    if (parent) {
      !adjMap.has(root.val) && adjMap.set(root.val, [])
      !adjMap.has(parent.val) && adjMap.set(parent.val, [])
      adjMap.get(root.val)!.push(parent.val)
      adjMap.get(parent.val)!.push(root.val)
    }

    dfs(root.left, root)
    dfs(root.right, root)
  }

  dfs(root, null)
  return adjMap
}

export default 1

console.log(distanceK(deserializeNode([1, 2, 3, 45])))
console.log(distanceK(deserializeNode([1])))

// 注意叶子节点root满足 adjMap.get(root)!.length===1
