# 873. 最长的斐波那契子序列的长度
# https://leetcode.cn/problems/length-of-longest-fibonacci-subsequence/description/
# 给定一个严格递增的正整数数组形成序列 arr ，找到 arr 中最长的斐波那契式的子序列的长度。
# 如果一个不存在，返回  0 。
# n<=1000


from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def lenLongestFibSubseq(self, arr: List[int]) -> int:
        """dp[j][i]表示arr[j]和arr[i]结尾的最长斐波那契子序列的长度."""
        n = len(arr)
        mp = {v: i for i, v in enumerate(arr)}
        dp = [[0] * n for _ in range(n)]
        res = 0
        for i, x in enumerate(arr):
            for j in range(n - 1, -1, -1):
                if arr[j] * 2 <= x:
                    break
                if x - arr[j] in mp:
                    k = mp[x - arr[j]]
                    dp[j][i] = max2(dp[k][j] + 1, 3)
                    res = max2(res, dp[j][i])
        return res
