/**
 * @description
 * 如果 nums[mid] 等于目标值， 则提前返回 mid（只需要找到一个满足条件的即可）
   如果 nums[mid] 小于目标值， 说明目标值在 mid 右侧，这个时候解空间可缩小为 [mid + 1, right] （mid 以及 mid 左侧的数字被我们排除在外）
   如果 nums[mid] 大于目标值， 说明目标值在 mid 左侧，这个时候解空间可缩小为 [left, mid - 1] （mid 以及 mid 右侧的数字被我们排除在外）
 */
const biSearch = (arr: number[], target: number): number => {
  if (arr.length === 0) return -1

  let l = 0
  let r = arr.length - 1
  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = arr[mid]
    if (midElement === target) {
      return mid
    } else if (midElement < target) {
      l = mid + 1
    } else {
      r = mid - 1
    }
  }

  return -1
}

const arr = [1, 2, 3, 4, 5, 6, 7]
console.log(biSearch(arr, 3))

export {}
