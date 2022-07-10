# 给定一个严格递增的正整数数组形成序列 arr ，找到 arr 中最长的斐波那契式的子序列的长度。
# 如果一个不存在，返回  0 。
from collections import defaultdict
from itertools import combinations
from typing import List


class Solution:
    def lenLongestFibSubseq(self, arr: List[int]) -> int:
        """3 <= arr.length <= 1000"""
        n = len(arr)
        indexMap = {v: i for i, v in enumerate(arr)}
        dp = defaultdict(lambda: defaultdict(lambda: 2))
        res = 2
        for j, k in combinations(range(n), 2):
            pre = arr[k] - arr[j]
            if pre in indexMap:
                i = indexMap[pre]
                if i < j:
                    dp[j][k] = max(dp[j][k], dp[i][j] + 1)
                    res = max(res, dp[j][k])

        return res if res > 2 else 0
