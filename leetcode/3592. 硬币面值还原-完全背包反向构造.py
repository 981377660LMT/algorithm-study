# 3592. 硬币面值还原-完全背包反向构造
# https://leetcode.cn/problems/inverse-coin-change/
# 给你一个 从 1 开始计数 的整数数组 numWays，其中 numWays[i] 表示使用某些 固定 面值的硬币（每种面值可以使用无限次）凑出总金额 i 的方法数。
# 每种面值都是一个 正整数 ，并且其值 最多 为 numWays.length。
# 然而，具体的硬币面值已经 丢失 。你的任务是还原出可能生成这个 numWays 数组的面值集合。
# 返回一个按从小到大顺序排列的数组，其中包含所有可能的 唯一 整数面值。
# 如果不存在这样的集合，返回一个 空 数组。
#
# 对每一个价值，要么有要么没有
# 没有 => dp[i] = numsWays[i]
# 有 => dp[i] + 1 = numsWays[i] (最小的那个第一次贡献)

from typing import List


class Solution:
    def findCoins(self, numWays: List[int]) -> List[int]:
        n = len(numWays)
        dp = [1] + [0] * n
        res = []
        for i, v in enumerate(numWays, 1):
            if dp[i] == v:
                continue
            if dp[i] + 1 != v:
                return []
            res.append(i)
            for j in range(i, n + 1):
                dp[j] += dp[j - i]
        return res
