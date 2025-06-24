# 3589. 计数质数间隔平衡子数组
# https://leetcode.cn/problems/count-prime-gap-balanced-subarrays/description/
# 类似 3578. 统计极差最大为 K 的分割方式数
#
# 给定一个整数数组 nums 和一个整数 k。
#
# 子数组 被称为 质数间隔平衡，如果：
#
# !其包含 至少两个质数，并且
# !该 子数组 中 最大 和 最小 质数的差小于或等于 k。
# 返回 nums 中质数间隔平衡子数组的数量。
#
# 注意：
#
# 子数组 是数组中连续的 非空 元素序列。
# 质数是大于 1 的自然数，它只有两个因数，即 1 和它本身。
#
# !去除“差小于或等于 k”

from typing import List


from collections import Counter


from sortedcontainers import SortedList


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "minPrime"  # 每个数的最小质因数

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
        self.minPrime = minPrime

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self.minPrime[n] == n

    def getPrimeFactors(self, n: int) -> "Counter[int]":
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self.minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self.minPrime) if i >= 2 and i == x]


E = EratosthenesSieve(int(1e5))


class Solution:
    def primeSubarray(self, nums: List[int], k: int) -> int:
        sl1, sl2 = SortedList(), SortedList()
        l1, l2 = 0, 0
        res = 0
        for r, x in enumerate(nums):
            if E.isPrime(x):
                sl1.add(x)
                sl2.add(x)
            while l1 <= r and len(sl1) > 1 and sl1[-1] - sl1[0] > k:  # type: ignore
                if E.isPrime(nums[l1]):
                    sl1.remove(nums[l1])
                l1 += 1
            while l2 <= r and len(sl2) > 1:
                if E.isPrime(nums[l2]):
                    sl2.remove(nums[l2])
                l2 += 1
            res += l2 - l1
        return res
