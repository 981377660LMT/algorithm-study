# 求阶乘

from functools import lru_cache
from math import factorial


MOD = int(1e9 + 7)


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


if __name__ == '__main__':
    print(fac(10))
    # 不要用这个 无法取模容易超时
    print(factorial(10))
