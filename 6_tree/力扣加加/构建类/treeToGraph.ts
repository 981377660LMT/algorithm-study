import { BinaryTree } from '../Tree'

/**
 *
 * @param root 各个节点的值不同
 * @returns
 */
const treeToGraph = function (root: BinaryTree) {
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

  return adjMap
}

export { treeToGraph }
