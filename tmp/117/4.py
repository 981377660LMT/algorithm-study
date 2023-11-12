from heapq import heapify, heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始大小为 m * n 的整数矩阵 values ，表示 m 个不同商店里 m * n 件不同的物品。每个商店有 n 件物品，第 i 个商店的第 j 件物品的价值为 values[i][j] 。除此以外，第 i 个商店的物品已经按照价值非递增排好序了，也就是说对于所有 0 <= j < n - 1 都有 values[i][j] >= values[i][j + 1] 。

# 每一天，你可以在一个商店里购买一件物品。具体来说，在第 d 天，你可以：

# 选择商店 i 。
# 购买数组中最右边的物品 j ，开销为 values[i][j] * d 。换句话说，选择该商店中还没购买过的物品中最大的下标 j ，并且花费 values[i][j] * d 去购买。
# 注意，所有物品都视为不同的物品。比方说如果你已经从商店 1 购买了物品 0 ，你还可以在别的商店里购买其他商店的物品 0 。


# 请你返回购买所有 m * n 件物品需要的 最大开销 。


class Solution:
    def maxSpending(self, values: List[List[int]]) -> int:
        ROW, COL = len(values), len(values[0])
        pq = []
        for r, row in enumerate(values):
            pq.append((row[-1], r, COL - 1))
        heapify(pq)
        res = 0
        for day in range(1, ROW * COL + 1):
            min_, r, c = heappop(pq)
            res += min_ * day
            if c > 0:
                heappush(pq, (values[r][c - 1], r, c - 1))
        return res


# values = [[8,5,2],[6,4,1],[9,7,3]]

print(Solution().maxSpending([[8, 5, 2], [6, 4, 1], [9, 7, 3]]))
