from typing import List
from collections import Counter
from functools import lru_cache
from math import floor

MOD = int(1e9 + 7)


@lru_cache(None)
def getPrimeFactors(n: int) -> Counter:
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = floor(n ** 0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


@lru_cache(None)
def dfs(index: int, remain: int) -> int:
    """remain个物品放到index个背包的放法"""
    if remain < 0:
        return 0
    if index == 0:
        return 1 if remain == 0 else 0
    res = 0
    for select in range(remain + 1):
        res += dfs(index - 1, remain - select)
        res %= MOD
    return res


@lru_cache(None)
def cal(bag: int, product: int) -> int:
    """乘积为product的数分解质因数放到bag个背包的放法,对每个素因子分析,相乘即可"""
    res = 1
    for count in getPrimeFactors(product).values():
        res *= dfs(bag, count)
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

