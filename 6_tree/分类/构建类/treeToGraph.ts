import { BinaryTree } from '../Tree'

// todo
// 思路是先用dfs序活得每个结点的id 再建图

/**
 *@description 二叉树转无向图
 */
function treeToGraph(root: BinaryTree) {
  const adjMap: Map<number, Set<number>> = new Map()
  const valueById: Map<number, number> = new Map()
  const idByNode: WeakMap<BinaryTree, number> = new WeakMap()
  let dfsId = 0

  dfsForId(root)
  dfsForGraph(root, null)

  return {
    adjMap,
    valueById,
    // idByNode,
  }

  function dfsForId(root: BinaryTree | null): void {
    if (!root) return
    dfsForId(root.left)
    dfsForId(root.right)
    idByNode.set(root, dfsId)
    dfsId++
  }

  function dfsForGraph(root: BinaryTree | null, parent: BinaryTree | null) {
    if (!root) return

    const rootId = idByNode.get(root)!
    valueById.set(rootId, root.val)

    if (parent) {
      const parentId = idByNode.get(parent)!
      !adjMap.has(parentId) && adjMap.set(parentId, new Set())
      !adjMap.has(rootId) && adjMap.set(rootId, new Set())
      adjMap.get(parentId)!.add(rootId)
      adjMap.get(rootId)!.add(parentId)
    }

    dfsForGraph(root.left, root)
    dfsForGraph(root.right, root)
  }
}

export { treeToGraph }
