/**
 * @param {number[]} arr  一个排序好的数组 arr
 * @param {number} k
 * @param {number} x
 * @return {number[]}  从数组中找到最靠近 x（两数之差最小）的 k 个数
 * 找到最接近 x 的 k 个数
 * a 比整数 b 更接近 x 需要满足 |a - x| < |b - x| 或 |a - x| === |b - x| && a < b
 * !答案一定是一个连续的窗口
 * !从窗口左右两端依次排除数字，最后留下k个数
 */
function findClosestElements(arr: number[], k: number, x: number): number[] {
  const n = arr.length
  let left = 0
  let right = n - 1
  let remove = n - k

  for (let _ = 0; _ < remove; _++) {
    if (Math.abs(arr[left] - x) > Math.abs(arr[right] - x)) {
      left++
    } else {
      right--
    }
  }

  return arr.slice(left, right + 1)
}

export {}
