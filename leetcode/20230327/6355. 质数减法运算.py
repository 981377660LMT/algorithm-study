# 给你一个下标从 0 开始的整数数组 nums ，数组长度为 n 。
# 你可以执行无限次下述运算：
# 选择一个之前未选过的下标 i ，并选择一个 严格小于 nums[i] 的质数 p ，从 nums[i] 中减去 p 。
# !如果你能通过上述运算使得 nums 成为严格递增数组，则返回 true ；否则返回 false 。
# 严格递增数组 中的每个元素都严格大于其前面的元素。


from bisect import bisect_left
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

    def getPrimeFactors(self, n: int) -> Counter[int]:
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


E = EratosthenesSieve(1010)
P = E.getPrimes()


class Solution:
    def primeSubOperation(self, nums: List[int]) -> bool:
        assign = nums[:]
        pre = 0
        for i, x in enumerate(nums):
            # !x-p>pre的最大质数
            pos = bisect_left(P, x - pre) - 1
            assign[i] = x - (P[pos] if pos >= 0 else 0)
            pre = assign[i]
        return all(x < y for x, y in zip(assign, assign[1:]))


assert not Solution().primeSubOperation([5, 8, 3])
