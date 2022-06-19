/**
 *
 * @param nums digit 数组
 * @description 返回值第二个参数ok
 */
function nextPermutation<T>(nums: T[], inplace = false): [res: T[], ok: boolean] {
  if (nums.length === 0) return [[], false]
  if (!inplace) nums = nums.slice()

  let left = nums.length - 1
  while (left > 0 && nums[left - 1] >= nums[left]) left--
  if (left === 0) return [[], false]
  const last = left - 1 // 最后一个递增位置

  let right = nums.length - 1
  while (nums[right] <= nums[last]) right--
  ;[nums[last], nums[right]] = [nums[right], nums[last]] // 找到最小的可交换的right，交换这两个数

  reverseRange(nums, last + 1, nums.length - 1)
  return [nums, true]
}

/**
 *
 * @param nums digit 数组
 * @description 返回值第二个参数带ok
 */
function prePermutation<T>(nums: T[], inplace = false): [res: T[], ok: boolean] {
  if (nums.length === 0) return [[], false]
  if (!inplace) nums = nums.slice()

  let left = nums.length - 1
  while (left > 0 && nums[left - 1] <= nums[left]) left--
  if (left === 0) return [[], false]
  const last = left - 1 // 最后一个递减位置

  let right = nums.length - 1
  while (nums[right] >= nums[last]) right--
  ;[nums[last], nums[right]] = [nums[right], nums[last]] // 找到最小的可交换的right，交换这两个数

  reverseRange(nums, last + 1, nums.length - 1)
  return [nums, true]
}

function reverseRange<T>(nums: T[], i: number, j: number) {
  while (i < j) {
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
    i++
    j--
  }
}

export { nextPermutation, prePermutation }
