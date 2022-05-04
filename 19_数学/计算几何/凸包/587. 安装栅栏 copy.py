from typing import List
from scipy.spatial import ConvexHull

Point = List[int]


def calCrossProduct(A: Point, B: Point, C: Point) -> int:
    """"计算AB与AC的叉乘"""

    AB = [B[0] - A[0], B[1] - A[1]]
    AC = [C[0] - A[0], C[1] - A[1]]
    return AB[0] * AC[1] - AB[1] * AC[0]


class Solution:
    def outerTrees(self, trees: List[List[int]]) -> List[List[int]]:
        """Andrew 算法 求出所有的凸包
        
        逆时针转一圈 删除顺时针的点
        """
        if len(trees) <= 3:
            return trees

        points = sorted(trees)
        stack = []

        # 寻找凸壳的下半部分
        for i in range(len(points)):
            while len(stack) >= 2 and calCrossProduct(stack[-2], stack[-1], points[i]) < 0:
                stack.pop()
            stack.append(tuple(points[i]))

        # 寻找凸壳的上半部分

        for i in range(len(points) - 1, -1, -1):
            while len(stack) >= 2 and calCrossProduct(stack[-2], stack[-1], points[i]) < 0:
                stack.pop()
            stack.append(tuple(points[i]))
        return list(set(stack))


print(Solution().outerTrees([[1, 1], [2, 2], [2, 0], [2, 4], [3, 3], [4, 2]]))

# 输出: [[1,1],[2,0],[4,2],[3,3],[2,4]]
