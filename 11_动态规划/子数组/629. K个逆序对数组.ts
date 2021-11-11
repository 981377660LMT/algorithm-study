/**
 *
 * @param n  有包含从 1 到 n 的数字
 * @param k  恰好拥有 k 个逆序对的不同的数组的个数
 *
 * 考虑加入的第i个元素对结果的贡献
 * (1)∗(1+x)∗(1+x+x2)∗……∗(1+x+x2+……+xn−1)
 */
function kInversePairs(n: number, k: number): number {
  const MOD = 10 ** 9 + 7

  let contribute = [1] // 初始项x^0
  for (let i = 1; i < n + 1; i++) {
    contribute = convolve(contribute, Array(i).fill(1)).slice(0, 1001)
  }

  return k < contribute.length ? contribute[k] : 0

  function convolve(volume: number[], kernel: number[]): number[] {
    if (volume.length === 0 || kernel.length === 0) throw new Error('数组不能为空')

    const res = Array(volume.length + kernel.length - 1).fill(0)
    for (let i = 0; i < volume.length; i++) {
      for (let j = 0; j < kernel.length; j++) {
        res[i + j] = (res[i + j] + volume[i] * kernel[j]) % MOD
      }
    }

    return res
  }
}

// O(nk)
// 考虑最大的数放在哪个位置
// https://leetcode-cn.com/problems/k-inverse-pairs-array/solution/python-ji-yi-hua-di-gui-by-himymben-bujz/
function kInversePairs2(n: number, k: number): number {
  const MOD = 10 ** 9 + 7
  if (k === 0) return 1

  // n个数 k个逆序对
  const dp = Array.from<unknown, number[]>({ length: n + 1 }, () => Array(k + 1).fill(0))

  for (let i = 0; i < n + 1; i++) {
    dp[i][0] = 1
  }

  for (let i = 1; i < n + 1; i++) {
    for (let j = 1; j < k + 1; j++) {
      dp[i][j] = dp[i][j - 1] + dp[i - 1][j] - (j >= i ? dp[i - 1][j - i] : 0)
      dp[i][j] = (dp[i][j] + MOD) % MOD
    }
  }

  return dp[n][k]
}

// 考虑最大的数放在哪个位置
// # dp[x][y] = dp[x-1][y-(x-1)] + ... + dp[x-1][y]
// # dp[x][y-1] = dp[x-1][y-1-(x-1)] + ... + dp[x-1][y-1]
// # 1-2得 dp[x][y] - dp[x][y-1] = -dp[x-1][y-x] + dp[x-1][y]
// # dp[x][y] = dp[x][y-1] + dp[x-1][y] - dp[x-1][y-x]

console.log(kInversePairs2(3, 1))

// export {}
