# 给定 n 个正整数 ai，请你输出这些数的乘积的约数个数，答案对 109+7 取模。


# 约数个数：(a1+1)(a2+1)...(ak+1)  考虑每个质因子的贡献即可

from collections import Counter
from functools import lru_cache
from math import floor


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


# 求 1-n 乘积的约数个数
n = int(input())
counter = Counter()
for _ in range(n):
    counter += getPrimeFactors(int(input()))

res = 1
for count in counter.values():
    res *= count + 1
    res %= int(1e9 + 7)

print(res)

