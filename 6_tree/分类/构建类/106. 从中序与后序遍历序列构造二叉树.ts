import { BinaryTree } from '../Tree'

/**
 * @param {number[]} inorder
 * @param {number[]} postorder
 * @return {BinaryTree}
 */
function buildTree(inorder: number[], postorder: number[]): BinaryTree | null {
  if (!inorder.length || !postorder.length) return null

  const rootValue = postorder[postorder.length - 1]
  const root = new BinaryTree(rootValue)
  const index = inorder.indexOf(rootValue)
  root.left = buildTree(inorder.slice(0, index), postorder.slice(0, index))
  root.right = buildTree(inorder.slice(index + 1), postorder.slice(index, -1))

  return root
}

console.log(buildTree([9, 3, 15, 20, 7], [9, 15, 7, 20, 3]))
