const qs = (arr: number[]) => qucikSort(arr, 0, arr.length - 1)

/**
 *
 * @param arr
 * @returns
 * @description
 *（1）在数据集之中，选择一个元素作为"基准"（pivot）。
 *（2）所有小于"基准"的元素，都移到"基准"的左边；所有大于"基准"的元素，都移到"基准"的右边。
 *（3）对"基准"左边和右边的两个子集，不断重复第一步和第二步，直到所有子集只剩下一个元素为止。
 * @description 存在的问题:完全有序的数组 退化成 时间复杂度O(n^2) 空间复杂度O(n)
 * 解决:随机partion 生成[l,r]的随机值
 * ```js
 * const randIndex = randint(l, r)
  ;[[arr[l], arr[randIndex]]] = [[arr[randIndex], arr[l]]]
 * ```
 * @description 还是存在的问题 元素有很多是相同的 退化成 时间复杂度O(n^2) 空间复杂度O(n)
 * 解决:重新设计partition算法 双路快排
 */
const qucikSort = (arr: number[], l: number, r: number): void => {
  if (arr.length <= 1) return
  if (l >= r) return
  // 最基础的partition
  const pivotIndex = partition(arr, l, r)
  qucikSort(arr, l, pivotIndex - 1)
  qucikSort(arr, pivotIndex + 1, r)
}

/**
 *
 * @description arr[l]是pivot
 * @description 将小于pivot的数移到pivot前面 形成一个递增的序列
 */
const partition = (arr: number[], l: number, r: number) => {
  // 优化，随机取标定点，以解决近乎有序的列表
  const randIndex = randint(l, r)
  ;[[arr[l], arr[randIndex]]] = [[arr[randIndex], arr[l]]]

  let pivotIndex = l
  const pivot = arr[l]
  for (let i = l + 1; i <= r; i++) {
    if (arr[i] < pivot) {
      pivotIndex++
      ;[[arr[i], arr[pivotIndex]]] = [[arr[pivotIndex], arr[i]]]
    }
  }

  // pivot放中间
  ;[[arr[l], arr[pivotIndex]]] = [[arr[pivotIndex], arr[l]]]

  return pivotIndex
}

/**
 * @description 生成[start,end]的随机整数
 */
const randint = (start: number, end: number) => {
  if (start > end) throw new Error('invalid interval')
  const amplitude = end - start
  return Math.floor((amplitude + 1) * Math.random()) + start
}

const arr = [4, 1, 2, 5, 3, 6, 7]
qs(arr)
console.log(arr)

export {}
