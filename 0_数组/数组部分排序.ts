// 原地排序并返回引用
const partialSort = <T = number>(
  nums: T[],
  start: number,
  end: number,
  compareFn?: (a: T, b: T) => number
) => {
  const preSorted = nums.slice(0, start)
  const postSorted = nums.slice(end)
  const sorted = nums.slice(start, end).sort(compareFn)
  nums.length = 0
  nums.push.apply(nums, preSorted.concat(sorted, postSorted))
  return nums
}

if (require.main === module) {
  const array = [5, 2, 6, 4, 0, 1, 9, 3, 8, 7]
  partialSort(array, 3, 7)
  console.log(array)
}

export { partialSort }
