from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def findCoins(self, numWays: List[int]) -> List[int]:
        n = len(numWays)
        dp = [0] * (n + 1)
        dp[0] = 1
        res = []
        for i in range(1, n + 1):
            if dp[i] > numWays[i - 1]:
                return []
            elif dp[i] < numWays[i - 1]:
                if numWays[i - 1] - dp[i] != 1:
                    return []
                res.append(i)
                for j in range(i, n + 1):
                    dp[j] += dp[j - i]
        if numWays != dp[1:]:
            return []
        return res
