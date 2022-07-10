# 求排列数
from functools import lru_cache


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


def A(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return (fac(n) * ifac(n - k)) % MOD


if __name__ == "__main__":
    print(A(n=4, k=4))
    print(A(n=4, k=5))
