/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * 算法需要将给定数字序列重新排列成字典序中下一个更大的排列。
 * 如果不存在下一个更大的排列，则将数字重新排列成最小的排列（即升序排列）。
 * @summary 倒序遍历  尽早交换  发现大的就叫唤
 */
const nextPermutation = function (nums: number[]) {
  let mono = true
  const partialSort = <T = number>(
    nums: T[],
    start: number,
    end: number,
    compareFn?: (a: T, b: T) => number
  ) => {
    const preSorted = nums.slice(0, start)
    const postSorted = nums.slice(end)
    const sorted = nums.slice(start, end).sort(compareFn)
    nums.length = 0
    nums.push.apply(nums, preSorted.concat(sorted, postSorted))
    return nums
  }
  loop: for (let i = nums.length - 1; i >= 0; i--) {
    for (let j = nums.length - 1; j > i; j--) {
      // 后面大于前面
      if (nums[j] > nums[i]) {
        console.log(i, j)
        // 交换玩排序
        ;[nums[i], nums[j]] = [nums[j], nums[i]]
        partialSort(nums, i + 1, nums.length, (a, b) => a - b)
        mono = false
        break loop
      }
    }
  }

  mono && nums.reverse()
  return nums
}

// console.log(nextPermutation([1, 2, 3]))
// // 输出：[1,3,2]
// console.log(nextPermutation([3, 2, 1]))
// // 输出：[1,2,3]
// console.log(nextPermutation([1, 2, 4, 3]))
// 1324
console.log(nextPermutation([2, 3, 1]))
export default 1
