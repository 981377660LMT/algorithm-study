from typing import List

INF = 0x7FFFFFFF
# 给定在 xy 平面上的一组点，确定由这些点组成的任何矩形的最小面积，其中矩形的边不一定平行于 x 轴和 y 轴。
# 如果没有任何矩形，就返回 0。

# 向量点乘判断是否是直角，向量叉乘得面积
# O(n^3)

# 1 <= points.length <= 50
class Solution:
    def minAreaFreeRect(self, points: List[List[int]]) -> float:
        n = len(points)
        pSet = set(map(tuple, points))
        res = INF

        # 三个点,i为交点
        for i in range(n):
            for j in range(i + 1, n):
                a, b = points[j][0] - points[i][0], points[j][1] - points[i][1]
                for k in range(j + 1, n):
                    c, d = points[k][0] - points[i][0], points[k][1] - points[i][1]
                    if a * c + b * d == 0 and (points[i][0] + a + c, points[i][1] + b + d) in pSet:
                        res = min(res, abs(a * d - b * c))

        return 0 if res == INF else res


print(Solution().minAreaFreeRect([[1, 2], [2, 1], [1, 0], [0, 1]]))
