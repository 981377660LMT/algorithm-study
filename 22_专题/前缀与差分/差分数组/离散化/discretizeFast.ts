/**
 * 将nums中的元素进行离散化，返回新的数组和对应的原始值.
 * `origin[newNums[i]] == nums[i]`.
 * @param nums 原始数组
 * @param uniqueIdForSameValue 是否对相同的值进行唯一化.
 */
function discretizeFast<T extends ArrayLike<number>>(
  nums: T,
  uniqueIdForSameValue = false
): { newNums: Uint32Array; origin: number[] } {
  const newNums = new Uint32Array(nums.length)
  const origin = Array<number>(nums.length)
  let ptr = 0
  const order = argSort(nums.length, (i, j) => nums[i] - nums[j])
  if (uniqueIdForSameValue) {
    for (let i = 0; i < nums.length; i++) {
      origin[i] = nums[order[i]]
      newNums[order[i]] = i
    }
    return { newNums, origin }
  }

  for (let i = 0; i < nums.length; i++) {
    if (ptr === 0 || origin[ptr - 1] !== nums[order[i]]) {
      origin[ptr++] = nums[order[i]]
    }
    newNums[order[i]] = ptr - 1
  }
  origin.length = ptr
  return { newNums, origin }
}

function argSort(n: number, compareFn: (i: number, j: number) => number): Uint32Array {
  const order = new Uint32Array(n)
  for (let i = 0; i < n; i++) order[i] = i
  order.sort(compareFn)
  return order
}

export { discretizeFast }

if (require.main === module) {
  const nums = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5]
  const { newNums, origin } = discretizeFast(nums)
  console.log(newNums, origin)
  const { newNums: newNums2, origin: origin2 } = discretizeFast(nums, true)
  console.log(newNums2, origin2)

  for (let i = 0; i < nums.length; i++) {
    if (origin[newNums[i]] !== nums[i]) {
      throw new Error('error')
    }
  }
}
