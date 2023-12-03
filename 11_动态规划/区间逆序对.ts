/**
 * 区间逆序对.
 * 时空复杂度 O(n^2) 预处理, O(1) 查询.
 */
function rangeInv(arr: ArrayLike<number>): (start: number, end: number) => number {
  const n = arr.length
  const dp = new Uint16Array((n + 1) * (n + 1))
  for (let left = n; left >= 0; left--) {
    for (let right = left; right <= n; right++) {
      if (right - left <= 1) {
        continue
      }
      dp[left * (n + 1) + right] =
        dp[left * (n + 1) + right - 1] +
        dp[(left + 1) * (n + 1) + right] -
        dp[(left + 1) * (n + 1) + right - 1] +
        +(arr[left] > arr[right - 1])
    }
  }

  const cal = (start: number, end: number): number => {
    if (start >= end) {
      return 0
    }
    return dp[start * (n + 1) + end]
  }

  return cal
}

export { rangeInv }

if (require.main === module) {
  const randNums = (n: number) => Array.from({ length: n }, () => Math.floor(Math.random() * n))
  const arr = randNums(10)
  const bruteForce = (arr: number[], start: number, end: number): number => {
    let res = 0
    for (let i = start; i < end; i++) {
      for (let j = i + 1; j < end; j++) {
        res += +(arr[i] > arr[j])
      }
    }
    return res
  }
  const cal = rangeInv(arr)
  for (let i = 0; i < arr.length; i++) {
    for (let j = i + 1; j <= arr.length; j++) {
      if (cal(i, j) !== bruteForce(arr, i, j)) {
        console.error('error', arr, i, j)
        throw new Error()
      }
    }
  }

  console.info('ok')
}
