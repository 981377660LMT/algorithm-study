/* eslint-disable arrow-body-style */
// https://leetcode.cn/problems/string-transformation/
// 字符串轮转
// 给你两个长度都为 n 的字符串 s 和 t 。你可以对字符串 s 执行以下操作：

// 将 s 长度为 l （0 < l < n）的 后缀字符串 删除，并将它添加在 s 的开头。
// 比方说，s = 'abcd' ，那么一次操作中，你可以删除后缀 'cd' ，并将它添加到 s 的开头，得到 s = 'cdab' 。
// 给你一个整数 k ，请你返回 恰好 k 次操作将 s 变为 t 的方案数。

// 由于答案可能很大，返回答案对 1e9 + 7 取余 后的结果。
//
// 2 <= s.length <= 5 * 105
// 1 <= k <= 1015
// s.length == t.length
// s 和 t 都只包含小写英文字母。
//
// !记dp[i][0/1]为 `i` 次操作后 `等于/不等于t` 的方案数，count 为 `t` 在 `s` 的循环轮转中出现的次数
// !dp[i][0] = dp[i-1][0]*(count-1) + dp[i-1][1]*count
// !dp[i][1] = dp[i-1][0]*(n-count) + dp[i-1][1]*(n-count-1)
// 因此有状态转移方程
// dp[i] = T*dp[i-1],
//
// T = [
//    [count-1, count],
//    [n-count, n-count-1]
//   ]

import { indexOfAll } from '../../../../17_模式匹配/kmp/kmp'
import { matMul, matPow } from '../matqpow'

const MOD = 1e9 + 7

function numberOfWays(s: string, t: string, k: number): number {
  // 统计t在s的循环轮转中出现的次数,不包含(s==t)(即在(s+s)[1:]中出现的次数)
  const countTInCyclicS = (s: string, t: string) => {
    return indexOfAll(s + s, t, 1).length
  }

  const n = s.length
  const count = countTInCyclicS(s, t)
  const init = s === t ? [[1], [0]] : [[0], [1]]
  const T = [
    [count - 1, count],
    [n - count, n - count - 1]
  ]
  const resT = matPow(T, k, MOD)
  const res = matMul(resT, init, MOD)
  return res[0][0] % MOD
}
