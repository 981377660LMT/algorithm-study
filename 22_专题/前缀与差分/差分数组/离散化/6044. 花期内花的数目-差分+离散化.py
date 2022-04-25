from bisect import bisect_right
from collections import defaultdict
from itertools import accumulate
from typing import List

# 10^9值域 差分数组


class Solution:
    def fullBloomFlowers(self, flowers: List[List[int]], persons: List[int]) -> List[int]:
        # 开map即可
        diff = defaultdict(int)
        for left, right in flowers:
            diff[left] += 1
            diff[right + 1] -= 1

        # 离散化的keys、原数组前缀和
        keys = sorted(diff)
        preSum = [0] + list(accumulate(diff[key] for key in keys))
        return [preSum[bisect_right(keys, p)] for p in persons]


print(
    Solution().fullBloomFlowers(flowers=[[1, 6], [3, 7], [9, 12], [4, 13]], persons=[2, 3, 7, 11])
)
