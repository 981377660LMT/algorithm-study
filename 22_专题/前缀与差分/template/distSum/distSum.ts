/**
 * **有序数组**所有点到`x=k`的距离之和.
 */
function distSum(sortedNums: ArrayLike<number>): (k: number) => number {
  const bisectRight = (nums: ArrayLike<number>, target: number) => {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = left + ((right - left) >>> 1)
      if (nums[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  const preSum = Array(sortedNums.length + 1)
  preSum[0] = 0
  for (let i = 0; i < sortedNums.length; i++) {
    preSum[i + 1] = preSum[i] + sortedNums[i]
  }

  return (k: number): number => {
    const pos = bisectRight(sortedNums, k)
    const leftSum = k * pos - preSum[pos]
    const rightSum = preSum[preSum.length - 1] - preSum[pos] - k * (sortedNums.length - pos)
    return leftSum + rightSum
  }
}

/**
 * **有序数组**中所有点对两两距离之和.一共有`n*(n-1)//2`对点对.
 */
function distSumOfAllPairs(sortedNums: ArrayLike<number>): number {
  let res = 0
  let preSum = 0
  for (let i = 0; i < sortedNums.length; i++) {
    res += sortedNums[i] * i - preSum
    preSum += sortedNums[i]
  }
  return res
}

export { distSum, distSumOfAllPairs }

if (require.main === module) {
  const sortedNums = [1, 2, 3, 4, 5, 6, 7, 8, 9]
  const Q = distSum(sortedNums)
  console.log(Q(5))
  console.log(distSumOfAllPairs(sortedNums))
}
