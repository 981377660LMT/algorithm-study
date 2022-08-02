# 求阶乘

##########################################
# 1.阶乘打表

MOD = int(1e9 + 7)
fac = [1]
for i in range(1, int(2e5) + 10):
    fac.append(fac[-1] * i % MOD)

print(fac[20], fac[10])
##########################################
# 2.记忆化
from functools import lru_cache


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


##########################################
from math import factorial

if __name__ == "__main__":
    print(fac(10))
    # 不要用这个 无法取模容易超时
    print(factorial(10))
