/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * 将它原地重新排列成 nums[0] < nums[1] > nums[2] < nums[3]... 的顺序
 */
var wiggleSort = function (nums: number[]): void {
  nums.sort((a, b) => a - b)
  console.log(nums)
  // 小大小大... 小不能少于大
  let mid = nums.length >> 1
  mid += nums.length % 2 == 0 ? 0 : 1
  const small = nums.slice(0, mid)
  const big = nums.slice(mid)
  for (let i = 0; i < nums.length; i++) {
    if (i & 1) {
      nums[i] = big.pop()!
    } else {
      nums[i] = small.pop()!
    }
  }
}

const a = [1, 5, 1, 1, 6, 4]
wiggleSort(a)
console.log(a)

export default 1
