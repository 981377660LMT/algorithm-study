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
  /**
   *
   * @description arr[l]是pivot
   * @description 要分为小于pivot 等于pivot 大于pivot三个部分
   */
  const partition = (arr: number[], l: number, r: number): number => {
    /**
     * @description 生成[start,end]的随机整数
     */
    const randint = (start: number, end: number) => {
      if (start > end) throw new Error('invalid interval')
      const amplitude = end - start
      return Math.floor((amplitude + 1) * Math.random()) + start
    }

    // 优化，随机取标定点，以解决近乎有序的列表
    const randIndex = randint(l, r)
    ;[[arr[l], arr[randIndex]]] = [[arr[randIndex], arr[l]]]

    let pivotIndex = l
    const pivot = arr[l]
    for (let i = l + 1; i <= r; i++) {
      if (arr[i] < pivot) {
        // 这里要先移pivotIndex是因为不能动最左边的比较元素 比较元素要最后移到自己的位置
        pivotIndex++
        ;[[arr[i], arr[pivotIndex]]] = [[arr[pivotIndex], arr[i]]]
      }
    }

    // pivot放到中间应有的位置
    ;[[arr[l], arr[pivotIndex]]] = [[arr[pivotIndex], arr[l]]]

    return pivotIndex
  }

  if (l < r) {
    // 最基础的partition
    const pivotIndex = partition(arr, l, r)
    qucikSort(arr, l, pivotIndex - 1)
    qucikSort(arr, pivotIndex + 1, r)
  }
}

const arr = [4, 1, 2, 5, 3, 6, 7]
qucikSort(arr, 0, arr.length - 1)
console.log(arr)

export {}
