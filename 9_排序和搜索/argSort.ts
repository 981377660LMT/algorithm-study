/* eslint-disable max-len */

/**
 * 返回arr排序后的索引数组.
 */
function argSort<T>(arr: T[], compareFn: (a: T, b: T) => number = (a: any, b: any) => a - b): number[] {
  const order = Array(arr.length)
  for (let i = 0; i < order.length; i++) {
    order[i] = i
  }
  order.sort((a, b) => compareFn(arr[a], arr[b]))
  return order
}

/**
 * 将arr按照order的顺序重新排列.
 */
function reArrange<T>(arr: T[], order: ArrayLike<number>): T[] {
  const res = Array(arr.length)
  for (let i = 0; i < res.length; i++) {
    res[i] = arr[order[i]]
  }
  return res
}

export { reArrange }

if (require.main === module) {
  const nums = [5, -7, 3, 5]
  const order = argSort(nums, (a, b) => a - b)
  console.log(order)
  console.log(reArrange(nums, order))
}
