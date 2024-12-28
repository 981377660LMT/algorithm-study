# !按值域统计，对数量比较多的 转换为卷积

import numpy as np
from collections import defaultdict
from typing import List


def fft_convolution(a, b):
    nA = len(a)
    nB = len(b)
    size = nA + nB - 1
    size_fft = 1
    while size_fft < size:
        size_fft <<= 1
    A = np.fft.rfft(a, size_fft)
    B = np.fft.rfft(b, size_fft)
    C = A * B
    c = np.fft.irfft(C, size_fft)
    return np.round(c[:size]).astype(np.int64)


def circular_correlation(a, b):
    n = len(a)
    b_rev = b[::-1]
    lin_corr = fft_convolution(a, b_rev)
    res = [0] * n
    for i, val in enumerate(lin_corr):
        res[i % n] += val
    return res


class Solution:
    def maximumMatchingIndices(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        if n == 0:
            return 0
        if nums1 == nums2:
            return n

        distinct_vals = list(set(nums1 + nums2))
        val2id = {v: idx for idx, v in enumerate(distinct_vals)}
        M = len(distinct_vals)

        freq1 = defaultdict(int)
        freq2 = defaultdict(int)
        for v in nums1:
            freq1[v] += 1
        for v in nums2:
            freq2[v] += 1

        circular_corr_global = [0] * n

        import math

        threshold = n * math.log2(n) if n > 1 else 1

        for val in distinct_vals:
            count1 = freq1[val]
            count2 = freq2[val]
            if count1 == 0 or count2 == 0:
                continue

            A_j = np.zeros(n, dtype=np.int64)
            B_j = np.zeros(n, dtype=np.int64)

            for i, v in enumerate(nums1):
                if v == val:
                    A_j[i] = 1
            for i, v in enumerate(nums2):
                if v == val:
                    B_j[i] = 1

            if count1 * count2 > threshold:
                corr_j = circular_correlation(A_j, B_j)
                for i in range(n):
                    circular_corr_global[i] += corr_j[i]
            else:
                posA = [i for i in range(n) if A_j[i] == 1]
                posB = [i for i in range(n) if B_j[i] == 1]
                for a in posA:
                    for b in posB:
                        k = (b - a) % n
                        circular_corr_global[k] += 1

        return max(circular_corr_global)
