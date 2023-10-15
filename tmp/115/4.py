from functools import lru_cache
from math import comb
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def countSubMultisets(self, nums: List[int], l: int, r: int) -> int:
        @lru_cache(None)
        def dfs(index: int, curSum: int) -> int:
            if curSum > r:
                return 0
            if index == m:
                return 1 if l <= curSum <= r else 0

            # 选几个当前的数
            res = 0
            curNum, curCount = elements[index]
            for i in range(curCount + 1):
                sum_ = i * curNum
                if curSum + sum_ > r:
                    break
                res += dfs(index + 1, curSum + sum_)
                res %= MOD
            return res

        counter = Counter(nums)
        elements = [(num, count) for num, count in counter.items() if num != 0]
        elements.sort(reverse=True)
        m = len(elements)
        res = dfs(0, 0)
        dfs.cache_clear()
        res *= (counter[0]) + 1
        return res % MOD


# nums = [1,2,2,3], l = 6, r = 6
print(Solution().countSubMultisets([1, 2, 2, 3], 6, 6))
# nums = [2,1,4,2,7], l = 1, r = 5
print(Solution().countSubMultisets([2, 1, 4, 2, 7], 1, 5))
# nums = [1,2,1,3,5,2], l = 3, r = 5
print(Solution().countSubMultisets([1, 2, 1, 3, 5, 2], 3, 5))
# [0,0,1,2,3]
# 2
# 3
print()

print(Solution().countSubMultisets([0, 0, 1, 2, 3], 2, 3))
print(Solution().countSubMultisets([0, 0], 0, 0))
