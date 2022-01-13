// 我们给出一个由字符 '0' 和 '1' 组成的字符串 S，
// 我们可以将任何 '0' 翻转为 '1' 或者将 '1' 翻转为 '0'。
// 返回使 S 单调递增的最小翻转次数。

// 遍历字符串，找到一个分界点，使得该分界点之前1的个数和分界点之后0的个数之和最小，把分界点之前的1变成0，之后的0变成1
function minFlipsMonoIncr(s: string): number {
  const len = s.length
  // dp[i][j] 表示 在 i 位置的时候，其为 j 状态的时候，保持单调的最小翻转的次数
  const dp = Array.from<unknown, [allZero: number, zeroThenOne: number]>(
    { length: len + 1 },
    () => [0, 0]
  )

  for (let i = 0; i < len; i++) {
    if (s[i] === '0') {
      dp[i + 1][0] = dp[i][0]
      dp[i + 1][1] = Math.min(dp[i][0], dp[i][1]) + 1
    } else {
      dp[i + 1][0] = dp[i][0] + 1
      dp[i + 1][1] = Math.min(dp[i][0], dp[i][1])
    }
  }

  return Math.min.apply(null, dp[len])
}

console.log(minFlipsMonoIncr('00110'))
// 解释：我们翻转最后一位得到 00111.
console.log(minFlipsMonoIncr('010110'))
// 解释：我们翻转得到 011111，或者是 000111。
console.log(minFlipsMonoIncr('00011000'))
// 解释：我们翻转得到 00000000。
// 在末尾新添加一个字符时，
// 如果该字符为“1”，出现“1”的次数为 one+1，最佳解不变仍为 t 次。
// 如果该字符为“0”，出现“1”的次数不变仍为 one，最佳解有两种情况：
// 1）将末尾“0”转为“1”，则共需要 t+1 次；
// 2）末尾“0”不变，将字符串中的“1”全部转为“0”，则共需要 one 次。
// 取两种情况的最小值，即 min(one, t + 1)。
