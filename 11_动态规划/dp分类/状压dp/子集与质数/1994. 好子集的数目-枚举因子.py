# !如果 nums 的一个子集中，所有元素的乘积可以用若干个 `互不相同的质数` 相乘得到
# 那么我们称它为 好子集 。
# 请你返回 nums 中不同的 好 子集的数目对 109 + 7 取余 的结果。（可以取空集）

# 1 <= nums.length <= 1e5
# 1 <= nums[i] <= 30，小于等于30的质数正好是10个，暗示状压

# 每个质数p只能在好子集中出现0或1次，对应着选或不选
# 遍历可能作为好子集元素的数(一定是0-30的素数)，乘积做组合，求每种情况的频率
# 看到限制里面数字最大也不会超过 30 ，立刻想到暴力。把`组合全部弄出来`，每种组合的次数就是数字出现次数的乘积。
# !需要小心的是，[1, 1, 1] 这种是不算的(其余的数是空集),并且1的集合需要特殊处理.


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
    def numberOfGoodSubsets(self, nums: List[int]) -> int:
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
        return (res - 1) * pow(2, counter[1], MOD) % MOD  # !去掉其余数都不取的情况


print(Solution().numberOfGoodSubsets(nums=[1, 2, 3, 4]))
# 输出：6
# 解释：好子集为：
# - [1,2]：乘积为 2 ，可以表示为质数 2 的乘积。
# - [1,2,3]：乘积为 6 ，可以表示为互不相同的质数 2 和 3 的乘积。
# - [1,3]：乘积为 3 ，可以表示为质数 3 的乘积。
# - [2]：乘积为 2 ，可以表示为质数 2 的乘积。
# - [2,3]：乘积为 6 ，可以表示为互不相同的质数 2 和 3 的乘积。
# - [3]：乘积为 3 ，可以表示为质数 3 的乘积。
