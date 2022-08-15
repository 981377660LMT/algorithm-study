from heapq import heapify, heappop, heappush
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)

# 6022. 将数组和减半的最少操作次数-pq
class Solution:
    def halveArray(self, nums: List[int]) -> int:
        """用double没问题,都是整数除以2的幂次,浮点数可以精确表示,尾数也够"""
        target = sum(nums) / 2
        curSum = 0
        res = 0
        pq = [-n for n in nums]
        heapify(pq)
        while curSum < target:
            top = -heappop(pq)
            curSum += top / 2
            heappush(pq, -top / 2)
            res += 1
        return res

    def halveArray2(self, nums: List[int]) -> int:
        """避免使用浮点数 将每个数都乘上一个 2^20 可以证明每个数除 2 的次数不会超过 20。"""
        ...
