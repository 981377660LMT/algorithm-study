interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
  val: 1,
  left: {
    val: 1,
    left: {
      val: 4,
      left: null,
      right: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
    },
  },
  right: {
    val: 3,
    left: {
      val: 6,
      left: null,
      right: null,
    },
    right: {
      val: 7,
      left: null,
      right: null,
    },
  },
}

const bt2: TreeNode = {
  val: 1,
  left: { val: 2, left: null, right: null },
  right: null,
}
// 两个节点之间的路径长度由它们之间的边数表示。
// 给定一个二叉树，找到最长的路径，这个路径中的每个节点具有相同值。
// 这条路径可以经过也可以不经过根节点。
// 思路:从叶子节点开始遍历(但是找不到父节点)

// 正确解法:没有思路的时候用递归(左子树右子树递归+元信息)
const longestUnivaluePath = (root: TreeNode) => {
  if (!root) return 0
  let res = 0
  const dfs = (root: TreeNode | null): number => {
    if (!root) return 0

    let left = dfs(root.left)
    let right = dfs(root.right)
    if (root.left && root.val === root.left.val) left++
    else left = 0
    if (root.right && root.val === root.right.val) right++
    else right = 0

    res = Math.max(res, left + right)
    return Math.max(left, right)
  }
  dfs(root)
  return res
}

console.log(longestUnivaluePath(bt))

export {}
