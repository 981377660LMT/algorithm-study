# 给你一个 m * n 的矩阵 mat，以及一个整数 k ，矩阵中的每一行都以非递减的顺序排列。
# 你可以从每一行中选出 1 个元素形成一个数组。返回所有可能数组中的第 k 个 最小 数组和。

# 1 <= m, n <= 40
# 1 <= k <= min(200, n ^ m)
# 1 <= mat[i][j] <= 5000
# mat[i] 是一个非递减数组
# https://leetcode.cn/problems/find-the-kth-smallest-sum-of-a-matrix-with-sorted-rows/

from typing import List


class Solution:
    def kthSmallest(self, mat: List[List[int]], k: int) -> int:
        """
        时间复杂度O(row*k*logU).
        二分+dfs,找到k个就不再递归.
        !注意这里的dfs起点是最小的选择.
        """

        def check(mid: int) -> bool:
            """选取的数组和不超过mid的个数是否>=k."""

            def dfs(index: int, curSum: int) -> None:
                nonlocal count
                if index == ROW:
                    count += 1
                    return
                if count >= k:  # !找到k个就不再递归
                    return
                for v in mat[index]:
                    cand = curSum + v - mat[index][0]  # 选v而不选mat[index][0]
                    if cand > mid:  # !剪枝,后面的都无法选择了
                        break
                    dfs(index + 1, cand)

            count = 0
            dfs(0, min_)  # !先默认选取每行的第一个元素，便于剪枝
            return count >= k

        ROW = len(mat)
        min_, max_ = sum(row[0] for row in mat), sum(row[-1] for row in mat)
        left, right = min_, max_
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


print(Solution().kthSmallest(mat=[[1, 3, 11], [2, 4, 6]], k=5))
# 输出：7
# 解释：从每一行中选出一个元素，前 k 个和最小的数组分别是：
# [1,2], [1,4], [3,2], [3,4], [1,6]。其中第 5 个的和是 7 。
