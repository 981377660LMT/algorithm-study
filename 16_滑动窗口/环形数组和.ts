// 找出环形数组中的连续子数组 和最接近某个数
const minSumDiff = (nums: number[]) => {
  const n = nums.length
  const total = nums.reduce((pre, cur) => pre + cur, 0)
  const half = total >> 1

  let res = Infinity
  let i = 0
  let j = 0
  let sum = 0

  // 左指针小于n 右指针小于2n
  while (i < n) {
    if (sum === half) return total & 1
    if (sum > half) sum -= nums[i++]
    else sum += nums[j++ % n]
    res = Math.min(res, Math.abs(total - 2 * sum))
  }

  return res
}

console.log(minSumDiff([1, 2, 3, 4]))
console.log(minSumDiff([10, 2, 8, 3]))
console.log(minSumDiff([0, 0, 0, 0, -1]))
console.log(minSumDiff([5, 8]))
