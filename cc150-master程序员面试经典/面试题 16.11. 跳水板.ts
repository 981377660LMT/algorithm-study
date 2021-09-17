// 你必须正好使用k块木板。编写一个方法，生成跳水板所有可能的长度。
// 0 <= k <= 100000
// 这个题不是 DP 或者是 DFS 什么的。看给出的 k 的范围是 100000，我们知道需要用 O(n) 的解法。因此需要找规律了。
function divingBoard(shorter: number, longer: number, k: number): number[] {
  if (k === 0) return []
  if (shorter === longer) return [shorter * k]

  // 此时必有k+1种
  const res = Array<number>(k + 1).fill(0)
  for (let i = 0; i <= k; i++) {
    res[i] = longer * i + shorter * (k - i)
  }

  return res
}

console.log(divingBoard(1, 2, 3))
// 输出： [3,4,5,6]
