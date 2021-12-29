import { BinaryTree } from '../../6_tree/力扣加加/Tree'

/**
 *
 * @param root 找出 BST 中的所有众数(可能不止一个)
 */
function findMode(root: BinaryTree | null): number[] {
  if (!root) return []
  let pre: number | null = null
  let curCount = 0
  let maxCount = 0
  let res: number[] = []

  const inorder = (root: BinaryTree | null) => {
    if (!root) return
    inorder(root.left)

    pre === root.val ? curCount++ : (curCount = 1)
    if (curCount === maxCount) {
      res.push(root.val)
    } else if (curCount > maxCount) {
      res = [root.val]
      maxCount = curCount
    }
    pre = root.val

    inorder(root.right)
  }

  inorder(root)
  return res
}

// 1. 判断与之前是否相等
// 2. 判断curCount是否等于maxCount 添加候选
// 3. 判断curCount是否大于maxCount 重置候选
