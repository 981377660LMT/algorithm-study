class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
    this.val = val === undefined ? 0 : val
    this.left = left === undefined ? null : left
    this.right = right === undefined ? null : right
  }
}

// 572. 另一棵树的子树
function isSubtree(root: TreeNode | null, subRoot: TreeNode | null): boolean {
  if (!root && !subRoot) return true
  if (!root || !subRoot) return false

  return (
    isSameTree(root, subRoot) || isSubtree(root.left, subRoot) || isSubtree(root.right, subRoot)
  )
}

function isSameTree(root1: TreeNode | null, root2: TreeNode | null): boolean {
  if (!root1 && !root2) return true
  if (!root1 || !root2) return false

  return (
    root1.val === root2.val &&
    isSameTree(root1.left, root2.left) &&
    isSameTree(root1.right, root2.right)
  )
}
