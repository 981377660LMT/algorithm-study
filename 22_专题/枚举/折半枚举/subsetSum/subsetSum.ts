/**
 * 返回arr的各个子集的元素和.
 * @param sorted 是否返回排序后的子集和.默认为false.
 */
function subsetSum(arr: ArrayLike<number>, sorted = false): number[] {
  return sorted ? subsetSumSorted(arr) : subsetSumUnsorted(arr)
}

/**
 * O(2^n)计算nums所有子集和.
 * 2^25(3e7) -> 200ms
 */
function subsetSumUnsorted(arr: ArrayLike<number>): number[] {
  const n = arr.length
  const res = Array<number>(1 << n).fill(0)
  for (let i = 0; i < n; i++) {
    for (let pre = 0; pre < 1 << i; pre++) {
      res[pre | (1 << i)] = res[pre] + arr[i]
    }
  }
  return res
}

/**
 * O(2^n)返回nums的各个子集的元素和的排序后的结果.
 * !比求出所有的子集的元素和再排序要快很多
 * 2^25(3e7) -> 650ms
 */
function subsetSumSorted(arr: ArrayLike<number>): number[] {
  const n = arr.length
  let dp = [0]
  for (let i = 0; i < n; i++) {
    const ndp = Array<number>(dp.length)
    for (let j = 0; j < dp.length; j++) {
      ndp[j] = dp[j] + arr[i]
    }
    dp = merge(dp, ndp)
  }
  return dp

  function merge(arr1: ArrayLike<number>, nums2: ArrayLike<number>): number[] {
    const n1 = arr1.length
    const n2 = nums2.length
    const res = Array<number>(n1 + n2)
    let i = 0
    let j = 0
    let k = 0
    while (i < n1 && j < n2) {
      if (arr1[i] < nums2[j]) res[k++] = arr1[i++]
      else res[k++] = nums2[j++]
    }
    while (i < n1) res[k++] = arr1[i++]
    while (j < n2) res[k++] = nums2[j++]
    return res
  }
}

export { subsetSum, subsetSumSorted, subsetSumUnsorted }

if (require.main === module) {
  const nums = Array.from({ length: 25 }, (_, i) => i + 1)
  console.time('subsetSum')
  const a = subsetSum(nums)
  console.timeEnd('subsetSum')

  console.time('subsetSumSorted')
  const b = subsetSumSorted(nums)
  console.timeEnd('subsetSumSorted')
}
