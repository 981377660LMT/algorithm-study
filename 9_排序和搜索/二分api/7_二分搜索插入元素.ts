import { bisectLeft } from './7_二分搜索寻找最左插入位置'

/**
 * @description
 * 如果 nums[mid] 等于目标值， 则提前返回 mid（只需要找到一个满足条件的即可）
   如果 nums[mid] 小于目标值， 说明目标值在 mid 右侧，这个时候解空间可缩小为 [mid + 1, right] （mid 以及 mid 左侧的数字被我们排除在外）
   如果 nums[mid] 大于目标值， 说明目标值在 mid 左侧，这个时候解空间可缩小为 [left, mid - 1] （mid 以及 mid 右侧的数字被我们排除在外）
 */
const bisectInsort = (arr: number[], target: number): void => {
  const index = bisectLeft(arr, target)
  arr.splice(index, 0, target)
}

if (require.main === module) {
  // const arr = [7, 7, 7, 7, 7, 7]
  // const arr = [6, 7, 8, 9, 10]
  const arr: number[] = []
  bisectInsort(arr, 11)
  console.log(arr)
}

export { bisectInsort }
