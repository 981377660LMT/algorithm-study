# 1155. 掷骰子等于目标和的方法数
# https://leetcode.cn/problems/number-of-dice-rolls-with-target-sum/description/
# 这里有 n 个一样的骰子，每个骰子上都有 k 个面，分别标号为 1 到 k 。
# 给定三个整数 n ,  k 和 target ，
# 返回可能的方式(从总共 k^n 种方式中)滚动骰子的数量，
# 使正面朝上的数字之和等于 target 。
# 答案可能很大，你需要对 1e9 + 7 取模 。


from itertools import accumulate
from typing import List


MOD = int(1e9 + 7)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def groupKnapsackCount(count: List[int], k: int, mod=int(1e9 + 7)) -> List[int]:
    """分组背包求方案数.
    每个物品有count[i]个, 求选择k个物品的方案数.
    """
    n = len(count)
    dp = [0] * (k + 1)
    dp[0] = 1
    for i in range(n):
        ndp = [0] * (k + 1)
        dpPresum = [0] + list(accumulate(dp))
        for j in range(k + 1):
            upper = min2(j + 1, count[i])
            ndp[j] = (dpPresum[j + 1] - dpPresum[j + 1 - upper]) % mod
        dp = ndp
    return dp


class Solution:
    def numRollsToTarget(self, n: int, k: int, target: int) -> int:
        dp = groupKnapsackCount([k] * n, target, mod=MOD)
        return dp[target - n]
