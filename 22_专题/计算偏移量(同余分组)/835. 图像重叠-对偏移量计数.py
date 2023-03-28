from typing import List
from collections import Counter

# O(n^4)
# n<=30

# 直接对偏移量计数
# 835. 图像重叠-对偏移量计数


class Solution:
    def largestOverlap(self, img1: List[List[int]], img2: List[List[int]]) -> int:
        diffCounter = Counter()  # 偏移量元组
        for r1, row1 in enumerate(img1):
            for c1, v1 in enumerate(row1):
                if v1 == 0:
                    continue
                for r2, row2 in enumerate(img2):
                    for c2, v2 in enumerate(row2):
                        if v2 == 0:
                            continue
                        diffCounter[(r1 - r2, c1 - c2)] += 1
        return max(diffCounter.values(), default=0)


print(
    Solution().largestOverlap(
        img1=[[1, 1, 0], [0, 1, 0], [0, 1, 0]], img2=[[0, 0, 0], [0, 1, 1], [0, 0, 1]]
    )
)
# 输出：3
# 解释：将 img1 向右移动 1 个单位，再向下移动 1 个单位。

# 二维fft
# https://leetcode.cn/problems/image-overlap/solution/ni-ke-neng-wu-fa-xiang-xiang-de-on2lognd-gc5j/
import numpy as np


class Solution:
    def largestOverlap(self, img1: List[List[int]], img2: List[List[int]]) -> int:
        N = len(img1)
        N2 = 1 << (N.bit_length() + 1)
        img1_fft = np.fft.fft2(np.array(img1), (N2, N2))
        img2_fft = np.fft.fft2(np.array(img2)[::-1, ::-1], (N2, N2))
        img1_fft *= img2_fft
        conv = np.fft.ifft2(img1_fft)
        return int(np.round(np.max(conv)))
