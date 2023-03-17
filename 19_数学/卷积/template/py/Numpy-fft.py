"""numpy fft"""

from typing import Any, List
import numpy as np


def convolution(nums1: Any, nums2: Any) -> "np.ndarray":
    """fft求卷积(可能精度不够)"""
    n, m = len(nums1), len(nums2)
    ph = 1 << (n + m - 2).bit_length()
    T = np.fft.rfft(nums1, ph) * np.fft.rfft(nums2, ph)
    res = np.fft.irfft(T, ph)[: n + m - 1]
    return np.rint(res).astype(np.int64)


def convolution_fft_large(a: Any, b: Any) -> List[int]:
    """精度不够用这个"""
    a, b = np.array(a, dtype=np.int64), np.array(b, dtype=np.int64)
    d = 1 << 10
    a1, a2 = np.divmod(a, d * d)
    a2, a3 = np.divmod(a2, d)
    b1, b2 = np.divmod(b, d * d)
    b2, b3 = np.divmod(b2, d)
    aa = convolution(a1, b1)
    bb = convolution(a2, b2)
    cc = convolution(a3, b3)
    dd = convolution(a1 + a2, b1 + b2) - (aa + bb)  # type: ignore
    ee = convolution(a2 + a3, b2 + b3) - (bb + cc)  # type: ignore
    ff = convolution(a1 + a3, b1 + b3) - (aa + cc)  # type: ignore
    h = ((aa * d * d)) * d * d + ((dd * d * d)) * d + (bb + ff) * d * d + ee * d + cc
    return h.tolist()


if __name__ == "__main__":
    nums1 = [1, 2]
    nums2 = [1, 2, 1]
    print(convolution(nums1, nums2))

    n, m = map(int, input().split())
    a = list(map(int, input().split()))
    b = list(map(int, input().split()))
    print(*convolution(a, b))
