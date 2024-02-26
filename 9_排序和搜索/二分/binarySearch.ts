/**
 * 查询非递减数组中目标值的第一个或最后一个位置.
 * @param arr 非递减数组
 * @param target 目标值
 * @param findFirst 是否查询第一个位置. 默认为 true.
 * @returns 目标值的第一个或最后一个位置. 如果目标值不存在, 返回 undefined.
 */
function binarySearch(
  arr: ArrayLike<number>,
  target: number,
  findFirst = true
): number | undefined {
  if (!arr.length || arr[0] > target || arr[arr.length - 1] < target) return -1
  if (findFirst) {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = left + Math.floor((right - left) / 2)
      if (arr[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left < arr.length && arr[left] === target ? left : -1
  } else {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = left + Math.floor((right - left) / 2)
      if (arr[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left > 0 && arr[left - 1] === target ? left - 1 : -1
  }
}

export { binarySearch }

if (require.main === module) {
  // https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/
  function searchRange(nums: number[], target: number): number[] {
    const pos1 = binarySearch(nums, target) ?? -1
    const pos2 = binarySearch(nums, target, false) ?? -1
    return [pos1, pos2]
  }
}
