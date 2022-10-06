import { BinaryTree } from '../分类/Tree'

/**
 * 如果当前节点的左右子树高度相同，那么左子树是一个满二叉树，右子树是一个完全二叉树。
 * 否则（左边的高度大于右边），那么左子树是一个完全二叉树，右子树是一个满二叉树。
 * 满二叉树的结点数是2^h-1
 * 时间复杂度：$O(logN * log N)$
 * 空间复杂度：$O(logN)$
 */
function countNodes(root: BinaryTree | null): number {
  if (root == null) {
    return 0
  }

  const leftDepth = getDepth(root.left)
  const rightDepth = getDepth(root.right)
  if (leftDepth === rightDepth) {
    // 左子树一定是满二叉树
    return 2 ** leftDepth + countNodes(root.right)
  }

  // 此时最后一层不满，但倒数第二层已经满了，可以直接得到右子树的节点个数
  return 2 ** rightDepth + countNodes(root.left)

  // !计算完全二叉树的高度
  // const getDepth = (root: BinaryTree | null): number => {
  //   if (!root) return 0
  //   return getDepth(root.left) + 1
  // }
  function getDepth(root: BinaryTree | null): number {
    let res = 0
    while (root) {
      res++
      root = root.left
    }
    return res
  }
}

export default 1

// 由于完全二叉树的性质，其子树一定有一棵是满的，所以一定会触发 leftDepth == rightDepth，
// 只消耗 O(logN) 的复杂度而不会继续递归。
// 综上，算法的递归深度就是树的高度 O(logN)，每次递归所花费的时间就是 while 循环，
// 需要 O(logN)，所以总体的时间复杂度是 O(logN*logN)。
