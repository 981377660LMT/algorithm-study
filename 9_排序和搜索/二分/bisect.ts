/* eslint-disable no-inner-declarations */
import assert from 'assert'

interface BisectOptions<E, T> {
  /**
   * 二分查找的起始索引
   */
  left?: number

  /**
   * 二分查找的结束索引
   */
  right?: number

  /**
   * 将数组中的元素转换为比较的值
   */
  key?: (e: E) => T
}

/**
 * 返回 `target` 在 `arrayLike` 中最左边的插入位置。
 * 存在多个相同的值时，返回最左边的位置。
 *
 * @param arrayLike 某个参数有序的数组
 * @param target 查找的目标值
 * @param options {@link BisectOptions}
 */
function bisectLeft<E, T>(
  arrayLike: ArrayLike<E>,
  target: T,
  options?: BisectOptions<E, T>
): number {
  const n = arrayLike.length
  if (n === 0) return 0

  let { left = 0, right = n - 1, key = (e: E) => e } = options ?? {}
  while (left <= right) {
    const mid = (left + right) >>> 1
    const midElement = key(arrayLike[mid])
    if (midElement < target) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }

  return left
}

/**
 * 返回 `target` 在 `arrayLike` 中最右边的插入位置。
 * 存在多个相同的值时，返回最右边的位置。
 *
 * @param arrayLike 某个参数有序的数组
 * @param target 查找的目标值
 * @param options {@link BisectOptions}
 */
function bisectRight<E, T>(
  arrayLike: ArrayLike<E>,
  target: T,
  options?: BisectOptions<E, T>
): number {
  const n = arrayLike.length
  if (n === 0) return 0

  let { left = 0, right = n - 1, key = (e: E) => e } = options ?? {}
  while (left <= right) {
    const mid = (left + right) >>> 1
    const midElement = key(arrayLike[mid])
    if (midElement <= target) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }

  return left
}

/**
 * 在 `array` 中插入 `target`，并保持 `array` 有序。
 * 如果 `array` 中存在多个相同的值，插入到`最左边`的位置。
 *
 * @param array 某个参数有序的数组
 * @param target 插入的目标值
 * @param options {@link BisectOptions}
 */
function insortLeft<E>(array: E[], target: E, options?: BisectOptions<E, E>): void {
  const pos = bisectLeft(array, target, options)
  array.splice(pos, 0, target)
}

/**
 * 在 `array` 中插入 `target`，并保持 `array` 有序。
 * 如果 `array` 中存在多个相同的值，插入到`最右边`的位置。
 * !适合用于维护长度较小(2000左右)的有序数组。
 *
 * @param array 某个参数有序的数组
 * @param target 插入的目标值
 * @param options {@link BisectOptions}
 *
 * @alias bisectInsort
 */
function insortRight<E>(array: E[], target: E, options?: BisectOptions<E, E>): void {
  const pos = bisectRight(array, target, options)
  array.splice(pos, 0, target)
}

const bisectInsort = insortRight

if (require.main === module) {
  const arr0 = [-3, -1, 1, 3]
  assert.strictEqual(bisectLeft(arr0, 1), 2)
  assert.strictEqual(bisectRight(arr0, 1), 3)
  assert.strictEqual(bisectLeft(arr0, 5), 4)
  assert.strictEqual(bisectRight(arr0, 5), 4)

  const arr1 = [1, 2, 2, 2, 3, 3, 4, 5, 6, 7]
  assert.strictEqual(bisectRight(arr1, 3), 6)
  assert.strictEqual(bisectLeft(arr1, 3), 4)
  assert.strictEqual(bisectRight(arr1, 2), 4)
  assert.strictEqual(bisectLeft(arr1, 2), 1)

  const arr3 = [1, 2, 2, 2, 3]
  insortLeft(arr3, 2)
  assert.deepStrictEqual(arr3, [1, 2, 2, 2, 2, 3])
  insortRight(arr3, 2)
  assert.deepStrictEqual(arr3, [1, 2, 2, 2, 2, 2, 3])

  // 2763. 所有子数组中不平衡数字之和 (n^3 -> 1400ms)
  // https://leetcode.cn/problems/sum-of-imbalance-numbers-of-all-subarrays/
  function sumImbalanceNumbers(nums: number[]): number {
    let res = 0
    const n = nums.length
    for (let left = 0; left < n; left++) {
      const sl: number[] = []
      for (let right = left; right < n; right++) {
        insortLeft(sl, nums[right])
        for (let i = 1; i < sl.length; i++) {
          res += +(sl[i] - sl[i - 1] > 1)
        }
      }
    }
    return res
  }
}

export { bisectLeft, bisectRight, insortLeft, insortRight, bisectInsort }
