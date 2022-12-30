// 给你一个整数数组 target 和一个数组 initial ，initial 数组与 target  数组有同样的维度，且一开始全部为 0 。
// 请你返回从 initial 得到  target 的最少操作次数，每次操作需遵循以下规则：
// 在 initial 中选择 任意 子数组，并将子数组中每个元素增加 1 。

// 1526. 形成目标数组的子数组最少增加次数
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
