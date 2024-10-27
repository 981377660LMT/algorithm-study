# 3337. 字符串转换后的长度 II
# https://leetcode.cn/problems/total-characters-in-string-after-transformations-ii/description/
# 给你一个由小写英文字母组成的字符串 s，一个整数 t 表示要执行的 转换 次数，以及一个长度为 26 的数组 nums。每次 转换 需要根据以下规则替换字符串 s 中的每个字符：
# 将 s[i] 替换为字母表中后续的 nums[s[i] - 'a'] 个连续字符。例如，如果 s[i] = 'a' 且 nums[0] = 3，则字符 'a' 转换为它后面的 3 个连续字符，结果为 "bcd"。
# 如果转换超过了 'z'，则 回绕 到字母表的开头。例如，如果 s[i] = 'y' 且 nums[24] = 3，则字符 'y' 转换为它后面的 3 个连续字符，结果为 "zab"。
# 返回 恰好 执行 t 次转换后得到的字符串的 长度。
# 由于答案可能非常大，返回其对 109 + 7 取余的结果。
# !每个字符互不影响，每个字符的转换规则不同
#
# dp[i][k] 表示字符i转换k次后的长度
# !dp[i][k] -> dp[i+1][k-1] + dp[i+2][k-1] + ... + dp[(i+nums[i])%26][k-1]

from typing import List

MOD = int(1e9 + 7)


class Solution:
    def lengthAfterTransformations(self, s: str, t: int, nums: List[int]) -> int:
        T = [[0] * 26 for _ in range(26)]
        for i in range(26):
            for j in range(1, nums[i] + 1):
                T[i][(i + j) % 26] = 1

        init = [[1] for _ in range(26)]
        resT = matpow(T, t, MOD)
        res = matmul(resT, init, MOD)
        ords = [ord(c) - ord("a") for c in s]
        return sum(res[ord_][0] for ord_ in ords) % MOD


def matmul(mat1: List[List[int]], mat2: List[List[int]], mod: int) -> List[List[int]]:
    """矩阵相乘"""
    i_, j_, k_ = len(mat1), len(mat2[0]), len(mat2)
    res = [[0] * j_ for _ in range(i_)]
    for i in range(i_):
        for k in range(k_):
            if mat1[i][k]:
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
