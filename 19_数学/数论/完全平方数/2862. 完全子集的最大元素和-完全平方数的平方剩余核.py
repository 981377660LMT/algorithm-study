# 2862. 完全子集的最大元素和-完全平方数的平方剩余核
# https://leetcode.cn/problems/maximum-element-sum-of-a-complete-subset-of-indices/
# 给你一个下标从 1 开始、由 n 个整数组成的数组。
# !如果一组数字中每对元素的乘积都是一个完全平方数，则称这组数字是一个 完全集 。
# 下标集 {1, 2, ..., n} 的子集可以表示为 {i1, i2, ..., ik}，
# 我们定义对应该子集的 元素和 为 nums[i1] + nums[i2] + ... + nums[ik] 。
# !返回下标集 {1, 2, ..., n} 的 完全子集 所能取到的 最大元素和 。
# 完全平方数是指可以表示为一个整数和其自身相乘的数。
# 1 <= n == nums.length <= 1e4
# 1 <= nums[i] <= 1e9


# !记kernel(x)为x中去除所有完全平方数因子后的结果，则每个组的kernal[i]都必须相同
# kernel(8)=8/4=2
# kernel(12)=12/4=3
# kernel(25)=25/25=1


from typing import List
from collections import Counter, defaultdict


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

    def getPrimeFactors(self, n: int) -> "Counter[int]":
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


E = EratosthenesSieve(int(1e5) + 10)


def getSquareFreeKernel(E: "EratosthenesSieve", n: int) -> int:
    res = 1
    for p, c in E.getPrimeFactors(n).items():
        if c & 1:
            res *= p
    return res


KERNELS = [getSquareFreeKernel(E, i) for i in range(int(1e5) + 10)]


class Solution:
    def maximumSum(self, nums: List[int]) -> int:
        mp = defaultdict(int)
        n = len(nums)
        for i in range(n):
            mp[KERNELS[i + 1]] += nums[i]
        return max(mp.values())
