# https://leetcode.cn/problems/kjpLFZ/

# 提取装置初始位于矩阵的左上角 [0,0]，可以通过每次操作移动到上、下、左、右相邻的 1 格位置中。
# 提取装置每次移动或每次提取均记为一次操作。
# 远征队需要按照顺序，从矩阵中逐一取出字母以组成 mantra，
# 才能够成功的启动升降机。请返回他们 最少 需要消耗的操作次数。如果无法完成提取，返回 -1。
# 注意：
# 提取装置可对同一位置的字母重复提取，每次提取一个
# 提取字母时，需按词语顺序依次提取

# !状态为(index,row,col) 多维bfs求最短路


from collections import deque
from typing import List


INF = int(1e18)
DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def extractMantra(self, matrix: List[str], mantra: str) -> int:
        ROW, COL, n = len(matrix), len(matrix[0]), len(mantra)
        dist = [[INF] * (n + 1) for _ in range(ROW * COL)]
        dist[0][0] = 0
        queue = deque([(0, 0, 0, 0)])  # (step, index, x, y)
        while queue:
            step, index, x, y = queue.popleft()
            if index == n:
                continue
            if step > dist[x * COL + y][index]:
                continue

            if matrix[x][y] == mantra[index]:  # !提取
                cand = step + 1
                if cand < dist[x * COL + y][index + 1]:
                    dist[x * COL + y][index + 1] = cand
                    queue.append((step + 1, index + 1, x, y))  # type: ignore

            for dx, dy in DIR4:  # !移动
                nx, ny = x + dx, y + dy
                if 0 <= nx < ROW and 0 <= ny < COL:
                    cand = step + 1
                    if cand < dist[nx * COL + ny][index]:
                        dist[nx * COL + ny][index] = cand
                        queue.append((cand, index, nx, ny))  # type: ignore

        res = min(d[n] for d in dist)
        return res if res < INF else -1
