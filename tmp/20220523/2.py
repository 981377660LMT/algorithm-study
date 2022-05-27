from itertools import accumulate
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumBags(self, capacity: List[int], rocks: List[int], additionalRocks: int) -> int:
        diff = [cap - cur for cap, cur in zip(capacity, rocks)]
        diff.sort()
        preSum = [0] + list(accumulate(diff))

        def check(mid: int) -> bool:
            return preSum[mid] <= additionalRocks

        left, right = 0, len(rocks)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right  # 最右二分

    # 更好的写法
    def maximumBags2(self, capacity: List[int], rocks: List[int], additionalRocks: int) -> int:
        diff = [cap - cur for cap, cur in zip(capacity, rocks)]
        diff.sort()
        res = 0
        for num in diff:
            if num > additionalRocks:
                break
            res += 1
            additionalRocks -= num
        return res
