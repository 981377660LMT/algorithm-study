// 给定一组城市和每对城市之间的火车票的价钱，找到每个城市只访问一次并返回起点的最小车费花销。
// 复杂度O(n*2^n)
// https://www.acwing.com/problem/content/93/
// . 最短哈密尔顿路径
function tsp(n: number, weight: number[][], start = 0): number {
  const target = (1 << n) - 1
  const memo = new Map<string, number>()

  const dfs = (cur: number, visited: number): number => {
    if (cur === start && visited === target) return 0
    const key = `${cur}#${visited}`
    if (memo.has(key)) return memo.get(key)!

    let res = Infinity
    for (let next = 0; next < n; next++) {
      if (visited & (1 << next)) continue
      res = Math.min(res, weight[cur][next] + dfs(next, visited | (1 << next)))
    }

    memo.set(key, res)
    return res
  }

  return dfs(start, 0)
}

console.log(
  tsp(4, [
    [0, 2, 6, 5],
    [2, 0, 4, 4],
    [6, 4, 0, 2],
    [5, 4, 2, 0],
  ])
)
