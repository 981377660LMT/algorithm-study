"""
1 <= queries.length <= 104
1 <= ni, ki <= 104

分解质因数+n个物品放到k个槽的方案数
"""

from typing import List
from collections import Counter
from functools import lru_cache
from math import floor

MOD = int(1e9 + 7)


@lru_cache(None)
def getPrimeFactors(
    n: int,
) -> Counter[int]:
    """返回 n 的质因子分解"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    if n > 1:
        res[n] += 1
    return res


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


def put(n: int, k: int) -> int:
    """
    n个物品放入k个槽(槽可空)的方案数
    """
    return C(n + k - 1, k - 1)


@lru_cache(None)
def cal(bag: int, product: int) -> int:
    """乘积为product的数分解质因数放到bag个背包的放法,对每个素因子分析,相乘即可"""
    res = 1
    for count in getPrimeFactors(product).values():
        res *= put(count, bag)
        res %= MOD
    return res


# 1 <= queries.length <= 104
# 1 <= ni, ki <= 104
# 把queries里的乘积“k”做素因子分解，把所有的素数因子分配到n个槽里，同一个槽内的素因子要相乘。槽允许为空，空槽是“1”。
class Solution:
    def waysToFillArray(self, queries: List[List[int]]) -> List[int]:
        return [cal(bag, product) for bag, product in queries]


print(Solution().waysToFillArray(queries=[[2, 6], [5, 1], [73, 660]]))
# 输出：[4,1,50734910]
# 解释：每个查询之间彼此独立。
# [2,6]：总共有 4 种方案得到长度为 2 且乘积为 6 的数组：[1,6]，[2,3]，[3,2]，[6,1]。
# [5,1]：总共有 1 种方案得到长度为 5 且乘积为 1 的数组：[1,1,1,1,1]。
# [73,660]：总共有 1050734917 种方案得到长度为 73 且乘积为 660 的数组。1050734917 对 109 + 7 取余得到 50734910 。
