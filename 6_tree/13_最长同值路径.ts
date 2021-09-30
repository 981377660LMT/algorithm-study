interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
  val: 1,
  left: {
    val: 2,
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
  const helper = (parentValue: number, root: TreeNode | null): number => {
    if (!root) return 0
    // left包括了curNode在内和左边的长度
    const left = helper(root.val, root.left)
    const right = helper(root.val, root.right)
    res = Math.max(res, left + right)
    return root.val === parentValue ? Math.max(left, right) + 1 : 0
  }
  helper(root.val, root)
  return res
}

console.log(longestUnivaluePath(bt))

export {}
