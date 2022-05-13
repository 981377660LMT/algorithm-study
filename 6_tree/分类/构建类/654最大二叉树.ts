import { BinaryTree } from '../Tree'

/**
 * @param {number[]} nums
 * @return {BinaryTree}
 * @description 
 * 二叉树的根是数组 nums 中的最大元素。
   左子树是通过数组中 最大值左边部分 递归构造出的最大二叉树。
   右子树是通过数组中 最大值右边部分 递归构造出的最大二叉树。
 */
const constructMaximumBinaryTree = (nums: number[]): BinaryTree | null => {
  /**
   *
   * @param arr
   * @param left
   * @param right
   * @returns  以这一段构建树
   */
  const helper = (arr: number[], left: number, right: number): BinaryTree | null => {
    if (left > right) return null
    let maxIndex = left
    for (let index = left; index <= right; index++) {
      if (arr[index] > arr[maxIndex]) maxIndex = index
    }
    const root = new BinaryTree(arr[maxIndex])
    root.left = helper(arr, left, maxIndex - 1)
    root.right = helper(arr, maxIndex + 1, right)
    return root
  }
  return helper(nums, 0, nums.length - 1)
}

console.dir(constructMaximumBinaryTree([3, 2, 1, 6, 0, 5]), { depth: null })
