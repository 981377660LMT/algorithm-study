# 这里有 n 个一样的骰子，每个骰子上都有 k 个面，分别标号为 1 到 k 。
# 给定三个整数 n ,  k 和 target ，
# 返回可能的方式(从总共 k^n 种方式中)滚动骰子的数量，
# 使正面朝上的数字之和等于 target 。
# 答案可能很大，你需要对 1e9 + 7 取模 。

MOD = int(1e9 + 7)


class Solution:
    def numRollsToTarget(self, n: int, k: int, target: int) -> int:
        poly = [1] * k
        for _ in range(n - 1):
            poly = convoleWithMod(poly, [1] * k, MOD)

        try:
            return int(poly[target - n] % MOD)
        except IndexError:
            return 0


from typing import Any
import numpy as np


def convolve(nums1: Any, nums2: Any) -> np.ndarray:
    """fft求卷积"""
    fftLen = 1
    while 2 * fftLen < len(nums1) + len(nums2) - 1:
        fftLen *= 2
    fftLen *= 2
    Fa = np.fft.rfft(nums1, fftLen)
    Fb = np.fft.rfft(nums2, fftLen)
    Fc = Fa * Fb
    res = np.fft.irfft(Fc, fftLen)
    res = np.rint(res).astype(np.int64)
    return res[: len(nums1) + len(nums2) - 1]


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


print(Solution().numRollsToTarget(200, 6, 7))
