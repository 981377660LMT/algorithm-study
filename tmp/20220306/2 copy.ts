function minimalKSum(nums: number[], k: number): number {
  nums = [...new Set(nums)].sort((a, b) => a - b)
  const mex = findMex(nums, k)
  const index = bisectLeft(nums, mex)
  const allSum = ((1 + mex) * mex) / 2
  return allSum - nums.slice(0, index).reduce((a, b) => a + b, 0)

  function findMex(nums: number[], k: number): number {
    let [left, right] = [0, nums.length - 1]
    while (left <= right) {
      const mid = (left + right) >> 1
      const diff = nums[mid] - (mid + 1)
      if (diff >= k) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }

    return left + k
  }
}

function bisectLeft(arr: number[], target: number): number {
  if (arr.length === 0) return 0

  let l = 0
  let r = arr.length - 1
  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    const midElement = arr[mid]
    // 尽量右移
    if (midElement < target) {
      // mid 根本就不是答案，直接更新 l = mid + 1，从而将 mid 从解空间排除
      l = mid + 1
    } else if (midElement >= target) {
      // midElement >= target :将 mid 从解空间排除，继续看看有没有更好的
      r = mid - 1
    }
  }

  return l
}

console.log(
  minimalKSum(
    [
      96, 44, 99, 25, 61, 84, 88, 18, 19, 33, 60, 86, 52, 19, 32, 47, 35, 50, 94, 17, 29, 98, 22,
      21, 72, 100, 40, 84,
    ],
    35
  )
)
console.log(minimalKSum([1, 4, 25, 10, 25], 2))
console.log(Number.MAX_SAFE_INTEGER, 2 ** 53 - 1)
console.log(minimalKSum([1000000000], 1000000000))
// 500000000500000000

console.log(500000000500000000)
