// # 给你一个字符串 s ，每个字符是数字 '1' 到 '9' ，再给你两个整数 k 和 minLength 。
// # 如果对 s 的分割满足以下条件，那么我们认为它是一个 完美 分割：

// # !1.s 被分成 k 段互不相交的子字符串。
// # !2.每个子字符串长度都 至少 为 minLength 。
// # !3.每个子字符串的第一个字符都是一个 质数 数字，最后一个字符都是一个 非质数 数字。
// # 质数数字为 '2' ，'3' ，'5' 和 '7' ，剩下的都是非质数数字。

// # 请你返回 s 的 完美 分割数目。由于答案可能很大，请返回答案对 1e9 + 7 取余 后的结果。
// # 一个 子字符串 是字符串中一段连续字符串序列。

// # !k,len(s),minLength<=1000
// # 期望是O(n^2)的解法 所以需要优化掉dp范围转移的复杂度

const isPrime = new Uint8Array(10)
for (const num of [2, 3, 5, 7]) isPrime[num] = 1
const MOD = 1e9 + 7

// 这样写是O(n^3)的(压缩维度之后还能更快)
// 需要前缀和优化掉dp转移的复杂度
function beautifulPartitions(s: string, k: number, minLength: number): number {
  const n = s.length
  // !dp[count][index]表示前index个字符分成count段的方案数 (因为要用前缀和优化index维度 所以把index维度作为第二维度更加方便)
  const dp = Array.from({ length: k + 1 }, () => new Uint32Array(n + 1))
  const nums = s.split('').map(Number)
  dp[0][0] = 1

  for (let count = 1; count <= k; count++) {
    for (let index = 1; index <= n; index++) {
      // i为结尾时
      if (isPrime[nums[index - 1]]) continue
      for (let len = minLength; len <= index; len++) {
        if (isPrime[nums[index - len]]) {
          dp[count][index] += dp[count - 1][index - len]
          dp[count][index] %= MOD
        }
      }
    }
  }

  return dp[k][n]
}

// !注意到最内层的转移可以前缀和优化 => dp由一连串的index转移过来 所以考虑把index作为第二维度遍历
function beautifulPartitions2(s: string, k: number, minLength: number): number {
  const n = s.length
  // !dp[count][index]表示前index个字符分成count段的方案数 (因为要用前缀和优化index维度 所以把index维度作为第二维度更加方便)
  const dp = Array.from({ length: k + 1 }, () => new Uint32Array(n + 1))
  const nums = s.split('').map(Number)
  dp[0][0] = 1

  for (let c = 1; c <= k; c++) {
    const dpSum = new Uint32Array(n + 1)
    for (let i = 1; i <= n; i++) {
      dpSum[i] = dpSum[i - 1] + (isPrime[nums[i - 1]] ? dp[c - 1][i - 1] : 0)
      dpSum[i] %= MOD
    }

    for (let i = 1; i <= n; i++) {
      if (isPrime[nums[i - 1]]) continue
      dp[c][i] = dpSum[Math.max(0, i - minLength + 1)] // !优化掉了内层求区间和的循环
    }
  }

  return dp[k][n]
}

console.log(beautifulPartitions2('23542185131', 3, 2)) // 3

export {}
