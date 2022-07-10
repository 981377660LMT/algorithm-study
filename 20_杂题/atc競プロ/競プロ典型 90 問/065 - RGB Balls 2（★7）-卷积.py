# 三种颜色的球 一共选K个
# R+G <=X
# G+B<=Y
# B+R<=Z
# 一共有多少种选法

# !条件转换,考虑反面:
# B>=K-X
# R>=K-Y
# G>=K-Z

# !卷积可以求解 畳み込み (画成二维矩阵 X+Y为定值K时 卷积值为Y=-X+K的对角线的和)
import sys
import numpy as np
from functools import lru_cache

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = 998244353


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


def convolve(a: np.ndarray, b: np.ndarray) -> np.ndarray:
    """fft求卷积"""
    fftLen = 1
    while 2 * fftLen < len(a) + len(b) - 1:
        fftLen *= 2
    fftLen *= 2
    Fa = np.fft.rfft(a, fftLen)
    Fb = np.fft.rfft(b, fftLen)
    Fc = Fa * Fb
    res = np.fft.irfft(Fc, fftLen)
    res = np.rint(res).astype(np.int64)
    return res[: len(a) + len(b) - 1]


def convoleWithMod(a: object, b: object, mod: int) -> np.ndarray:
    """fft求卷积 取模"""
    npa = np.array(a, np.int64)
    npb = np.array(b, np.int64)

    a1, a2 = np.divmod(npa, 1 << 15)
    b1, b2 = np.divmod(npb, 1 << 15)

    x = convolve(a1, b1) % mod
    z = convolve(a2, b2) % mod
    y = (convolve(a1 + a2, b1 + b2) - (x + z)) % mod

    c = (x << 30) + (y << 15) + z
    return c % mod


R, G, B, K = map(int, input().split())
X, Y, Z = map(int, input().split())

nums1 = [0] * (B + 1)
for i in range(K - X, B + 1):
    nums1[i] = C(B, i)
nums2 = [0] * (R + 1)
for i in range(K - Y, R + 1):
    nums2[i] = C(R, i)
nums3 = [0] * (G + 1)
for i in range(K - Z, G + 1):
    nums3[i] = C(G, i)

print(convoleWithMod(nums1, convoleWithMod(nums2, nums3, MOD), MOD)[K])
