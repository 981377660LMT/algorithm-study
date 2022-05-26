import { BinaryTree } from '../分类/Tree'

function sortedArrayToBST(nums: number[]): BinaryTree | null {
  if (!nums.length) return null
  const mid = nums.length >> 1
  const root = new BinaryTree(nums[mid])

  root.left = sortedArrayToBST(nums.slice(0, mid))
  root.right = sortedArrayToBST(nums.slice(mid + 1))

  return root
}

console.dir(sortedArrayToBST([-10, -3, 0, 5, 9]))
