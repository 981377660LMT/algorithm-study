from typing import List
from queue import PriorityQueue

# 每个单元格不是 0（空）就是 1（障碍物）
# 您 最多 可以消除 k 个障碍物，请找出从左上角 (0, 0) 到
# 右下角 (m-1, n-1) 的最短路径，并返回通过该路径所需的步数。如果找不到这样的路径，则返回 -1。

# 1 <= m, n <= 40
'''
A* 搜索最佳状态
'''

from typing import List
from queue import PriorityQueue


class Solution:
    def shortestPath(self, grid: List[List[int]], k: int) -> int:
        que = PriorityQueue()
        m, n = len(grid), len(grid[0])

        que.put((m - 1 + n - 1, 0, -k, 0, 0))  # (预估距离总和，当前开销, -当前剩余消除障碍次数, 行位置，列位置)
        best_stat = {(0, 0, k): 0}

        while not que.empty():
            predict, cost, cur_k, i, j = que.get()
            cur_k = 0 - cur_k

            if i == m - 1 and j == n - 1:
                return cost

            new_cost = cost + 1
            for ii, jj in [(i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1)]:
                if ii >= 0 and ii < m and jj >= 0 and jj < n:
                    if grid[ii][jj] == 1:
                        if cur_k > 0:
                            new_k = cur_k - 1
                        else:
                            continue
                    else:
                        new_k = cur_k

                    if (ii, jj, new_k) not in best_stat or new_cost < best_stat[(ii, jj, new_k)]:
                        best_stat[(ii, jj, new_k)] = new_cost
                        que.put((new_cost + (m - 1 - ii) + (n - 1 - jj), new_cost, -new_k, ii, jj))

        return -1


print(Solution().shortestPath([[0, 0, 0], [1, 1, 0], [0, 0, 0], [0, 1, 1], [0, 0, 0]], 1))
