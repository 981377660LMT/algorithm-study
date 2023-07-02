function sumImbalanceNumbers(N: number[]): number {
  if (N.length < 2) return 0

  const nums = new Uint16Array(N)
  const copy = nums.slice()
  const max = Math.max(...nums)
  const counter = new Uint16Array(max + 1)

  let res = 0

  // 枚举子数组
  for (let i = 0; i < nums.length; i++) {
    for (let j = i; j < nums.length; j++) {
      // 优化：如果子数组长度小于 200，直接排序
      if (j + 1 - i <= 200) {
        nums.subarray(i, j + 1).sort()
        for (let k = i; k < j; k++) res += +(nums[k + 1] - nums[k] > 1)
        nums.set(copy.subarray(i, j + 1), i)
      } else {
        // 否则计数排序
        for (let k = i; k <= j; k++) counter[nums[k]]++
        const cur = Array(j - i + 1)
        for (let k = 0, l = 0; k < counter.length; k++) {
          while (counter[k]--) cur[l++] = k
        }
        for (let k = 0; k < cur.length - 1; k++) res += +(cur[k + 1] - cur[k] > 1)
        counter.fill(0)
      }
    }
  }

  return res
}

/**
 * 所有元素都在[0, 2^32)范围内的数组排序.
 */
class SortRangeUint32 {
  private readonly _origin: Uint32Array

  constructor(nums: Uint32Array) {}

  sort(): void {}
}
