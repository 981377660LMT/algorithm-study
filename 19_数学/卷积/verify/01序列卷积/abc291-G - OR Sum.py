# 给定两个长度为 n 的数组 ，你可以执行任意次以下操作
# 将A元素所有元素向左轮转一位
# 求能得到的最大的 sum(A[i] | B[i]) 的值
# !2<=n<=5e5,0<=ai,bi<=31


# !按位卷积
# 1.翻转数组,把zip变成卷积乘法的形式
# 2.按位or可以表示为a[i]+b[i]-a[i]&b[i],01序列即为a[i]+b[i]-a[i]*b[i]
# 3.卷积
from typing import List
from typing import Any
import numpy as np


def convolve(nums1: Any, nums2: Any) -> "np.ndarray":
    """fft求卷积"""
    n, m = len(nums1), len(nums2)
    ph = 1 << (n + m - 2).bit_length()
    T = np.fft.rfft(nums1, ph) * np.fft.rfft(nums2, ph)
    res = np.fft.irfft(T, ph)[: n + m - 1]
    return np.rint(res).astype(np.int64)


def xorSum(nums1: List[int], nums2: List[int]) -> int:
    n = len(nums1)
    nums2 = nums2[::-1]
    res = [0] * n  # 每个轮转对应的和
    for bit in range(5):
        A = [(nums1[i] >> bit) & 1 for i in range(n)]
        B = [(nums2[i] >> bit) & 1 for i in range(n)]
        sum_ = A.count(1) + B.count(1)
        conv = convolve(A, B)
        for i in range(n, n + n - 1):
            conv[i % n] += conv[i]
        for i in range(n):
            res[i] += (sum_ - conv[i]) * (1 << bit)

    return max(res)


if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    print(xorSum(nums1, nums2))
