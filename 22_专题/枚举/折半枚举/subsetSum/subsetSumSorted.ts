/**
 * O(2^n) 返回nums的各个子集的元素和的排序后的结果.
 */
function subsetSumSorted(arr: ArrayLike<number>): number[] {
  const n = arr.length
  let dp = [0]
  for (let i = 0; i < n; i++) {
    const ndp = dp.slice()
    for (let j = 0; j < dp.length; j++) {
      ndp[j] += arr[i]
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

export { subsetSumSorted }
