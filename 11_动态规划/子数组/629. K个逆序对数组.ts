// /**
//  *
//  * @param n  有包含从 1 到 n 的数字
//  * @param k  恰好拥有 k 个逆序对的不同的数组的个数
//  */
// function kInversePairs(n: number, k: number): number {
//   let mod = 10 ** 9 + 7
//   if (k > (n * (n - 1)) / 2 || k < 0) return 0
//   const dp = Array.from({ length: n + 1 }, () => Array(k + 1).fill(0))
//   for (let i = 0; i <= n; i++) {
//     dp[i][0] = 1
//   }

//   for (let i = 0; i < array.length; i++) {
//     const element = array[i]
//   }
//   return dp[n][k]
// }

// console.log(kInversePairs(3, 0))

// export {}
