/**
 *
 * @param arr
 * @returns
 * @description
 *（1）在数据集之中，选择一个元素作为"基准"（pivot）。
 *（2）所有小于"基准"的元素，都移到"基准"的左边；所有大于"基准"的元素，都移到"基准"的右边。
 *（3）对"基准"左边和右边的两个子集，不断重复第一步和第二步，直到所有子集只剩下一个元素为止。
 * @description 存在的问题:完全有序的数组 时间复杂度O(n^2) 空间复杂度O(n)
 * 解决:随机partion 生成[l,r]的随机值
 */
const qucikSort = (arr: number[], l: number, r: number): void => {
  /**
   *
   * @param arr
   * @param l
   * @param r
   * @description arr[l]是pivot
   * @description 要分为小于pivot 等于pivot 大于pivot三个部分
   */
  const partition = (arr: number[], l: number, r: number) => {
    let pivotIndex = l
    const pivot = arr[l]

    for (let i = l + 1; i <= r; i++) {
      if (arr[i] < pivot) {
        pivotIndex++
        ;[[arr[i], arr[pivotIndex]]] = [[arr[pivotIndex], arr[i]]]
      }
    }

    // 待比较数据放入最终位置
    ;[[arr[l], arr[pivotIndex]]] = [[arr[pivotIndex], arr[l]]]

    // 返回待比较数据最终位置
    return pivotIndex
  }

  if (l >= r) return
  if (arr.length <= 1) return

  // 最基础的partition
  const pivotIndex = partition(arr, l, r)
  qucikSort(arr, l, pivotIndex - 1)
  qucikSort(arr, pivotIndex + 1, r)
}

const arr = [4, 1, 2, 5, 3, 6, 7]
qucikSort(arr, 0, arr.length - 1)
console.log(arr)

export {}
