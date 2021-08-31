// O(n)时间复杂
// 最大面积一定是当前的面积或者通过移动短的端点得到
const maxArea = (height: number[]) => {
  let max = 0
  let left = 0
  let right = height.length - 1
  while (left < right) {
    const cur = (right - left) * Math.min(height[left], height[right])
    if (cur > max) max = cur
    // 移动的方法
    if (height[left] < height[right]) {
      left++
    } else {
      right--
    }
  }

  return max
}

console.log(maxArea([4, 3, 2, 1, 4]))
// 输出：16
export {}
