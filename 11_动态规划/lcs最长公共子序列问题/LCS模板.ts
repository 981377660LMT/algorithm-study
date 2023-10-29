import { LCSBit } from './最长公共子序列/位运算加速LCS'
import { getLCSFast } from './最长公共子序列/位运算加速LCS-getLCS'

function LCS<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): number {
  const n1 = arr1.length
  const n2 = arr2.length
  const dp = new Uint16Array((n1 + 1) * (n2 + 1))
  let res = 0
  for (let i = 1; i <= n1; ++i) {
    for (let j = 1; j <= n2; ++j) {
      if (arr1[i - 1] === arr2[j - 1]) {
        dp[i * (n2 + 1) + j] = dp[(i - 1) * (n2 + 1) + j - 1] + 1
        res = Math.max(res, dp[i * (n2 + 1) + j])
      } else {
        dp[i * (n2 + 1) + j] = Math.max(dp[(i - 1) * (n2 + 1) + j], dp[i * (n2 + 1) + j - 1])
      }
    }
  }
  return res
}

function getLCS<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): T[] {
  const n1 = arr1.length
  const n2 = arr2.length
  const dp = new Uint16Array((n1 + 1) * (n2 + 1))
  const pre = new Int8Array((n1 + 1) * (n2 + 1)) // 1:左上 2:左 3:上 0:无
  for (let i = 0; i < n1 + 1; i++) {
    for (let j = 0; j < n2 + 1; j++) {
      if (i < n1 && j < n2 && arr1[i] === arr2[j]) {
        if (dp[(i + 1) * (n2 + 1) + j + 1] < dp[i * (n2 + 1) + j] + 1) {
          dp[(i + 1) * (n2 + 1) + j + 1] = dp[i * (n2 + 1) + j] + 1
          pre[(i + 1) * (n2 + 1) + j + 1] = 1
        }
      }
      if (i < n1) {
        if (dp[(i + 1) * (n2 + 1) + j] < dp[i * (n2 + 1) + j]) {
          dp[(i + 1) * (n2 + 1) + j] = dp[i * (n2 + 1) + j]
          pre[(i + 1) * (n2 + 1) + j] = 2
        }
      }
      if (j < n2) {
        if (dp[i * (n2 + 1) + j + 1] < dp[i * (n2 + 1) + j]) {
          dp[i * (n2 + 1) + j + 1] = dp[i * (n2 + 1) + j]
          pre[i * (n2 + 1) + j + 1] = 3
        }
      }
    }
  }

  const res: T[] = []
  let i = n1
  let j = n2
  while (i && j && pre[i * (n2 + 1) + j]) {
    const cur = pre[i * (n2 + 1) + j]
    if (cur === 1) {
      i--
      j--
      res.push(arr1[i])
    } else if (cur === 2) {
      i--
    } else if (cur === 3) {
      j--
    }
  }

  return res.reverse()
}

export { LCS, getLCS, LCSBit, getLCSFast }

if (require.main === module) {
  // https://leetcode.cn/problems/longest-common-subsequence/
  // eslint-disable-next-line no-inner-declarations
  function longestCommonSubsequence(text1: string, text2: string): number {
    return getLCS(text1, text2).length
  }
}
