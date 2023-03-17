# https://yukicoder.me/problems/no/1618
# 给定长为n的数组A和B,求出长为2n的数组C
# !其中 C[k] = ∑i*A[i]+j*B[j] (i+j=k)

# n<=2e5 0<=nums[i]<2e5


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


n = int(input())
A = list(map(int, input().split()))
B = list(map(int, input().split()))


A = [0] + A
B = [0] + B
f = list(range(n + 1))
A = convolution_fft_large(A, f)
B = convolution_fft_large(B, f)
C = [a + b for a, b in zip(A, B)]
C = C[1:]
print(*C)
