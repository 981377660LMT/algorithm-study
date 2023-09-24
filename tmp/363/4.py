from functools import lru_cache

from typing import DefaultDict, FrozenSet, List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 1 开始、由 n 个整数组成的数组。

# 如果一组数字中每对元素的乘积都是一个完全平方数，则称这组数字是一个 完全集 。

# 下标集 {1, 2, ..., n} 的子集可以表示为 {i1, i2, ..., ik}，我们定义对应该子集的 元素和 为 nums[i1] + nums[i2] + ... + nums[ik] 。

# 返回下标集 {1, 2, ..., n} 的 完全子集 所能取到的 最大元素和 。


# 完全平方数是指可以表示为一个整数和其自身相乘的数。

from collections import Counter
from typing import List


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "_minPrime"  # 每个数的最小质因数

    def __init__(self, maxN: int):
        """预处理 O(nloglogn)"""
        minPrime = list(range(maxN + 1))
        upper = int(maxN**0.5) + 1
        for i in range(2, upper):
            if minPrime[i] < i:
                continue
            for j in range(i * i, maxN + 1, i):
                if minPrime[j] == j:
                    minPrime[j] = i
        self._minPrime = minPrime

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self._minPrime[n] == n

    def getPrimeFactors(self, n: int) -> "DefaultDict":
        """求n的质因数分解 O(logn)"""
        res, f = defaultdict(int), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


E = EratosthenesSieve(int(1e4 + 10))
P = [E.getPrimeFactors(i) for i in range(int(1e4 + 10))]
for d in P:
    for k, v in d.items():
        d[k] = v % 2


class Solution:
    def maximumSum(self, nums: List[int]) -> int:
        #  每个数都为完全平方数
        if len(nums) == 1:
            return nums[0]

        nums = [0] + nums
        res = max(nums)

        @lru_cache(None)
        def dfs(value: int, select: FrozenSet[int]) -> int:
            """如果当前因数为偶数,则必须在select中,否则必须不在select中"""
            if value == len(nums):
                return 0

            # 不选
            res = dfs(value + 1, select)
            # 选
            factors = P[value]
            ok = True
            for factor, count in factors.items():
                if count == 1 and factor not in select:
                    ok = False
                    break
            if ok:
                select = select.union(factors.keys())
                cand = dfs(value + 1, select) + nums[value]
                if cand > res:
                    res = cand
            return res

        res = dfs(1, frozenset())
        dfs.cache_clear()
        return res


nums = [5, 10, 3, 10, 1, 13, 7, 9, 4]
nums = [8, 7, 3, 5, 7, 2, 4, 9]
print(Solution().maximumSum(nums))
