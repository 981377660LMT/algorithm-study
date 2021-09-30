/**
 * @param {number[]} height
 * @return {number}
 */
var maxArea = function (height) {
  let max = 0
  let left = 0
  let right = height.length - 1
  while (left < right) {
    const cur = (right - left) * Math.min(height[left], height[right])
    if (cur > max) max = cur
    // 移动小的
    if (height[left] < height[right]) {
      left++
    } else {
      right--
    }
  }

  return max
}

console.log(maxArea([1, 8, 6, 2, 5, 4, 8, 3, 7]))
// 输出：49
