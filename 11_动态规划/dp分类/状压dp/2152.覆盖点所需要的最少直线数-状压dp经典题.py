from typing import List


# 1 <= points.length <= 10

Point = List[int]


class Solution:
    def minimumLines(self, points: List[List[int]]) -> int:
        n = len(points)
        dp = [n] * (1 << n)

        for state in range(1 << n):
            if self._isOnOneLine([points[i] for i in range(n) if ((state >> i) & 1)]):
                dp[state] = 1
                continue

            selected = [i for i in range(n) if (state >> i) & 1]
            for group1 in range(1 << len(selected)):
                group2 = state ^ group1
                dp[state] = min(dp[state], dp[group1] + dp[group2])

        return dp[-1]

    @classmethod
    def _isOnOneLine(cls, points: List[Point]) -> bool:
        """ "判断k点共线"""

        if len(points) <= 2:
            return True

        p1, p2, *restP = points
        return all(cls._calCrossProduct(p1, p2, p3) == 0 for p3 in restP)

    @staticmethod
    def _calCrossProduct(A: Point, B: Point, C: Point) -> int:
        """ "计算三点叉乘"""

        AB = [B[0] - A[0], B[1] - A[1]]
        AC = [C[0] - A[0], C[1] - A[1]]
        return AB[0] * AC[1] - AB[1] * AC[0]


# 2 4 4
# print(Solution().bestLine([[0, 0], [1, 1], [1, 0], [2, 0]]))
print(Solution().minimumLines(points=[[0, 1], [2, 3], [4, 5], [4, 3]]))
print(
    Solution().minimumLines(
        points=[[-2, 2], [4, -1], [-5, -3], [1, 0], [-1, -3], [-2, 0], [-4, -4]]
    )
)
print(
    Solution().minimumLines(
        points=[
            [4, -1],
            [2, -4],
            [2, -1],
            [1, -1],
            [3, 3],
            [2, 2],
            [-4, 4],
            [-5, 1],
            [0, 4],
            [-1, -5],
        ]
    )
)
