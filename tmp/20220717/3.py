from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def smallestTrimmedNumbers(self, nums: List[str], queries: List[List[int]]) -> List[int]:
        res = []
        # 裁剪 到剩下 最右边 trimi 个数位
        for k, trim in queries:
            copy = nums[:]
            sl = []
            for i in range(len(copy)):
                copy[i] = copy[i][-trim:]
                sl.append((int(copy[i]), i))
            sl.sort()
            res.append(sl[k - 1][1])
        return res


print(
    Solution().smallestTrimmedNumbers(
        nums=["102", "473", "251", "814"], queries=[[1, 1], [2, 3], [4, 2], [1, 2]]
    )
)
