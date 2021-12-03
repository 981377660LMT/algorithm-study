from typing import List

# 在给定的二维二进制数组 A 中，存在两座岛。（岛是由四面相连的 1 形成的一个最大组。）
# 现在，我们可以将 0 变为 1，以使两座岛连接起来，变成一座岛。
# 返回必须翻转的 0 的最小数目。（可以保证答案至少是 1 。）

# 思路:
# 1.找起点
# 2.dfs将岛全部加入queue 原地标记-1
# 3.bfs最短路径 找到1就返回
class Solution:
    def shortestBridge(self, A: List[List[int]]) -> int:
        def dfs(i, j):
            A[i][j] = -1
            queue.append((i, j))
            for x, y in ((i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1)):
                if 0 <= x < n and 0 <= y < n and A[x][y] == 1:
                    dfs(x, y)

        def first():
            for i in range(n):
                for j in range(n):
                    if A[i][j]:
                        return i, j

        n, step, queue = len(A), 0, []
        dfs(*first())

        while queue:
            nextQueue = []
            for i, j in queue:
                for x, y in ((i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1)):
                    if 0 <= x < n and 0 <= y < n:
                        if A[x][y] == 1:
                            return step
                        elif not A[x][y]:
                            A[x][y] = -1
                            nextQueue.append((x, y))
            step += 1
            queue = nextQueue


print(
    Solution().shortestBridge(
        [[1, 1, 1, 1, 1], [1, 0, 0, 0, 1], [1, 0, 1, 0, 1], [1, 0, 0, 0, 1], [1, 1, 1, 1, 1]]
    )
)
