from typing import List

# 2 <= n <= 250
# 你可以执行以下操作 任意次 ：
# 选择 matrix 中 相邻 两个元素，并将它们都 乘以 -1 。
# 你的目的是 最大化 方阵元素的和。请你在执行以上操作之后，返回方阵的 最大 和。


# 总结：
# 只需要确定负数的数量是奇数还是偶数就好。
# 如果是偶数，经过有限次变换一定能把负数都抵消掉
# 如果是奇数，让绝对值最小的那个数为负数，对和的影响最小。


class Solution:
    def maxMatrixSum(self, matrix: List[List[int]]) -> int:
        res = minusCount = 0
        minVal = int(1e5 + 1)
        for i in range(len(matrix)):
            for j in range(len(matrix)):
                res += abs(matrix[i][j])
                minVal = min(minVal, abs(matrix[i][j]))
                if matrix[i][j] < 0:
                    minusCount ^= 1
        return res - 2 * minusCount * minVal


print(Solution().maxMatrixSum(matrix=[[1, -1], [-1, 1]]))
# 输出：4
# 解释：我们可以执行以下操作使和等于 4 ：
# - 将第一行的 2 个元素乘以 -1 。
# - 将第一列的 2 个元素乘以 -1 。
