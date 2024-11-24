from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minArraySum(self, nums: List[int], k: int, op1: int, op2: int) -> int:
        sum_ = sum(nums)
        dp = [[INF] * (op2 + 1) for _ in range(op1 + 1)]
        dp[0][0] = sum_
        for v in nums:
            select = []
            select.append((0, 0, v))
            v1 = (v + 1) // 2
            select.append((1, 0, v1))
            if v >= k:
                v2 = v - k
                select.append((0, 1, v2))
                v4 = (v2 + 1) // 2
                select.append((1, 1, v4))
            v3 = (v + 1) // 2
            if v3 >= k:
                v3 -= k
                select.append((1, 1, v3))
            ndp = [[INF] * (op2 + 1) for _ in range(op1 + 1)]
            for i in range(op1 + 1):
                for j in range(op2 + 1):
                    if dp[i][j] == INF:
                        continue
                    for a, b, c in select:
                        ni, nj = i + a, j + b
                        if ni <= op1 and nj <= op2:
                            ndp[ni][nj] = min2(ndp[ni][nj], dp[i][j] - v + c)
            dp = ndp
        return min(min(row) for row in dp)
