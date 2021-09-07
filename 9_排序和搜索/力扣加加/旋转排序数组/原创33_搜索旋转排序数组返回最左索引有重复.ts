/**
 * 元素可能相同
 * @summary 如果存在相等 返回最左边的索引
 *
 */
const search = function (arr: number[], target: number): number {
  let l = 0
  let r = arr.length - 1

  while (l <= r) {
    const mid = (l + r) >> 1

    // # 重点1：当left符合时直接返回, 因为找的是最小的索引
    if (arr[l] === target) return l
    // # 重点2：当中间值等于目标值，将右边界移到中间，因为左边可能还有相等的值
    if (arr[mid] === target) r = mid

    if (arr[l] < arr[mid]) {
      if (arr[mid] > target && target >= arr[l]) {
        r = mid - 1
      } else {
        l = mid + 1
      }
    } else if (arr[l] > arr[mid]) {
      if (arr[r] >= target && target > arr[mid]) {
        l = mid + 1
      } else {
        r = mid - 1
      }
      // # 重点3：当中间数字与左边数字相等时，将左边界右移
    } else if (arr[l] === arr[mid]) {
      l++
    }
  }

  // 兜底
  return arr.indexOf(target)
}

// 遇到low和mid相等的情况下，因为有重复元素，并不能确认前半部分是有序的
console.log(search([5, 5, 5, 1, 2, 3, 4, 5], 5))
console.log(search([15, 16, 19, 20, 25, 1, 3, 4, 5, 7, 10, 14], 5))
// 0
