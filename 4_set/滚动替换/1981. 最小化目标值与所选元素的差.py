from typing import List

# 从矩阵的 每一行 中选择一个整数，你的目标是 最小化 所有选中元素之 和 与目标值 target 的 绝对差 。
# 返回 最小的绝对差 。
# 1 <= m, n <= 70
# mat[i][j] <= 70,
# 表明最后的和不超过4900 => 对矩阵每一行，求和并更新
class Solution:
    def minimizeTheDifference(self, mat: List[List[int]], target: int) -> int:
        sums = set([0])
        for row in mat:
            row.sort()

        for row in mat:
            tmp = set()
            for s in sums:
                for v in row:
                    tmp.add(s + v)
                    if s + v >= target:
                        break
            sums = tmp

        return min(abs(s - target) for s in sums)


print(Solution().minimizeTheDifference(mat=[[1, 2, 3], [4, 5, 6], [7, 8, 9]], target=13))
# 输入：mat = [[1,2,3],[4,5,6],[7,8,9]], target = 13
# 输出：0
# 解释：一种可能的最优选择方案是：
# - 第一行选出 1
# - 第二行选出 5
# - 第三行选出 7
# 所选元素的和是 13 ，等于目标值，所以绝对差是 0 。

