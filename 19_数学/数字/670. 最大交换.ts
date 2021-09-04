/**
 * @param {number} num
 * @return {number}
 * 给定一个非负整数，你至多可以交换一次数字中的任意两位。返回你能得到的最大值。
 * @summary 让最高位的小数和后面的大数交换 大数要尽可能后
 * 需要交换的位置处于第一条单调递减线之中，由此可以引入单调栈进行求解。
 * @link https://leetcode-cn.com/problems/maximum-swap/solution/dan-diao-zhan-yu-tan-xin-gui-lu-by-user5707f/
 */
var maximumSwap = function (num: number): number {
  const res = num.toString().split('').map(Number)
  // 单减的栈
  const stack: number[] = [res[0]]
  let record = 1

  while (res[record] <= stack[stack.length - 1]) {
    stack.push(res[record])
    record++
  }

  let max = res[record]
  // 寻找后面的最大值 相等时取后面的
  for (let i = record; i < res.length; i++) {
    if (res[i] >= max) {
      max = res[i]
      record = i
    }
  }

  // 找到最前面的可交换位置
  while (stack.length && stack[stack.length - 1] < max) {
    stack.pop()
  }

  ;[res[stack.length], res[record]] = [res[record], res[stack.length]]

  return parseInt(res.join(''))
}

console.log(maximumSwap(2736))

export default 1
