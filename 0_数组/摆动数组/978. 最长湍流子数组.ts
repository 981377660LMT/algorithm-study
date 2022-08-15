/**
 * @param {number[]} arr
 * @return {number}
 * @description 如果比较符号在子数组中的每个相邻元素对之间翻转，则该子数组是湍流子数组。
 * @summary 其实是dp
 * 这个交替数组的问题很经典
 */
const maxTurbulenceSize = function (arr: number[]): number {
  let res = 1
  let up = 1
  let down = 1

  for (let i = 1; i < arr.length; i++) {
    if (arr[i] > arr[i - 1]) {
      up = down + 1
      down = 1
    } else if (arr[i] < arr[i - 1]) {
      down = up + 1
      up = 1
    } else {
      up = 1
      down = 1
    }

    res = Math.max(res, up, down)
  }

  return res
}

console.log(maxTurbulenceSize([9, 4, 2, 10, 7, 8, 8, 1, 9]))
// 输出：5
// 解释：(A[1] > A[2] < A[3] > A[4] < A[5])
