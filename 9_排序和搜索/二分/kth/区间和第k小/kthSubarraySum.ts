/**
 * 区间和的第 k 小.数组元素均为非负,1<=k<=n*(n+1)/2.
 */
function kthSubarraySum(arr: ArrayLike<number>, k: number): number {
  k-- // 从0开始
  const n = arr.length
  const preSum = Array(n + 1)
  preSum[0] = 0
  for (let i = 0; i < n; i++) {
    preSum[i + 1] = preSum[i] + arr[i]
  }

  // countNGT：和小于等于 target 的子数组个数 <= k.
  const check = (target: number): boolean => {
    let res = 0
    let left = 0
    let curSum = 0
    for (let right = 0; right < n; right++) {
      curSum += arr[right]
      while (left <= right && curSum > target) {
        curSum -= arr[left]
        left++
      }
      res += right - left + 1
    }
    return res <= k
  }

  let left = 0
  let right = preSum[n]
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }
  return left
}

export { kthSubarraySum }

if (require.main === module) {
  console.log(kthSubarraySum([2, 1, 3], 4))
}
