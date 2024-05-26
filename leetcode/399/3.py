from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个整数数组 nums1 和 nums2，长度分别为 n 和 m。同时给你一个正整数 k。

# 如果 nums1[i] 可以被 nums2[j] * k 整除，则称数对 (i, j) 为 优质数对（0 <= i <= n - 1, 0 <= j <= m - 1）。


# 返回 优质数对 的总数。


def getFactors(n: int) -> List[int]:
    """n 的所有因数 O(sqrt(n))"""
    if n <= 0:
        return []
    small, big = [], []
    upper = int(n**0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


class Solution:
    def numberOfPairs(self, nums1: List[int], nums2: List[int], k: int) -> int:
        res = 0
        counter2 = dict()
        for v2 in nums2:
            counter2[v2] = counter2.get(v2, 0) + 1
        counter1 = Counter(nums1)
        for v1, cnt in counter1.items():
            if v1 % k == 0:
                div = v1 // k
                factors = getFactors(div)
                for factor in factors:
                    res += counter2.get(factor, 0) * cnt
        return res
