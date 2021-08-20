// // 我们用一个 boolean dp[l][r] 表示字符串从 i 到 j 这段是否为回文
// // 试想如果 dp[l][r]=true，我们要判断 dp[l-1][r+1] 是否为回文。
// const longestPalindrome = (str: string) => {
//   if (str.length <= 1) return str
//   const n = str.length

//   const dp = Array.from({ length: n }, (_, k) =>
//     Array(n)
//       .fill(0)
//       .map((_, index) => (index === k ? true : false))
//   )
//   for (let i = 0; i < n - 1; i++) {
//     for (let j = i + 1; j < n - 1; j++) {
//       dp[i][j] = str[i] === str[j] && dp[i + 1][j - 1]
//     }
//   }

//   return dp
// }

// console.log(longestPalindrome('babad'))

// export {}
