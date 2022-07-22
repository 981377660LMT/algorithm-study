# https://leetcode.cn/problems/finding-the-number-of-visible-mountains/
# !如果一座山峰不在另一座山峰之内(包括其他山峰的边界) ，那么这座山峰就被认为是可见的

from collections import defaultdict
from typing import List


class Solution:
    def visibleMountains(self, peaks: List[List[int]]) -> int:
        """逆时针旋转点 + 二维偏序看这个点是否被其他山峰包含"""
        points = [(x - y, x + y) for x, y in peaks]
        adjMap = defaultdict(list)  # !按照x坐标分组
        for x, y in points:
            adjMap[x].append((x, y))

        keys = sorted(adjMap)
        res, maxY = 0, -int(1e20)
        for key in keys:
            group = adjMap[key]
            cur = 0
            for _, py in group:
                if py > maxY:
                    maxY = py
                    cur = 1
                elif py == maxY:
                    cur = 0
            res += cur
        return res


print(Solution().visibleMountains(peaks=[[2, 2], [6, 3], [5, 4]]))
print(Solution().visibleMountains(peaks=[[1, 3], [1, 3]]))
print(Solution().visibleMountains(peaks=[[38, 26], [3, 32], [2, 1]]))
