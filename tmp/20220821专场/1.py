from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def numberOfPairs(self, nums: List[int]) -> int:
        """忘记取模WA了 以后看到取模先写res%MOD"""

        def cal(x: int) -> int:
            s1, s2 = str(x), str(x)[::-1]
            return int(s1) - int(s2)

        adjMap = defaultdict(list)
        for num in nums:
            adjMap[cal(num)].append(num)
        res = 0
        for v in adjMap.values():
            if len(v) > 1:
                res += len(v) * (len(v) - 1) // 2
                res %= MOD
        return res % MOD
