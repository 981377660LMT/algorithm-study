// 我们设将 a 行 b 列涂成黑色，会形成 sum 个黑色格子，
// 而 sum =a⋅n+b⋅n−a⋅b 。
// 遍历统计所有 sum = k 的涂色方案即可。

function paintingPlan(n: number, k: number): number {
  if (k === n * n) return 1
  let res = 0

  for (let a = 0; a < n; a++) {
    for (let b = 0; b < n; b++) {
      const sum = a * n + b * n - a * b
      if (sum === k) {
        res += combination(n, a) * combination(n, b)
      }
    }
  }

  return res

  function combination(n: number, k: number): number {
    let res = 1
    for (let i = n; i > k; i--) res *= i
    for (let i = n - k; i > 0; i--) res /= i
    return res
  }
}

console.log(paintingPlan(2, 2))

// 输出：4

// 解释：一共有四种不同的方案：
// 第一种方案：涂第一列；
// 第二种方案：涂第二列；
// 第三种方案：涂第一行；
// 第四种方案：涂第二行。

// 小扣可以选择任意多行以及任意多列的格子涂成黑色（选择的整行、整列均需涂成黑色），
// 所选行数、列数均可为 0。
// 小扣希望最终的成品上需要有 k 个黑色格子，请返回小扣共有多少种涂色方案。

// 输入：n = 2, k = 1
// 输出：0
// 解释：不可行，因为第一次涂色至少会涂两个黑格。
