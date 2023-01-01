from math import floor
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数数组 nums ，对 nums 所有元素求积之后，找出并返回乘积中 不同质因数 的数目。

# 注意：

# 质数 是指大于 1 且仅能被 1 及自身整除的数字。
# 如果 val2 / val1 是一个整数，则整数 val1 是另一个整数 val2 的一个因数。


def getPrimeFactors(n: int) -> Counter:
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


class Solution:
    def distinctPrimeFactors(self, nums: List[int]) -> int:
        res = Counter()
        for num in nums:
            res += getPrimeFactors(num)
        return len(res)
