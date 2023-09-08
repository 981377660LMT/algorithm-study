/* eslint-disable max-len */
/**
 * 返回arr的各个子集的元素和.
 * @param sorted 是否返回排序后的子集和.默认为false.
 */
function subsetSum(arr: ArrayLike<number>, sorted = false): number[] {
  return sorted ? subsetSumSorted(arr) : subsetSumUnsorted(arr)
}

/**
 * O(2^n)计算nums所有子集和.
 * 2^25(3e7) -> 200ms.
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
 * !比求出所有的子集的元素和再排序要快很多.
 * 2^25(3e7) -> 650ms.
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

/**
 * O(2^n)返回nums的各个子集的元素和的排序后的结果, 并且记录状态.
 * 2^25(3e7) -> 3s.
 */
function subsetSumSortedWithState(arr: ArrayLike<number>): [sum: number, state: number][] {
  let dp: [sum: number, state: number][] = [[0, 0]]
  for (let i = 0; i < arr.length; i++) {
    const ndp: [sum: number, state: number][] = Array(dp.length)
    for (let j = 0; j < dp.length; j++) {
      ndp[j] = [dp[j][0] + arr[i], dp[j][1] | (1 << i)]
    }
    dp = merge(dp, ndp)
  }

  return dp

  function merge(arr1: [sum: number, state: number][], arr2: [sum: number, state: number][]): [sum: number, state: number][] {
    const n1 = arr1.length
    const n2 = arr2.length
    const res: [sum: number, state: number][] = Array(n1 + n2)
    let i = 0
    let j = 0
    let k = 0
    while (i < n1 && j < n2) {
      if (arr1[i][0] < arr2[j][0]) res[k++] = arr1[i++]
      else res[k++] = arr2[j++]
    }
    while (i < n1) res[k++] = arr1[i++]
    while (j < n2) res[k++] = arr2[j++]
    return res
  }
}

/**
 * O(2^n) 返回arr的各个子集的元素和, 按照子集的元素个数分组.
 * ```ts
 * groupSubsetSumBySize([1, 2, 3]) // [[0], [3, 2, 1], [5, 4, 3], [6]]
 * ```
 */
function groupSubsetSumBySize(arr: ArrayLike<number>): number[][] {
  const dfs = (index: number, sum: number, count: number) => {
    if (index === n) {
      sumGroup[count].push(sum)
      return
    }
    dfs(index + 1, sum, count)
    dfs(index + 1, sum + arr[index], count + 1)
  }

  const n = arr.length
  const sumGroup: number[][] = Array(n + 1)
  for (let i = 0; i <= n; i++) sumGroup[i] = []
  dfs(0, 0, 0)
  return sumGroup
}

export { subsetSum, subsetSumSorted, subsetSumUnsorted, subsetSumSortedWithState, groupSubsetSumBySize }

if (require.main === module) {
  const nums = Array.from({ length: 25 }, (_, i) => i + 1)
  console.time('subsetSum')
  const a = subsetSum(nums)
  a.sort((a, b) => a - b)
  console.timeEnd('subsetSum')

  console.time('subsetSumSorted')
  const b = subsetSumSorted(nums)
  console.timeEnd('subsetSumSorted')

  console.time('subsetSumSortedWithState')
  const c = subsetSumSortedWithState(nums)
  console.timeEnd('subsetSumSortedWithState')
  console.log(subsetSumSortedWithState([1, 2, 3]))
}
