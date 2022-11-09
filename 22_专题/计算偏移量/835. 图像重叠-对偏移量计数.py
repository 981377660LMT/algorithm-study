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
