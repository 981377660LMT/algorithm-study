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
        digits = list(map(int, s))
        valid = [c for c in range(1, n + 1) if STEPS[c] <= k - 1]
        res = 0
        for c in valid:
            count = self.solve(digits, c)
            res += count
        ones = s.count("1")
        if ones in valid:
            if STEPS[ones] <= k - 1:
                res = (res - 1) % MOD
        return res % MOD

    def solve(self, digits: List[int], c: int) -> int:
        n = len(digits)

        @lru_cache(None)
        def dfs(pos: int, count: int, isLimit: bool, isNum: bool) -> int:
            if count > c:
                return 0
            if pos == n:
                return int(isNum and count == c)
            res = 0
            upper = digits[pos] if isLimit else 1
            for digit in range(upper + 1):
                nextLimit = isLimit and (digit == upper)
                nextIsNum = isNum or digit != 0
                nextCount = count + digit if nextIsNum else count
                res = (res + dfs(pos + 1, nextCount, nextLimit, nextIsNum)) % MOD
            return res

        res = dfs(0, 0, True, False)
        dfs.cache_clear()
        return res
