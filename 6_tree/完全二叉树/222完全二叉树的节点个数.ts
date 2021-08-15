// 二叉树
interface BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
}

const bt: BinaryTree = {
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

/**
 * 如果当前节点的左右子树高度相同，那么左子树是一个满二叉树，右子树是一个完全二叉树。
 * 否则（左边的高度大于右边），那么左子树是一个完全二叉树，右子树是一个满二叉树。
 * 满二叉树的结点数是2^h-1
 * 时间复杂度：$O(logN * log N)$
 * 空间复杂度：$O(logN)$
 */
const countNodes = function (root: BinaryTree | null): number {
  if (root == null) return 0

  // 计算完全二叉树的高度，可以只递归左树
  const getDepth = (root: BinaryTree | null): number => {
    if (!root) return 0
    return getDepth(root.left) + 1
  }

  const leftDepth = getDepth(root.left)
  const rightDepth = getDepth(root.right)
  if (leftDepth === rightDepth) {
    return 2 ** leftDepth + countNodes(root.right)
  } else {
    return 2 ** rightDepth + countNodes(root.left)
  }
}

console.log(countNodes(bt))
export default 1
