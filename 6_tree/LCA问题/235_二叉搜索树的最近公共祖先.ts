interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
  val: 6,
  left: {
    val: 2,
    left: {
      val: 0,
      left: null,
      right: null,
    },
    right: {
      val: 4,
      left: {
        val: 3,
        left: null,
        right: null,
      },
      right: {
        val: 5,
        left: null,
        right: null,
      },
    },
  },
  right: {
    val: 8,
    left: {
      val: 7,
      left: null,
      right: null,
    },
    right: {
      val: 9,
      left: null,
      right: null,
    },
  },
}

const p: TreeNode = {
  val: 3,
  left: null,
  right: null,
}

const q: TreeNode = {
  val: 7,
  left: null,
  right: null,
}
// 给定一个二叉搜索树, 找到该树中两个指定节点的最近公共祖先。
// 注意二叉搜索树的性质
// 迭代法: O(h) O(1)
const lowestCommonAncestor = (root: TreeNode, p: TreeNode, q: TreeNode) => {
  while (root) {
    if (root.val > p.val && root.val > q.val) {
      root = root.left!
    } else if (root.val < p.val && root.val < q.val) {
      root = root.right!
    } else {
      break
    }
  }
  return root
}

// 递归法 O(h) O(h)
const lowestCommonAncestor2 = (
  root: TreeNode | null,
  p: TreeNode,
  q: TreeNode
): TreeNode | null => {
  if (!root) return root
  if (root.val < p.val && root.val < q.val) {
    return lowestCommonAncestor2(root.right, p, q)
  }
  if (root.val > p.val && root.val > q.val) {
    return lowestCommonAncestor2(root.left, p, q)
  }
  return root
}
console.dir(lowestCommonAncestor2(bt, p, q), { depth: null })

export {}
