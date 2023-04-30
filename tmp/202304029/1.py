from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个数组 start ，其中 start = [startX, startY] 表示你的初始位置位于二维空间上的 (startX, startY) 。另给你一个数组 target ，其中 target = [targetX, targetY] 表示你的目标位置 (targetX, targetY) 。

# 从位置 (x1, y1) 到空间中任一其他位置 (x2, y2) 的代价是 |x2 - x1| + |y2 - y1| 。

# 给你一个二维数组 specialRoads ，表示空间中存在的一些特殊路径。其中 specialRoads[i] = [x1i, y1i, x2i, y2i, costi] 表示第 i 条特殊路径可以从 (x1i, y1i) 到 (x2i, y2i) ，但成本等于 costi 。你可以使用每条特殊路径任意次数。


# 返回从 (startX, startY) 到 (targetX, targetY) 所需的最小代价。


class Solution:
    def minimumCost(
        self, start: List[int], target: List[int], specialRoads: List[List[int]]
    ) -> int:
        dist = defaultdict(lambda: INF)
        sr, sc = start
        tr, tc = target
        dist[(sr, sc)] = 0
        pq = [(0, sr, sc)]

        mp = defaultdict(list)
        for x1, y1, x2, y2, cost in specialRoads:
            mp[(x1, y1)].append((x2, y2, cost))

        while pq:
            curDist, r, c = heappop(pq)
            if (r, c) == (tr, tc):
                return curDist
            if dist[(r, c)] < curDist:
                continue

            cand1 = curDist + abs(tr - r) + abs(tc - c)
            if cand1 < dist[(tr, tc)]:
                dist[(tr, tc)] = cand1
                heappush(pq, (cand1, tr, tc))
            for nextX, nextY in mp:
                cand2 = curDist + abs(nextX - r) + abs(nextY - c)
                if cand2 < dist[(nextX, nextY)]:
                    dist[(nextX, nextY)] = cand2
                    heappush(pq, (cand2, nextX, nextY))

            for nr, nc, w in mp[(r, c)]:
                cand = curDist + w
                if cand < dist[(nr, nc)]:
                    dist[(nr, nc)] = cand
                    heappush(pq, (cand, nr, nc))

        return -1


# start = [1,1], target = [4,5], specialRoads = [[1,2,3,3,2],[3,4,4,5,1]]
print(Solution().minimumCost([1, 1], [4, 5], [[1, 2, 3, 3, 2], [3, 4, 4, 5, 1]]))
