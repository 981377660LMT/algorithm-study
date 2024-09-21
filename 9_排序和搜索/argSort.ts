/* eslint-disable max-len */

/**
 * 将 {@link n} 个元素按照 {@link less} 函数的比较结果排序, 返回排序后的索引数组.
 *
 * @returns order `order[i]` 表示排序后的第 `i` 个元素在原数组中的索引.
 * @example
 * ```ts
 * const arr = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3]
 * const order = argSort(arr.length, (i, j) => arr[i] < arr[j])
 * console.log(order) // [1, 3, 0, 9, 6, 2, 4, 8, 7, 5]
 * console.log(reArrange(arr, order)) // [1, 1, 2, 3, 3, 4, 5, 5, 6, 9]
 * ```
 */
function argSort(n: number, less: (i: number, j: number) => boolean): number[] {
  const order = Array(n)
  for (let i = 0; i < order.length; i++) order[i] = i
  order.sort((a, b) => (less(a, b) ? -1 : 1))
  return order
}

/**
 * 将 {@link arr} 按照 {@link order} 重新排列.
 */
function reArrange<T>(arr: T[], order: ArrayLike<number>): T[] {
  const res = Array(arr.length)
  for (let i = 0; i < res.length; i++) res[i] = arr[order[i]]
  return res
}

export { argSort, reArrange }

if (require.main === module) {
  const arr = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3]
  const order = argSort(arr.length, (i, j) => arr[i] < arr[j])
  console.log(order) // [1, 3, 0, 9, 6, 2, 4, 8, 7, 5]
  console.log(reArrange(arr, order)) // [1, 1, 2, 3, 3, 4, 5, 5, 6, 9]
}
