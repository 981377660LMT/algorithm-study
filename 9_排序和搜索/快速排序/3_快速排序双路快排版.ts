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
 * 解决:重新设计partition算法 双路快排  O(nlogn)
 */
const qucikSort = (arr: number[], l: number, r: number): void => {
  /**
   * @description 开始选取pivot arr[l]
   * @description 左右两个指针ij 左边遇到大于等于pivot的停 右边遇到小于等于pivot的停 然后交换他们 i>j则循环终止 标定点移到j
   * @description arr[l+1...i-1]<=v  arr[j+1...r]>=v
   * @summary 双路快排的partition 两重while循环但是O(n)复杂度
   */
  const partition = (arr: number[], l: number, r: number) => {
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

    const pivot = arr[l]
    let leftPoint = l + 1
    let rightPoint = r
    while (leftPoint <= rightPoint) {
      while (leftPoint <= rightPoint && arr[leftPoint] < pivot) {
        leftPoint++
      }
      while (leftPoint <= rightPoint && arr[rightPoint] > pivot) {
        rightPoint--
      }
      if (leftPoint >= rightPoint) break
      ;[[arr[leftPoint], arr[rightPoint]]] = [[arr[rightPoint], arr[leftPoint]]]
      leftPoint++
      rightPoint--
    }

    // pivot放会原来位置
    ;[[arr[l], arr[rightPoint]]] = [[arr[rightPoint], arr[l]]]

    return rightPoint
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
