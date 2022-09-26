// 两个节点之间的路径长度由它们之间的边数表示。
// !给定一个二叉树，找到最长的路径，这个路径中的每个节点具有相同值。
// !这条路径可以经过也可以不经过根节点。
// 思路:从叶子节点开始遍历(但是找不到父节点)

function longestUnivaluePath(root: TreeNode | null): number {
  if (!root) return 0

  let res = 0
  dfs(root, null)
  return res

  // !非自顶向下的需要在每个位置更新答案
  function dfs(cur: TreeNode | null, pre: TreeNode | null): number {
    if (!cur) return 0
    const left = dfs(cur.left, cur)
    const right = dfs(cur.right, cur)
    res = Math.max(res, left + right) // !经过当前节点的最长路径
    return cur.val === pre?.val ? Math.max(left, right) + 1 : 0
  }
}

export {}
