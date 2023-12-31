#  https://leetcode.cn/problems/string-transformation/description/
#  字符串轮转
#  给你两个长度都为 n 的字符串 s 和 t 。你可以对字符串 s 执行以下操作：
#
#  将 s 长度为 l （0 < l < n）的 后缀字符串 删除，并将它添加在 s 的开头。
#  比方说，s = 'abcd' ，那么一次操作中，你可以删除后缀 'cd' ，并将它添加到 s 的开头，得到 s = 'cdab' 。
#  给你一个整数 k ，请你返回 恰好 k 次操作将 s 变为 t 的方案数。
#
#  由于答案可能很大，返回答案对 1e9 + 7 取余 后的结果。
#
#  2 <= s.length <= 5 * 105
#  1 <= k <= 1015
#  s.length == t.length
#  s 和 t 都只包含小写英文字母。
#
#  !记dp[i][0/1]为 `i` 次操作后 `等于/不等于t` 的方案数，count 为 `t` 在 `s` 的循环轮转中出现的次数
#  !dp[i][0] = dp[i-1][0]*(count-1) + dp[i-1][1]*count
#  !dp[i][1] = dp[i-1][0]*(n-count) + dp[i-1][1]*(n-count-1)
#  因此有状态转移方程
#  dp[i] = T*dp[i-1],
#
#  T = [
#    [count-1, count],
#    [n-count, n-count-1]
#  ]


from typing import List

MOD = int(1e9 + 7)


class Solution:
    def numberOfWays(self, s: str, t: str, k: int) -> int:
        n = len(s)
        allPos = indexOfAll(s + s, t, start=1)
        count = len(allPos)
        init = [[1], [0]] if s == t else [[0], [1]]
        T = [[count - 1, count], [n - count, n - count - 1]]
        resT = matpow(T, k, MOD)
        res = matmul(resT, init, MOD)
        return res[0][0]


def getNext(needle: str) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组
    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    """
    next = [0] * len(needle)
    j = 0
    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]
        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1
        next[i] = j
    return next


def indexOfAll(longer, shorter, start=0) -> List[int]:
    """kmp O(n+m)求搜索串 `longer` 中所有匹配 `shorter` 的位置."""
    if not shorter:
        return []
    if len(longer) < len(shorter):
        return []
    res = []
    next = getNext(shorter)
    hitJ = 0
    for i in range(start, len(longer)):
        while hitJ > 0 and longer[i] != shorter[hitJ]:
            hitJ = next[hitJ - 1]
        if longer[i] == shorter[hitJ]:
            hitJ += 1
        if hitJ == len(shorter):
            res.append(i - len(shorter) + 1)
            hitJ = next[hitJ - 1]
    return res


def matmul(mat1: List[List[int]], mat2: List[List[int]], mod: int) -> List[List[int]]:
    """矩阵相乘"""
    i_, j_, k_ = len(mat1), len(mat2[0]), len(mat2)
    res = [[0] * j_ for _ in range(i_)]
    for i in range(i_):
        for k in range(k_):
            for j in range(j_):
                res[i][j] = (res[i][j] + mat1[i][k] * mat2[k][j]) % mod
    return res


def matpow(base: List[List[int]], exp: int, mod: int) -> List[List[int]]:
    n = len(base)
    e = [[0] * n for _ in range(n)]
    for i in range(n):
        e[i][i] = 1
    b = [row[:] for row in base]
    while exp:
        if exp & 1:
            e = matmul(e, b, mod)
        exp >>= 1
        b = matmul(b, b, mod)
    return e
