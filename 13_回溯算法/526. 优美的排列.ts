/**
 * @param {number} n  有从 1 到 n 的 n 个整数
 * @return {number}
 * 只要满足下述条件 之一 ，该数组就是一个 优美的排列 ：
    nums[i] 能够被 i 整除
    i 能够被 nums[i] 整除
 */
const countArrangement = function (n: number): number {
  if (n === 1) return 1
  if (n === 2) return 2
  let res = 0
  const nums = Array.from({ length: n }, (_, i) => i + 1)
  const visited = new Set<number>()

  const bt = (index: number) => {
    if (index === n + 1) return res++
    for (const next of nums) {
      if (visited.has(next)) continue

      if (index % next === 0 || next % index === 0) {
        visited.add(next)
        bt(index + 1)
        visited.delete(next)
      }
    }
  }
  bt(1)
  return res
}

console.log(countArrangement(2))
// 输出：2
// 解释：
// 第 1 个优美的排列是 [1,2]：
//     - nums[1] = 1 能被 i = 1 整除
//     - nums[2] = 2 能被 i = 2 整除
// 第 2 个优美的排列是 [2,1]:
//     - nums[1] = 2 能被 i = 1 整除
//     - i = 2 能被 nums[2] = 1 整除
export {}
