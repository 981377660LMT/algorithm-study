# O(nlogv)解法
# !将数组切分成若干个子数组，
# !使得每个子数组最左边和最右边数字的最大公约数大于 1，求解最少能切成多少个子数组。
# !n<=1e5
# !nums[i]<=1e6
from collections import defaultdict
from typing import DefaultDict, List

INF = int(1e9)


def splitArray(nums: List[int]) -> int:
    """
    dp[i][p]表示前i个数中,能被p整除的数字能划分的最小段数
    !每一个新的质数可以继承到之前的质数的dp值,或者新开一组
    !ndp[p]=min(dp[p],preMin+1)

    滚动dp优化一个维度
    """
    dp, preMin = defaultdict(lambda: INF), 0
    for num in nums:
        ndp, curMin = dp, INF
        for p in E.getPrimeFactors(num):
            ndp[p] = min(ndp[p], preMin + 1)  # !继承或者新加一组
            curMin = min(curMin, ndp[p])
        dp, preMin = ndp, curMin
    return preMin


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

    def getPrimeFactors(self, n: int) -> "DefaultDict[int, int]":
        """求n的质因数分解 O(logn)"""
        res, f = defaultdict(int), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


E = EratosthenesSieve(int(1e6 + 10))


assert splitArray(nums=[2, 6, 3, 4, 3]) == 2
assert splitArray(nums=[2, 3, 5, 7]) == 4
