from typing import List


# 1 <= m, n <= 100
# 将球导向右侧的挡板跨过左上角和右下角，在网格中用 1 表示。
# 将球导向左侧的挡板跨过右上角和左下角，在网格中用 -1 表示。
# 返回一个大小为 n 的数组 answer ，
# 其中 answer[i] 是球放在顶部的第 i 列后从底部掉出来的那一列对应的下标，如果球卡在盒子里，则返回 -1 。
class Solution:
    def findBall(self, grid: List[List[int]]) -> List[int]:
        m, n = len(grid), len(grid[0])

        def dfs(i: int, j: int) -> int:
            # 出界
            if j < 0 or j >= n:
                return -1
            if i >= m:
                return j

            # V型
            if j + 1 < n and grid[i][j] == 1 and grid[i][j + 1] == -1:
                return -1
            if j - 1 > -1 and grid[i][j] == -1 and grid[i][j - 1] == 1:
                return -1

            if grid[i][j] == 1:
                return dfs(i + 1, j + 1)
            else:
                return dfs(i + 1, j - 1)

        return [dfs(0, j) for j in range(n)]


print(
    Solution().findBall(
        grid=[
            [1, 1, 1, -1, -1],
            [1, 1, 1, -1, -1],
            [-1, -1, -1, 1, 1],
            [1, 1, 1, 1, -1],
            [-1, -1, -1, -1, -1],
        ]
    )
)
# 输出：[1,-1,-1,-1,-1]
