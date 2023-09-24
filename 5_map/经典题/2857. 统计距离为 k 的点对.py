# https://leetcode.cn/problems/count-pairs-of-points-with-distance-k/
# 2857. 统计距离为 k 的点对
# 给你一个 二维 整数数组 coordinates 和一个整数 k ，其中 coordinates[i] = [xi, yi] 是第 i 个点在二维平面里的坐标。
# 我们定义两个点 (x1, y1) 和 (x2, y2) 的 距离 为 (x1 XOR x2) + (y1 XOR y2) ，XOR 指的是按位异或运算。
# 请你返回满足 i < j 且点 i 和点 j之间距离为 k 的点对数目。


# 2 <= coordinates.length <= 5e4
# 0 <= xi, yi <= 1e6
# 0 <= k <= 100

from collections import defaultdict
from typing import List


class Solution:
    def countPairs(self, coordinates: List[List[int]], k: int) -> int:
        res = 0
        counter = defaultdict(int)
        for x, y in coordinates:
            for i in range(k + 1):
                res += counter[(x ^ i, y ^ (k - i))]
            counter[(x, y)] += 1
        return res
