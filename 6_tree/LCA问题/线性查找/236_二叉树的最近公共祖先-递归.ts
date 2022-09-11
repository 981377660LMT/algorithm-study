// 时空复杂度O(n)
function lowestCommonAncestor(root: TreeNode | null, p: TreeNode, q: TreeNode): TreeNode | null {
  if (root == null || root.val === p.val || root.val === q.val) return root

  const leftRes = lowestCommonAncestor(root.left, p, q)
  const rightRes = lowestCommonAncestor(root.right, p, q)

  // 排除有一个不存在的情况
  if (!leftRes) return rightRes
  if (!rightRes) return leftRes
  return root
}

export {}
