from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 二维 整数数组 coordinates 和一个整数 k ，其中 coordinates[i] = [xi, yi] 是第 i 个点在二维平面里的坐标。

# 我们定义两个点 (x1, y1) 和 (x2, y2) 的 距离 为 (x1 XOR x2) + (y1 XOR y2) ，XOR 指的是按位异或运算。


# 请你返回满足 i < j 且点 i 和点 j之间距离为 k 的点对数目。


class Solution:
    def countPairs(self, coordinates: List[List[int]], k: int) -> int:
        mp = defaultdict(lambda: defaultdict(int))  # x->y
        for x, y in coordinates:
            mp[x][y] += 1
        res = 0
        for x1, y1 in coordinates:
            for sum1 in range(0, k + 1):
                sum2 = k - sum1
                x2 = x1 ^ sum1
                y2 = y1 ^ sum2
                if (x1, y1) == (x2, y2):
                    res += mp[x2][y2] - 1
                else:
                    res += mp[x2][y2]
        return res // 2


# coordinates = [[1,2],[4,2],[1,3],[5,2]], k = 5

print(Solution().countPairs(coordinates=[[1, 2], [4, 2], [1, 3], [5, 2]], k=5))
# coordinates = [[1,3],[1,3],[1,3],[1,3],[1,3]], k = 0
print(Solution().countPairs(coordinates=[[1, 3], [1, 3], [1, 3], [1, 3], [1, 3]], k=0))
