/**
 * @param {number[]} nums
 * @return {number}
 * @description 给你一个整数数组 nums，请你找出并返回能被三整除的元素最大和。
 * @description 思路一:回溯即可
 * 这里用状态机的解法
 * 状态机使用非常广泛，比如正则表达式的引擎，编译器的词法和语法分析，网络协议，企业应用等很多领域都会用到。
 * 我们从左到右扫描数组的过程，将会不断改变状态机的状态。
 * state[0] 表示 mod 为 0 的 最大和
   state[1] 表示 mod 为 1 的 最大和
   state[2] 表示 mod 为 1 的 最大和
 */
const maxSumDivThree = function (nums: number[]): number {
  let state: [number, number, number] = [0, -Infinity, -Infinity]

  for (const num of nums) {
    if (num % 3 === 0) {
      state = [state[0] + num, state[1] + num, state[2] + num]
    } else if (num % 3 === 1) {
      const a = Math.max(state[2] + num, state[0])
      const b = Math.max(state[0] + num, state[1])
      const c = Math.max(state[1] + num, state[2])
      state = [a, b, c]
    } else if (num % 3 === 2) {
      const a = Math.max(state[1] + num, state[0])
      const b = Math.max(state[2] + num, state[1])
      const c = Math.max(state[0] + num, state[2])
      state = [a, b, c]
    }
  }
  console.log(state)
  return state[0]
}

console.log(maxSumDivThree([3, 6, 5, 1, 8]))
