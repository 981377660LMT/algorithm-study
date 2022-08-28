// 找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。
// 返回容器可以储存的最大水量。

function maxArea(height: number[]): number {
  const n = height.length
  let res = 0
  let left = 0
  let right = n - 1
  while (left < right) {
    res = Math.max(res, Math.min(height[left], height[right]) * (right - left))
    if (height[left] < height[right]) {
      left++
    } else {
      right--
    }
  }

  return res
}

console.log(maxArea([1, 8, 6, 2, 5, 4, 8, 3, 7]))
// 输出：49
