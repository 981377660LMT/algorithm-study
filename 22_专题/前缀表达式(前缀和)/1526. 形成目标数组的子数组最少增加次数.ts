// 在 initial 中选择 任意 子数组，并将子数组中每个元素增加 1
// 122. 买卖股票的最佳时机 II
// 升序加两个元素的差值，降序不管。累加求和，最后再加第一个元素即可
function minNumberOperations(target: number[]): number {
  let res = target[0]
  for (let i = 1; i < target.length; i++) {
    res += Math.max(0, target[i] - target[i - 1])
  }
  return res
}

console.log(minNumberOperations([3, 1, 1, 2]))
// 输出：4
// 这道题我们可以逆向思考返回从 target 得到 initial 的最少操作次数。
export default 1
