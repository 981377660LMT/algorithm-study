from typing import List
from KM算法模板 import KM


class Solution:
    def maxCompatibilitySum(self, students: List[List[int]], mentors: List[List[int]]) -> int:
        n, m = (len(students), len(students[0]))
        adjMatrix = [[0] * n for _ in range(n)]
        for i in range(n):
            for j in range(n):
                for k in range(m):
                    adjMatrix[i][j] += int(students[i][k] == mentors[j][k])

        return KM(adjMatrix).getResult()


print(
    Solution().maxCompatibilitySum(
        students=[[1, 1, 0], [1, 0, 1], [0, 0, 1]], mentors=[[1, 0, 0], [0, 0, 1], [1, 1, 0]]
    )
)
