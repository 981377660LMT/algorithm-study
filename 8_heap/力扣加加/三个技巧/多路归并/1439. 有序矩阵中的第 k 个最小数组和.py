from typing import List
from heapq import heappop, heappush

# 给你一个 m * n 的矩阵 mat，以及一个整数 k ，矩阵中的每一行都以非递减的顺序排列。
# 你可以从每一行中选出 1 个元素形成一个数组。返回所有可能数组中的第 k 个 最小 数组和。


class Solution:
    def kthSmallest(self, mat: List[List[int]], k: int) -> int:
        m, n = len(mat), len(mat[0])
        pointers = [0] * m
        pq = []
        curSum = sum(row[0] for row in mat)
        heappush(pq, (curSum, pointers))
        visited = set()
        visited.add(tuple(pointers))

        res = curSum

        for _ in range(k):
            c, p = heappop(pq)
            res = c
            for row, col in enumerate(p):
                if col < n - 1:
                    nextP = list(p)
                    nextP[row] = col + 1
                    nextP = tuple(nextP)
                    if nextP not in visited:
                        visited.add(nextP)
                        nextSum = res + mat[row][col + 1] - mat[row][col]
                        heappush(pq, (nextSum, nextP))

        return res


print(Solution().kthSmallest(mat=[[1, 3, 11], [2, 4, 6]], k=5))
# 输出：7
# 解释：从每一行中选出一个元素，前 k 个和最小的数组分别是：
# [1,2], [1,4], [3,2], [3,4], [1,6]。其中第 5 个的和是 7 。

