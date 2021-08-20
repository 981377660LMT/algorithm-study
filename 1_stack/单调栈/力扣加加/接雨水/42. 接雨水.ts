/**
 * @param {number[]} height
 * @return {number}
 * @description 计算按此排列的柱子，下雨之后能接多少雨水。
 * @summary
 * 每个柱子顶部可以储水的高度为：该柱子的左右两侧最大高度的较小者减去此柱子的高度。
 * 暴力 时间O(N) 空间O(N)
 *
 */
// const trap = function (height: number[]): number {
//   let res = 0
//   let max = 0
//   const leftMax: number[] = []
//   const rightMax: number[] = []

//   for (let i = 0; i < height.length; i++) {
//     max = Math.max(max, height[i])
//     leftMax[i] = max
//   }

//   max = 0

//   for (let i = height.length - 1; i >= 0; i--) {
//     max = Math.max(max, height[i])
//     rightMax[i] = max
//   }

//   for (let i = 0; i < height.length; i++) {
//     res += Math.min(leftMax[i], rightMax[i]) - height[i]
//   }

//   return res
// }
// /**
//  * @param {number[]} height
//  * @return {number}
//  * @description 计算按此排列的柱子，下雨之后能接多少雨水。
//  * @summary
//  * 每个柱子顶部可以储水的高度为：该柱子的左右两侧最大高度的较小者减去此柱子的高度。
//  * 双指针 时间O(N) 空间O(1)
//  *
//  */
// const trap = function (height: number[]): number {
//   let res = 0
//   let max = 0
//   const leftMax: number[] = []
//   const rightMax: number[] = []

//   for (let i = 0; i < height.length; i++) {
//     max = Math.max(max, height[i])
//     leftMax[i] = max
//   }

//   max = 0

//   for (let i = height.length - 1; i >= 0; i--) {
//     max = Math.max(max, height[i])
//     rightMax[i] = max
//   }

//   for (let i = 0; i < height.length; i++) {
//     res += Math.min(leftMax[i], rightMax[i]) - height[i]
//   }

//   return res
// }
/**
 * @param {number[]} height
 * @return {number}
 * @description 计算按此排列的柱子，下雨之后能接多少雨水。
 * @summary
 * 注意我们的目标是找到下降后第一个上升的点 => 单调栈
 * 形成凹槽才能接雨水(单调栈pop时)
 * 单调栈 时间O(N) 空间O(N))
 * 维护一个单调递减的栈
 *
 */
const trap = function (height: number[]): number {
  let res = 0
  // 存储下标
  const stack: number[] = []

  for (let i = 0; i < height.length; i++) {
    while (stack.length && height[stack[stack.length - 1]] < height[i]) {
      const bottomIndex = stack.pop()!
      // 如果栈顶元素一直相等，那么全都pop出去，只留第一个。
      while (stack.length && height[stack[stack.length - 1]] === height[bottomIndex]) {
        stack.pop()
      }
      if (stack.length) {
        // 高度取左右最小高度减去bottom的高度
        res +=
          (Math.min(height[stack[stack.length - 1]], height[i]) - height[bottomIndex]) *
          (i - 1 - stack[stack.length - 1])
      }
    }
    stack.push(i)
  }
  return res
}
console.log(trap([0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1]))
console.log(trap([4, 2, 0, 3, 2, 5]))
export default 1
