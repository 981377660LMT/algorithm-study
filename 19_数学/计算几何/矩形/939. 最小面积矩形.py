from typing import List

# 给定在 xy 平面上的一组点，确定由这些点组成的矩形的最小面积，其中矩形的边平行于 x 轴和 y 轴。
# 如果没有任何矩形，就返回 0。

INF = 0x7FFFFFFF
# 总结：遍历，找对角线
class Solution:
    def minAreaRect(self, points: List[List[int]]) -> int:
        n = len(points)
        pSet = set(map(tuple, points))
        res = INF

        for i in range(n):
            x1, y1 = points[i]
            for j in range(i + 1, n):
                x2, y2 = points[j]
                if x1 != x2 and y1 != y2:
                    p3 = (x2, y1)
                    p4 = (x1, y2)
                    if p3 in pSet and p4 in pSet:
                        cur = abs(x1 - x2) * abs(y1 - y2)  # 当前的面积
                        res = min(res, cur)

        return 0 if res == INF else res


print(Solution().minAreaRect([[1, 1], [1, 3], [3, 1], [3, 3], [2, 2]]))
