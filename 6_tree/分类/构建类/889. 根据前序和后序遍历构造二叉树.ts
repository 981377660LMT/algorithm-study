import { BinaryTree } from '../Tree'

// 如果存在多个答案，您可以返回其中 任何 一个。
function constructFromPrePost(preorder: number[], postorder: number[]): BinaryTree | null {
  if (preorder.length === 0 || postorder.length === 0) return null
  const root = new BinaryTree(preorder[0])
  if (preorder.length === 1) return root
  // preorder第二个元素（如果存在的话）一定是左子树，即后序遍历(0, index + 1)全是左子树
  const index = postorder.indexOf(preorder[1])

  root.left = constructFromPrePost(preorder.slice(1, index + 2), postorder.slice(0, index + 1))
  root.right = constructFromPrePost(preorder.slice(index + 2), postorder.slice(index + 1, -1))

  return root
}

console.log(constructFromPrePost([1, 2, 4, 5, 3, 6, 7], [4, 5, 2, 6, 7, 3, 1]))
// 输出：[1,2,3,4,5,6,7]
