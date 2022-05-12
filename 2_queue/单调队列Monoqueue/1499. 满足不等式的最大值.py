from typing import List
from collections import deque

# 数组中每个元素都表示二维平面上的点的坐标，并按照横坐标 x 的值从小到大排序。
# 2 <= points.length <= 10^5
# 请你找出 yi + yj + |xi - xj| 的 最大值，其中 |xi - xj| <= k 且 1 <= i < j <= points.length。
# yi + yj + |xi - xj| = (yi - xi) + (yj + xj)
# 因此，我们可以扫描所有的点，对于每一个点，看前面的点中，只要横坐标相差没有超过 k，其中最大的 y(i) - x(i) 是多少。
# 每次搜索到一个新的 point，我们可以通过单调队列，直接取出之前的点中，y - x 的最大值。不过，每次取出元素前，我们需要检验一下，这个元素的 x 值是否满足限制条件。


# 总结:
# lessons learned: when using a heap to maintain an ordered subsequence/subarray, consider using a monotonic stack.
# 当使用堆来维护有序的子序列/子数组时，考虑使用单调队列。


class Solution:
    def findMaxValueOfEquation(self, points: List[List[int]], k: int) -> int:
        queue = deque()
        res = -0x7FFFFFFF

        for x, y in points:
            # 1. 过期的数据
            while queue and queue[0][1] < x - k:
                queue.popleft()

            # 2.更新结果
            if queue:
                res = max(res, queue[0][0] + x + y)

            # 3.入队
            while queue and queue[-1][0] <= y - x:
                queue.pop()
            queue.append((y - x, x))

        return res


print(Solution().findMaxValueOfEquation(points=[[1, 3], [2, 0], [5, 10], [6, -10]], k=1))
# 输入：points = [[1,3],[2,0],[5,10],[6,-10]], k = 1
# 输出：4
# 解释：前两个点满足 |xi - xj| <= 1 ，代入方程计算，则得到值 3 + 0 + |1 - 2| = 4 。第三个和第四个点也满足条件，得到值 10 + -10 + |5 - 6| = 1 。
# 没有其他满足条件的点，所以返回 4 和 1 中最大的那个。

