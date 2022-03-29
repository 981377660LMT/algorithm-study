/**
 * @param {number[]} piles
 * @return {number}
 * @description
 * 亚历克斯先开始。最初，M = 1。
   在每个玩家的回合中，该玩家可以拿走剩下的 前 X 堆的所有石子，
   其中 1 <= X <= 2M。然后，令 M = max(M, X)。
 * 返回亚历克斯可以得到的最大数量的石头。
 */
const stoneGameII = function (piles: number[]): number {
  const len = piles.length
  // 从i堆开始时剩下的石子的个数/后缀和
  const remain = Array<number>(len).fill(0)
  remain[len - 1] = piles[len - 1]
  for (let i = len - 2; i >= 0; i--) {
    remain[i] = remain[i + 1] + piles[i]
  }
  const memo = new Map<string, number>()

  const dfs = (index: number, M: number): number => {
    // 可以拿完剩下的石头
    if (len - index <= 2 * M) return remain[index]

    const key = `${index}#${M}`
    if (memo.has(key)) return memo.get(key)!

    let res = -Infinity
    for (let x = 1; x <= 2 * M; x++) {
      const newM = Math.max(M, x)
      // 减去对手最多拿的
      res = Math.max(res, remain[index] - dfs(index + x, newM))
    }

    memo.set(key, res)
    return res
  }
  return dfs(0, 1)
}

console.log(stoneGameII([2, 7, 9, 4, 4]))

export {}
