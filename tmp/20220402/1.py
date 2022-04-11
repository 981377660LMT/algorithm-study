from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minBitFlips(self, start: int, goal: int) -> int:
        res = 0
        for i in range(34):
            if (start >> i) & 1 != (goal >> i) & 1:
                res += 1
        return res


print(Solution().minBitFlips(10, 7))
print(Solution().minBitFlips(3, 4))
