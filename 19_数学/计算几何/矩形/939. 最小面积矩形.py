# 939. 最小面积矩形
# https://leetcode.cn/problems/minimum-area-rectangle/
# 给定在 xy 平面上的一组点，确定由这些点组成的矩形的最小面积，其中矩形的边平行于 x 轴和 y 轴。
# 如果没有任何矩形，就返回 0。
# !总结：遍历，找对角线


from collections import defaultdict
from typing import List


INF = int(1e18)


class Solution:
    def minAreaRect(self, points: List[List[int]]) -> int:
        xToYs = defaultdict(list)
        for x, y in points:
            xToYs[x].append(y)
        for ys in xToYs.values():
            ys.sort()

        res = INF
        lastY1Y2 = dict()
        for x in sorted(xToYs):
            ys = xToYs[x]
            m = len(ys)
            for i in range(m):
                y1 = ys[i]
                for j in range(i + 1, m):
                    y2 = ys[j]
                    key = (y1, y2)
                    if key in lastY1Y2:
                        w, h = x - lastY1Y2[key], y2 - y1
                        area = w * h
                        if area < res:
                            res = area
                    lastY1Y2[key] = x

        return res if res < INF else 0


print(Solution().minAreaRect([[1, 1], [1, 3], [3, 1], [3, 3], [2, 2]]))
