from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


@lru_cache(None)
def getSteps(x: int) -> int:
    if x == 0:
        return INF
    if x == 1:
        return 0
    return 1 + getSteps(x.bit_count())


STEPS = [getSteps(i) for i in range(1000)]


class Solution:
    def countKReducibleNumbers(self, s: str, k: int) -> int:
        n = len(s)
        nums = [int(c) for c in s]
        validCounts = set(c for c in range(1, n + 1) if STEPS[c] <= k - 1)

        @lru_cache(None)
        def dfs(pos: int, isLimit: bool, isNum: bool, count: int) -> int:
            if pos == n:
                return int(isNum and count in validCounts)
            res = 0
            upperLimit = nums[pos] if isLimit else 1
            for digit in [0, 1]:
                if isLimit and digit > upperLimit:
                    continue
                nextLimit = isLimit and (digit == upperLimit)
                nextIsNum = isNum or (digit != 0)
                nextCount = count + digit if nextIsNum else count
                res = (res + dfs(pos + 1, nextLimit, nextIsNum, nextCount)) % MOD
            return res

        res = dfs(0, True, False, 0)
        dfs.cache_clear()
        ones = s.count("1")
        if ones in validCounts and STEPS[ones] <= k - 1:
            res -= 1
        return res % MOD
