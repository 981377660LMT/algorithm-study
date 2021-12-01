from typing import List, Optional, Tuple

# 你需要检查 grid 中是否存在 相同值 形成的环。
# 一个环是一条开始和结束于同一个格子的长度 大于等于 4 的路径

# 总结:环的条件:dfs又走回visited
# dfs判断环需要记录parent防止直接重新走回

# 也可并查集


class Solution:
    def containsCycle(self, grid: List[List[str]]) -> bool:
        def dfs(cur: Tuple[int, int], parent: Optional[Tuple[int, int]]) -> bool:
            if cur in visited:
                return True
            visited.add(cur)
            nx, ny = cur
            nexts = [
                (cx, cy)
                for cx, cy in [[nx + 1, ny], [nx - 1, ny], [nx, ny + 1], [nx, ny - 1]]
                if 0 <= cx < m
                and 0 <= cy < n
                and grid[cx][cy] == grid[nx][ny]
                and (cx, cy) != parent
            ]

            for next in nexts:
                if dfs(next, cur):
                    return True
            return False

        m, n = len(grid), len(grid[0])
        visited = set()
        for i in range(m):
            for j in range(n):
                if (i, j) in visited:
                    continue
                if dfs((i, j), None):
                    return True
        return False


print(
    Solution().containsCycle(
        grid=[
            ["a", "a", "a", "a"],
            ["a", "b", "b", "a"],
            ["a", "b", "b", "a"],
            ["a", "a", "a", "a"],
        ]
    )
)
