import re
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumSum(self, nums: List[int]) -> int:
        def cal(num: int) -> int:
            s = str(num)
            res = 0
            for n in s:
                res += int(n)
            return res

        adjMap = defaultdict(list)
        for num in nums:
            sum_ = cal(num)
            adjMap[sum_].append(num)
        for v in adjMap.values():
            v.sort()

        res = -1
        for v in adjMap.values():
            if len(v) <= 1:
                continue
            res = max(res, v[-1] + v[-2])

        return res


print(Solution().maximumSum(nums=[18, 43, 36, 13, 7]))
