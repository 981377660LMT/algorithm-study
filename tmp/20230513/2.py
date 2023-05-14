from operator import inv
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 枚举集合最小值

INV2 = pow(2, MOD - 2, MOD)


class Solution:
    def sumOfPower(self, nums: List[int]) -> int:
        nums.sort()
        cands = []
        curSum = 0
        n = len(nums)
        for i in range(1, n):
            v = nums[i]
            cur = v * v * pow(2, i - 1, MOD) % MOD
            cands.append(cur)
            curSum += cur
            curSum %= MOD

        res = 0
        for i in range(n - 1):
            min_ = nums[i]
            res += curSum * min_
            res %= MOD
            curSum -= cands[i] * pow(INV2, i, MOD)  # 这里不对
            curSum *= INV2
            curSum %= MOD

        res += sum(v * v * v for v in nums)
        res %= MOD
        return res


def bf(nums: List[int]) -> int:
    res = 0
    for sub in range(1, 1 << len(nums)):
        cur = []
        for i in range(len(nums)):
            if sub & (1 << i):
                cur.append(nums[i])
        res += min(cur) * max(cur) ** 2
    return res


# print(Solution().sumOfPower(nums=[2, 1, 4]))
# # [76,24,96,82,97]
# print(Solution().sumOfPower(nums=[76, 24, 96, 82, 97]))
# # 13928461
# print(bf([76, 24, 96, 82, 97]))
print(bf([1, 2, 3, 4]), Solution().sumOfPower([1, 2, 3, 4]))
