/**
 * @param {number} target
 * @param {number} inc 从 x 号站点移动至 x + 1 号站点需要花费的时间为 inc
 * @param {number} dec 从 x 号站点移动至 x - 1 号站点需要花费的时间为 dec
 * @param {number[]} jump
 * @param {number[]} cost
 * @return {number}
 */
var busRapidTransit = function (
  target: number,
  inc: number,
  dec: number,
  jump: number[],
  cost: number[]
): number {
  const MOD = 10e9 + 7
  const memo = new Map<number, number>()

  /**
   *
   * @param pos
   * @returns
   * 我们也可采取逆向思维，即从 target 出发返回 0
   * 从pos走到0的耗费
   */
  const dfs = (pos: number): number => {
    if (pos === 0) return 0
    if (pos === 1) return inc
    if (memo.has(pos)) return memo.get(pos)!

    let res = pos * inc
    for (let i = 0; i < jump.length; i++) {
      const [div, mod] = [~~(pos / jump[i]), pos % jump[i]]
      if (mod === 0) {
        res = Math.min(res, cost[i] + dfs(div))
      } else {
        // 往后走到公交下一站，往前走到公交站
        res = Math.min(
          res,
          cost[i] + dfs(div) + inc * mod,
          cost[i] + dfs(div + 1) + (jump[i] - mod) * dec
        )
      }
    }

    memo.set(pos, res)
    return res
  }

  return dfs(target) % MOD
}

console.log(busRapidTransit(31, 5, 3, [6], [10]))
