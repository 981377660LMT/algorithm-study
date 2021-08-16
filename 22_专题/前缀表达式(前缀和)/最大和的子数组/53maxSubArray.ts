// 将点的和画成折线图，找出图上差最大的的那两点即可
const maxSubArray = (nums: number[]) => {
  if (nums.length <= 1) return nums

  // S0到Sn
  const pre = Array(nums.length + 1).fill(0)
  let curSum = 0
  let min = 0
  let maxDiff = 0
  let minIndex = 0
  let maxIndex = 0

  for (let index = 0; index < nums.length; index++) {
    curSum += nums[index]
    pre[index + 1] = curSum
    if (curSum - min > maxDiff) {
      maxDiff = curSum - min
      maxIndex = index
    }
    if (curSum < min) {
      min = curSum
      minIndex = index
    }
  }
  console.log(min, minIndex, maxIndex)
  return [maxDiff, nums.slice(minIndex, maxIndex + 1)]
}

console.log(maxSubArray([-2, 1, -3, 4, -1, 2, 1, -5, 4]))
console.log(maxSubArray([1, -2, 3, 10, -4, 7, 2, -5]))
console.log(maxSubArray([1, 3, -2, 4, -5]))

export {}
