/**
 * @param {number[]} height
 * @return {number}
 * @description 计算按此排列的柱子，下雨之后能接多少雨水。
 * @summary 按照列来计算
 * 每一列雨水的高度，取决于，该列 左侧最高的柱子和右侧最高的柱子中最矮的那个柱子的高度
 */
const trap = function (height: number[]): number {
  let res = 0
  let max = 0
  const leftMax: number[] = []
  const rightMax: number[] = []

  for (let i = 0; i < height.length; i++) {
    max = Math.max(max, height[i])
    leftMax[i] = max
  }

  max = 0

  for (let i = height.length - 1; i >= 0; i--) {
    max = Math.max(max, height[i])
    rightMax[i] = max
  }

  for (let i = 0; i < height.length; i++) {
    res += Math.min(leftMax[i], rightMax[i]) - height[i]
  }

  return res
}
console.log(trap([0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1]))
console.log(trap([4, 2, 0, 3, 2, 5]))

export default 1
