# 伯努利数 (BernoulliNumber)
# n<=5e5
# https://tjkendev.github.io/procon-library/python/math/bernoulli-number.html


# calculate bernoulli numbers B_0 ~ B_p on modulo mod
from typing import List, Tuple


# O(p^2)求伯努利数,返回B_0 ~ B_p
def bernoulli(p: int, mod: int) -> Tuple[List[int], List[int]]:
    rev = [1] * (p + 2)
    cur = 1
    for i in range(1, p + 1):
        cur = cur * i % mod
        rev[i + 1] = cur
    cur = pow(cur * (p + 1), mod - 2, mod)
    for i in range(p + 1, -1, -1):
        rev[i] = (rev[i] * cur) % mod
        cur = cur * i % mod

    B = [0] * (p + 1)
    B[0] = 1
    if p > 1:
        B[1] = (-rev[2]) % mod
    for i in range(2, p + 1, 2):
        v = 0
        tmp = 1
        for j in range(i):
            v = (v + tmp * B[j]) % mod
            tmp = ((i + 1 - j) * rev[j + 1]) * tmp % mod
        B[i] = (-rev[i + 1] * v) % mod
    return B, rev


# calculate ∑_{k=1}^n k^p on modulo mod (求1加到n的p次方之和模mod)
# 时间复杂度O(p)
def calc(n: int, p: int, mod: int) -> int:
    B, rev = bernoulli(p, mod)

    res = 0
    tmp = pow(n, p + 1, mod)
    rev_n = pow(n, mod - 2, mod)
    for i in range(0, p + 1):
        res = (res + tmp * B[i]) % mod
        tmp = (-(p + 1 - i) * rev[i + 1] % mod) * (rev_n * tmp % mod) % mod
    return (res * rev[p + 1]) % mod


print(calc(10, 2, 10**9 + 7))
# => "385"
print(calc(100, 200, 10**9 + 7))
# => "710636539"
