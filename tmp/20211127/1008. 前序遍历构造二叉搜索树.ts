import { BinaryTree } from '../../6_tree/力扣加加/Tree'

// 序列里小于等于首元素的值递归进左子树，大于的递归进右子树
function bstFromPreorder(preorder: number[]): BinaryTree | null {
  if (preorder.length === 0) return null
  const root = new BinaryTree(preorder.shift()!)
  const left = preorder.filter(v => v < root.val)
  const right = preorder.filter(v => v > root.val)
  root.left = bstFromPreorder(left)
  root.right = bstFromPreorder(right)
  return root
}

console.log(bstFromPreorder([8, 5, 1, 7, 10, 12]))
