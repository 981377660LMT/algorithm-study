/**
 * 返回一个数字数组的排序后的索引.内部使用计数排序.
 */
function argSortCounting(nums: ArrayLike<number>, minArg: number, maxArg: number): Uint32Array {
  const counter = new Int32Array(maxArg - minArg + 1)
  for (let i = 0; i < nums.length; i++) {
    counter[nums[i] - minArg]++
  }
  for (let i = 1; i < counter.length; i++) {
    counter[i] += counter[i - 1]
  }
  const order = new Uint32Array(nums.length)
  // 值相等时，按照下标从小到大排序.
  for (let i = nums.length - 1; i >= 0; i--) {
    const v = nums[i] - minArg
    counter[v]--
    order[counter[v]] = i
  }
  return order
}

/**
 * 返回数组的排序索引.
 */
function argSort(arr: ArrayLike<number>): Uint32Array {
  const n = arr.length
  const order = new Uint32Array(n)
  for (let i = 0; i < n; i++) order[i] = i
  order.sort((i, j) => arr[i] - arr[j])
  return order
}

/**
 * 按照索引数组重新排列数组.
 */
function reArrage<T>(arr: ArrayLike<T>, order: ArrayLike<number>): T[] {
  const n = arr.length
  const res = Array(n)
  for (let i = 0; i < n; i++) res[i] = arr[order[i]]
  return res
}

export { argSortCounting, argSort, reArrage }

if (require.main === module) {
  const nums = [3, 1, 2, 4, 3, 1, 2, 4]
  console.log(argSortCounting(nums, 1, 4))
  console.log(argSort(nums))
  console.log(reArrage(nums, argSort(nums)))
}
