# 2572. 无平方子集计数
# 如果数组 nums 的子集中的元素乘积是一个 无平方因子数 ，则认为该子集是一个 无平方 子集。
# 无平方因子数 是无法被除 1 之外任何平方数整除的数字。
# 返回数组 nums 中 无平方 且 `非空` 的子集数目。因为答案可能很大，返回对 1e9 + 7 取余的结果。
# 1 <= nums.length <= 1000
# 1 <= nums[i] <= 30

from collections import Counter
from typing import List
from functools import lru_cache


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


MOD = int(1e9 + 7)
E = EratosthenesSieve(100)
N = 30
P = [x for x in range(2, N + 1) if E.isPrime(x)]  # !N以内的质数 (N=30时，质数有10个)
MP = {p: i for i, p in enumerate(P)}  # !每个质数的编号
F = [E.getPrimeFactors(x) for x in range(N + 1)]  # !N以内的数的质因数分解
FMAX = [max(fs.values(), default=0) for fs in F]  # !每个数的质因数分解中最多的质因数的个数
MASK = [sum(1 << MP[p] for p in f) for f in F]  # !每个数的质因数分解对应的mask


class Solution:
    def squareFreeSubsets(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, visited: int) -> int:
            if index == n:
                return 1

            # 不选
            res = dfs(index + 1, visited)
            # 选
            if visited & MASK[keys[index]] == 0 and FMAX[keys[index]] <= 1:
                res += dfs(index + 1, visited | MASK[keys[index]]) * counter[keys[index]]
            return res % MOD

        counter = Counter(nums)
        keys = sorted(set(num for num in nums if num != 1))
        n = len(keys)
        res = dfs(0, 0)
        dfs.cache_clear()
        return (res * pow(2, counter[1], MOD) - 1) % MOD  # 不能取空集
