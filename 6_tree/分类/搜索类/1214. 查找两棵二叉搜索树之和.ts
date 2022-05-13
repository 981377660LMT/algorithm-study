import { BinaryTree } from '../Tree'

/**
 *
 * @param root1
 * @param root2
 * @param target
 * 请你从两棵树中各找出一个节点，使得这两个节点的值之和等于目标值 Target
 */
function twoSumBSTs(root1: BinaryTree | null, root2: BinaryTree | null, target: number): boolean {
  if (!root1 || !root2) return false
  const nums1: number[] = []
  const nums2: number[] = []
  inorder(root1, nums1)
  inorder(root2, nums2)
  return findTarget(nums1, nums2, target)

  // 两个有序数组找和为target
  function findTarget(nums1: number[], nums2: number[], target: number) {
    let i = 0
    let j = nums2.length - 1

    while (i < nums1.length && j >= 0) {
      const sum = nums1[i] + nums2[j]
      if (sum === target) return true
      else if (sum < target) i++
      else j--
    }

    return false
  }

  function inorder(root: BinaryTree | null, nums: number[]) {
    if (!root) return
    inorder(root.left, nums)
    nums.push(root.val)
    inorder(root.right, nums)
  }
}
