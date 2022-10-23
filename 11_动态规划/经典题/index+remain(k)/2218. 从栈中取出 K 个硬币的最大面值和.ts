/* eslint-disable implicit-arrow-linebreak */
// 5269. 从栈中取出 K 个硬币的最大面值和

// 1 <= n <= 1000
// !1 <= k <= sum(piles[i].length) <= 2000
// 1 <= piles[i][j] <= 1e5

// !时间复杂度为O(k*sum(piles[i].length))
// 换句话说就是暴力解法，对每个栈都讨论了当前能取到的最大值

function maxValueOfCoins(piles: number[][], limit: number): number {
  const n = piles.length
  let dp = new Uint32Array(limit + 1).fill(0)
  for (let i = 0; i < n; i++) {
    let preSum = 0
    const ndp = new Uint32Array(limit + 1).fill(0)
    for (let j = 0; j < piles[i].length; j++) {
      preSum += piles[i][j]
      for (let k = j + 1; k <= limit; k++) {
        ndp[k] = Math.max(Math.max(dp[k], dp[k - 1 - j] + preSum), ndp[k])
      }
    }
    dp = ndp
  }

  return dp[limit]
}

if (require.main === module) {
  const piles = Array.from({ length: 2000 }, () =>
    Array.from({ length: 1000 }, () => Math.floor(Math.random() * 100000))
  )
  const limit = 1000
  console.time('start')
  maxValueOfCoins(piles, limit)
  console.timeEnd('start') // start: 1.932s
}

export {}
