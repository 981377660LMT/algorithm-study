# https://leetcode.cn/problems/cherry-pickup/solutions/453318/shi-yong-ayou-hua-yun-xing-su-du-by-ling-jian-2012/
#
# !1.按照剩余步数 * 2来估计最大得分，注意没必要真的将这个估值函数算出来，只要考虑状态转化之后估值函数下降的值就可以了
# !2.值域有限时，使用多个stack代替优先队列
# !3.为了解决没有通路的问题，提前进行一次O(n^2)的测试来判断起点和终点之间是否有通路。
from typing import List


class Solution:
    def cherryPickup(self, grid: List[List[int]]) -> int:
        N = len(grid)
        if N == 1:
            return grid[0][0]
        # quick test for connectivity
        c = [False] * N
        c[0] = True
        for i in range(1, N):
            if grid[0][i] < 0:
                break
            c[i] = True
        for i in range(1, N):
            c[0] = c[0] and grid[i][0] >= 0
            for j in range(1, N):
                c[j] = (c[j - 1] or c[j]) and grid[i][j] >= 0
        if not c[N - 1]:
            return 0
        # A*
        visited = set([(0, 0, 0, 0)])
        q = [[(grid[0][0], 0, 0, 0, 0)], [], []]
        while any(q):
            while not q[0]:
                q[0], q[1], q[2] = q[1], q[2], q[0]
            c, x1, y1, x2, y2 = q[0].pop()
            for dx1, dy1 in ((1, 0), (0, 1)):
                if x1 + dx1 < N and y1 + dy1 < N and grid[x1 + dx1][y1 + dy1] >= 0:
                    for dx2, dy2 in ((0, 1), (1, 0)):
                        if (
                            x2 + dx2 < N
                            and y2 + dy2 < N
                            and grid[x2 + dx2][y2 + dy2] >= 0
                            and (x1 + dx1, y1 + dy1, x2 + dx2, y2 + dy2) not in visited
                        ):
                            if x1 + dx1 == x2 + dx2 and y1 + dy1 == y2 + dy2:
                                dc = grid[x1 + dx1][y1 + dy1]
                                if x1 + dx1 == N - 1 and y1 + dy1 == N - 1:
                                    return c + dc
                            else:
                                dc = grid[x1 + dx1][y1 + dy1] + grid[x2 + dx2][y2 + dy2]
                            visited.add((x1 + dx1, y1 + dy1, x2 + dx2, y2 + dy2))
                            q[2 - dc].append((c + dc, x1 + dx1, y1 + dy1, x2 + dx2, y2 + dy2))
