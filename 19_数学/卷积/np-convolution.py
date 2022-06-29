"""时间复杂度 nlogn"""

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


if __name__ == '__main__':
    nums1 = [1, 2, 3]
    nums2 = [4, 5, 6]
    print(convolve(nums1, nums2))

