from typing import List
from functools import lru_cache

# 每个学生都会被分配给 一名 导师，而每位导师也会分配到 一名 学生。配对的学生与导师之间的兼容性评分等于学生和导师答案相同的次数。
# 请你找出最优的学生与导师的配对方案，以 最大程度上 提高 兼容性评分和 。
# 1 <= m, n <= 8


class Solution:
    def maxCompatibilitySum(self, students: List[List[int]], mentors: List[List[int]]) -> int:
        n, m = len(students), len(students[0])
        target = (1 << n) - 1
        weight = [[0] * n for _ in range(n)]
        for i in range(n):
            for j in range(n):
                for k in range(m):
                    weight[i][j] += int(students[i][k] == mentors[j][k])

        @lru_cache(None)
        def dfs(cur: int, state: int) -> int:
            if state == target:
                return 0

            return max(
                dfs(cur + 1, state | (1 << next)) + weight[cur][next]
                for next in range(n)
                if not state & (1 << next)
            )

        return dfs(0, 0)


print(
    Solution().maxCompatibilitySum(
        students=[[1, 1, 0], [1, 0, 1], [0, 0, 1]], mentors=[[1, 0, 0], [0, 0, 1], [1, 1, 0]]
    )
)
# 输出：8
# 解释：按下述方式分配学生和导师：
# - 学生 0 分配给导师 2 ，兼容性评分为 3 。
# - 学生 1 分配给导师 0 ，兼容性评分为 2 。
# - 学生 2 分配给导师 1 ，兼容性评分为 3 。
# 最大兼容性评分和为 3 + 2 + 3 = 8 。
