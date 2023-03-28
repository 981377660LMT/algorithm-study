# 1605. 给定行和列的和求矩阵
# 其中 rowSum[i] 是二维矩阵中第 i 行元素的和， colSum[j] 是第 j 列元素的和
# 请找到大小为 rowSum.length x colSum.length 的任意 非负整数 矩阵，且该矩阵满足 rowSum 和 colSum 的要求。

# !思路:遍历 matrix，每次都填这次能填的最大数, min(rowSum[i], colSum[j])
from typing import List


class Solution:
    def restoreMatrix(self, rowSum: List[int], colSum: List[int]) -> List[List[int]]:
        m, n = len(rowSum), len(colSum)
        res = [[0] * n for _ in range(m)]

        for r in range(m):
            for c in range(n):
                minVal = min(rowSum[r], colSum[c])
                res[r][c] = minVal
                rowSum[r] -= minVal
                colSum[c] -= minVal
                if rowSum[r] == 0 and colSum[c] == 0:
                    break

        return res


print(Solution().restoreMatrix(rowSum=[3, 8], colSum=[4, 7]))
# 输出：[[3,0],
#       [1,7]]
# 解释：
# 第 0 行：3 + 0 = 3 == rowSum[0]
# 第 1 行：1 + 7 = 8 == rowSum[1]
# 第 0 列：3 + 1 = 4 == colSum[0]
# 第 1 列：0 + 7 = 7 == colSum[1]
# 行和列的和都满足题目要求，且所有矩阵元素都是非负的。
# 另一个可行的矩阵为：[[1,2],
#                   [3,5]]
