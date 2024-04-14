from functools import lru_cache
from math import floor, gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


# 超级丑数
MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 coins 表示不同面额的硬币，另给你一个整数 k 。

# 你有无限量的每种面额的硬币。但是，你 不能 组合使用不同面额的硬币。


# 返回使用这些硬币能制造的 第 kth 小 金额。

# 容斥原理


def count(upper: int, nums: List[int]) -> int:
    """[1, upper]中能被primes中的至少一个数整除的数的个数"""
    m = len(nums)
    res = 0
    for state in range(1, (1 << m)):
        mul = 1
        for i in range(m):
            if state & (1 << i):
                gcd_ = gcd(mul, nums[i])
                mul *= nums[i] // gcd_

        if state.bit_count() & 1:
            res += upper // mul
        else:
            res -= upper // mul
    return res


class Solution:
    def findKthSmallest(self, coins: List[int], k: int) -> int:
        def check(mid: int) -> bool:
            return count(mid, coins) >= k

        left, right = 1, int(1e12)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


# coins = [3,6,9], k = 3
print(Solution().findKthSmallest([3, 6, 9], 3))
# coins = [5,2], k = 7
print(Solution().findKthSmallest([5, 2], 7))
# print(count(10, [2, 5]), 111)
