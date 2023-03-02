"""numpy fft"""

from typing import Any
import numpy as np


def convolve(nums1: Any, nums2: Any) -> "np.ndarray":
    """fft求卷积"""
    n, m = len(nums1), len(nums2)
    ph = 1 << (n + m - 2).bit_length()
    T = np.fft.rfft(nums1, ph) * np.fft.rfft(nums2, ph)
    res = np.fft.irfft(T, ph)[: n + m - 1]
    return np.rint(res).astype(np.int64)


if __name__ == "__main__":
    nums1 = [1, 2]
    nums2 = [1, 2, 1]
    print(convolve(nums1, nums2))

    n, m = map(int, input().split())
    a = list(map(int, input().split()))
    b = list(map(int, input().split()))
    print(*convolve(a, b))
