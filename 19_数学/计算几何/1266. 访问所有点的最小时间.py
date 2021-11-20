from typing import List

# 横竖斜
# 必须按照数组中出现的顺序来访问这些点。
# 水桶定理 最长边影响最小时间


class Solution:
    def minTimeToVisitAllPoints(self, points: List[List[int]]) -> int:
        return sum(
            max(abs(points[i][0] - points[i - 1][0]), abs(points[i][1] - points[i - 1][1]))
            for i in range(1, len(points))
        )

