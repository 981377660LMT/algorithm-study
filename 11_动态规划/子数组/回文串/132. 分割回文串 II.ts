/**
 * @param {string} s   s 的长度可以到达 2000
 * @return {number}
 * 类似于300. 最长递增子序列
 */
var minCut = function (s: string): number {
  const isPalindrome = (str: string) => str === str.split('').reverse().join('')
  const n = s.length
  const dp = Array(n).fill(n - 1)

  for (let i = 0; i < n; i++) {
    if (isPalindrome(s.slice(0, i + 1))) {
      dp[i] = 0
      continue
    }
    for (let j = 0; j < i; j++) {
      if (isPalindrome(s.slice(j + 1, i + 1))) dp[i] = Math.min(dp[i], dp[j] + 1)
    }
  }

  return dp[n - 1]
}

console.log(minCut('aab'))
// 输出：1
// 解释：只需一次分割就可将 s 分割成 ["aa","b"] 这样两个回文子串。
// 状态转移方程：dp[i] = min(dp[i], dp[j] + 1);dp[i]=min(dp[i],dp[j]+1); 如果 isPalindrome(s[j + 1..i])isPalindrome(s[j+1..i])
