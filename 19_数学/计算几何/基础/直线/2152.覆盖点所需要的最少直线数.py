from typing import List, Set
from collections import defaultdict

# 这道题是 149. 直线上最多的点数、面试题 16.14. 最佳直线 的加强版

# 1.bestLine 函数 即 面试题 16.14. 最佳直线 ，该函数返回所有的通过点的数目最多的直线，每条直线用点的index组成的集合表示
# 2.递归，在剩下的点中继续寻找通过点的数目最多的直线，直到最后剩余点数<=2

# 1 <= points.length <= 10


class Solution:
    def minimumLines(self, points: List[List[int]]) -> int:
        n = len(points)
        if n == 0:
            return 0
        if n in (1, 2):
            return 1

        bestLines = self.bestLine(points)
        res = 0x3F3F3F3F
        for selectedIds in bestLines:
            remainIds = set(range(n)) - selectedIds
            res = min(res, 1 + self.minimumLines([points[i] for i in remainIds]))
        return res

    @staticmethod
    def bestLine(points: List[List[int]]) -> List[Set[int]]:
        def gcd(a, b):
            return a if b == 0 else gcd(b, a % b)

        n = len(points)
        res = [set()]
        maxCount = 0

        for i in range(n):
            x1, y1 = points[i]
            groups = defaultdict(set)

            for j in range(i + 1, n):
                x2, y2 = points[j]
                A, B = y2 - y1, x2 - x1
                if B == 0:
                    key = (0, 0)
                else:
                    gcd_ = gcd(A, B)
                    key = (A / gcd_, B / gcd_)

                groups[key].add(j)
                count = len(groups[key])
                if count > maxCount:
                    maxCount = count
                    res = [groups[key] | {i}]
                elif count == maxCount:
                    res.append(groups[key] | {i})

        return res


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
