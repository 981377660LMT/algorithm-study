/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * 将它原地重新排列成 nums[0] < nums[1] > nums[2] < nums[3]... 的顺序
 */
function wiggleSort(nums: number[]): void {
  const n = nums.length
  nums.sort((a, b) => a - b)

  // 小大小大... 小不能少于大
  let half = (n + 1) >> 1
  const small = nums.slice(0, half)
  const big = nums.slice(half)

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
