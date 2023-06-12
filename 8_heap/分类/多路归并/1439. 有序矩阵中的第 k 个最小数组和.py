# 给你一个 m * n 的矩阵 mat，以及一个整数 k ，矩阵中的每一行都以非递减的顺序排列。
# 你可以从每一行中选出 1 个元素形成一个数组。返回所有可能数组中的第 k 个 最小 数组和。

# 1 <= m, n <= 40
# 1 <= k <= min(200, n ^ m)
# 1 <= mat[i][j] <= 5000
# mat[i] 是一个非递减数组

from functools import reduce
from typing import List
from kSmallestPairs import kSmallestPairs2


class Solution:
    def kthSmallest(self, mat: List[List[int]], k: int) -> int:
        """时间复杂度O(row*k*logmin(col,k))"""
        return reduce(lambda x, y: kSmallestPairs2(x, y, k), mat)[-1]


print(Solution().kthSmallest(mat=[[1, 3, 11], [2, 4, 6]], k=5))
# 输出：7
# 解释：从每一行中选出一个元素，前 k 个和最小的数组分别是：
# [1,2], [1,4], [3,2], [3,4], [1,6]。其中第 5 个的和是 7 。
