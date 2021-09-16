// 二叉树中寻找两个节点的最近公共祖先(LCA)
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
      val: 3,
      left: null,
      right: null,
    },
    right: {
      val: 4,
      left: {
        val: 5,
        left: null,
        right: null,
      },
      right: {
        val: -9,
        left: null,
        right: null,
      },
    },
  },
  right: {
    val: 9,
    left: {
      val: 7,
      left: null,
      right: null,
    },
    right: {
      val: 8,
      left: null,
      right: null,
    },
  },
}

const p: TreeNode = {
  val: -9,
  left: null,
  right: null,
}

const q: TreeNode = {
  val: 9,
  left: null,
  right: null,
}

/**
 * @description 若 root 为 p，q 的公共祖先，有三种情况

                p, q 分别在 root 的左右子树；
                p === root 切 q 在 p 的左右子树；
                q === root 且 p 在 q 的左右子树。

 * @description 在root之下递归寻找p与q的最近祖先
 *              递归返回的内容是包含目标节点的子节点
 *              如果某一个子节点同时包含两个目标节点，则返回该节点
 *              如果某一个子节点只包含一个目标节点，则返回该目标节点
 *              否则返回为null
 */
const lowestCommonAncestor = (root: TreeNode | null, p: TreeNode, q: TreeNode): TreeNode | null => {
  if (root === null || root.val === p.val || root.val === q.val) return root
  // 左子树中存在p或q
  const detectLeft = lowestCommonAncestor(root.left, p, q)
  // 右子树中存在p或q
  const detectRight = lowestCommonAncestor(root.right, p, q)

  // 排除有一个不存在的情况
  if (detectLeft === null) return detectRight
  if (detectRight === null) return detectLeft
  // p, q 分别在 root 的左右子树；
  return root
}

console.dir(lowestCommonAncestor(bt, p, q), { depth: null })

export {}
