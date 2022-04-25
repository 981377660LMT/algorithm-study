from typing import List
from scipy.spatial import ConvexHull

Point = List[int]


def isOnOneLine(points: List[Point]) -> bool:
    def calCrossProduct(A: Point, B: Point, C: Point) -> int:
        """"计算三点叉乘"""

        AB = [B[0] - A[0], B[1] - A[1]]
        AC = [C[0] - A[0], C[1] - A[1]]
        return AB[0] * AC[1] - AB[1] * AC[0]

    """"判断k点共线"""

    if len(points) <= 2:
        return True

    p1, p2, *restP = points
    return all(calCrossProduct(p1, p2, p3) == 0 for p3 in restP)


class Solution:
    def outerTrees(self, trees: List[List[int]]) -> List[List[int]]:
        # allPoints = set(trees)
        res1 = [trees[i] for i in ConvexHull(trees).vertices]  # 凸包的边
        print(res1)


print(Solution().outerTrees([[1, 1], [2, 2], [2, 0], [2, 4], [3, 3], [4, 2]]))

# 输出: [[1,1],[2,0],[4,2],[3,3],[2,4]]
