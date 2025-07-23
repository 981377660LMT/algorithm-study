/**
 * 查询区间`[start, end)`中值为`value`的元素个数.
 */
function rangeFreq<T>(arr: ArrayLike<T>): (start: number, end: number, value: T) => number {
  const mp = new Map<T, number[]>()
  for (let i = 0; i < arr.length; i++) {
    const v = arr[i]
    if (!mp.has(v)) mp.set(v, [])
    mp.get(v)!.push(i)
  }

  const bisectLeft = (nums: ArrayLike<number>, target: number): number => {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (nums[mid] < target) left = mid + 1
      else right = mid - 1
    }
    return left
  }

  const bisectRight = (nums: ArrayLike<number>, target: number): number => {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (nums[mid] <= target) left = mid + 1
      else right = mid - 1
    }
    return left
  }

  return (start: number, end: number, value: T): number => {
    const pos = mp.get(value)
    if (!pos) return 0
    return bisectRight(pos, end - 1) - bisectLeft(pos, start)
  }
}

if (require.main === module) {
  const arr = [1, 1, 1, 4, 5, 6, 7, 8, 9]
  const freq = rangeFreq(arr)
  console.log(freq(0, 2, 1))
}

export { rangeFreq }
