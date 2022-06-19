from typing import List
from functools import lru_cache


class Solution:
    # 1 <= questions.length <= 105
    # 从数据量看肯定状态只有一个
    def mostPoints(self, questions: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(cur: int) -> int:
            if cur >= n:
                return 0
            score, jump = questions[cur]
            return max(score + dfs(cur + jump + 1), dfs(cur + 1))

        n = len(questions)
        return dfs(0)


# 5 7 157
print(Solution().mostPoints(questions=[[3, 2], [4, 3], [4, 4], [2, 5]]))
print(Solution().mostPoints(questions=[[1, 1], [2, 2], [3, 3], [4, 4], [5, 5]]))
print(
    Solution().mostPoints(
        questions=[[21, 5], [92, 3], [74, 2], [39, 4], [58, 2], [5, 5], [49, 4], [65, 3]]
    )
)
