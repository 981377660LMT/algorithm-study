/**
 * @param {number} n
 * @return {number[]}
 * @description
 * 给定一个整数 n, 返回从 1 到 n 的字典顺序。
 * @description 字典树前序遍历 因为要先把1开头的全部看了
 * 字典序排数
 */
function lexicalOrder(n: number): number[] {
  const res: number[] = []
  dfs(0, n)
  return res

  function dfs(cur: number, limit: number): void {
    for (let i = 0; i <= 9; i++) {
      const next = cur * 10 + i
      // when larger than n, return to the previous level
      if (next > limit) {
        return
      }
      if (next === 0) {
        continue
      }
      res.push(next)
      dfs(next, limit)
    }
  }
}

const lexicalOrder2 = function (n: number): number[] {
  return Array.from({ length: n }, (_, i) => i + 1).sort()
}

console.log(lexicalOrder(13))
// 返回 [1,10,11,12,13,2,3,4,5,6,7,8,9]
