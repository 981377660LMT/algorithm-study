from typing import List
from KM算法模板 import KM


class Solution:
    def maxCompatibilitySum(self, students: List[List[int]], mentors: List[List[int]]) -> int:
        n, m = (len(students), len(students[0]))
        costMatrix = [[0] * n for _ in range(n)]
        for i in range(n):
            for j in range(n):
                for k in range(m):
                    costMatrix[i][j] += int(students[i][k] == mentors[j][k])

        return KM(costMatrix)[0]


print(
    Solution().maxCompatibilitySum(
        students=[[1, 1, 0], [1, 0, 1], [0, 0, 1]], mentors=[[1, 0, 0], [0, 0, 1], [1, 1, 0]]
    )
)
