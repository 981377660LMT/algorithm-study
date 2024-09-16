/**
 * O(n) 返回下一个字典序的排列.不含重复排列.
 */
function nextPermutation<T extends number | string>(nums: T[]): boolean {
  if (!nums.length) return false
  let left = nums.length - 1
  while (left > 0 && nums[left - 1] >= nums[left]) left--
  if (!left) return false
  const last = left - 1
  let right = nums.length - 1
  while (nums[right] <= nums[last]) right--
  const tmp = nums[last]
  nums[last] = nums[right]
  nums[right] = tmp
  reverseRange(nums, last + 1, nums.length - 1)
  return true
}

/**
 * O(n) 返回上一个字典序的排列.不含重复排列.
 */
function prePermutation<T extends number | string>(nums: T[]): boolean {
  if (!nums.length) return false
  let left = nums.length - 1
  while (left > 0 && nums[left - 1] <= nums[left]) left--
  if (!left) return false
  const last = left - 1 // 最后一个递减位置
  let right = nums.length - 1
  while (nums[right] >= nums[last]) right--
  const tmp = nums[last]
  nums[last] = nums[right]
  nums[right] = tmp
  reverseRange(nums, last + 1, nums.length - 1)
  return true
}

function reverseRange<T>(nums: T[], i: number, j: number): void {
  while (i < j) {
    const tmp = nums[i]
    nums[i] = nums[j]
    nums[j] = tmp
    i++
    j--
  }
}

export { nextPermutation, prePermutation }
