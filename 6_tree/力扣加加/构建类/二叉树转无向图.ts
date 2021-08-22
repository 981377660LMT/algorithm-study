import { BinaryTree } from '../Tree'

const distanceK = function (root: BinaryTree) {
  if (!root) return
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
}

export default 1
