# 100462. 找到初始输入字符串 II-前缀和优化dp(分组背包)
# https://leetcode.cn/problems/find-the-original-typed-string-ii/description/
# Alice 正在她的电脑上输入一个字符串。但是她打字技术比较笨拙，她 可能 在一个按键上按太久，导致一个字符被输入 多次 。
# 给你一个字符串 word ，它表示 最终 显示在 Alice 显示屏上的结果。同时给你一个 正 整数 k ，表示一开始 Alice 输入字符串的长度 至少 为 k 。
# 请你返回 Alice 一开始可能想要输入字符串的总方案数。
# 由于答案可能很大，请你将它对 109 + 7 取余 后返回。
#
# !减去长度<k的方案数.

from functools import reduce
from itertools import accumulate, groupby
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
    def possibleStringCount(self, word: str, k: int) -> int:
        def caclBad() -> int:
            """前i个物品, 每个物品有lens[i]个, 选<target个物品的方案数."""
            if len(lens) >= k:
                return 0
            target = k - len(lens) - 1
            dp = groupKnapsackCount(lens, target, mod=MOD)
            return sum(dp)

        lens = [len(list(group)) for _, group in groupby(word)]
        all_ = reduce(lambda x, y: x * y, lens, 1)
        bad = caclBad()
        return int((all_ - bad) % MOD)
