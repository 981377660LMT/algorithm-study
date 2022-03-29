from typing import List
from functools import lru_cache

# 你需要制定一份 d 天的工作计划表。工作之间存在依赖，要想执行第 i 项工作，你必须完成全部 j 项工作（ 0 <= j < i）。
# 你每天 至少 需要完成一项任务。工作计划的总难度是这 d 天每一天的难度之和，而一天的工作难度是当天应该完成工作的最大难度。
# 返回整个工作计划的 最小难度

# 1 <= jobDifficulty.length <= 300
# 1 <= d <= 10

INF = 0x7FFFFFFF


class Solution:
    def minDifficulty(self, jobDifficulty: List[int], d: int) -> int:
        n = len(jobDifficulty)
        if n < d:
            return -1

        @lru_cache(None)
        def dfs(cur: int, remain: int) -> int:
            if remain == 1:
                return max(jobDifficulty[cur:])

            res = INF
            dayMax = -INF
            for i in range(cur, n - 1):
                dayMax = max(dayMax, jobDifficulty[i])
                res = min(res, dayMax + dfs(i + 1, remain - 1))

            return res

        return dfs(0, d)


print(Solution().minDifficulty(jobDifficulty=[6, 5, 4, 3, 2, 1], d=2))

# 输出：7
# 解释：第一天，您可以完成前 5 项工作，总难度 = 6.
# 第二天，您可以完成最后一项工作，总难度 = 1.
# 计划表的难度 = 6 + 1 = 7

