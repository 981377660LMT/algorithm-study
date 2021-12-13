from typing import List

# 第 0 行的元素之和为 upper。
# 第 1 行的元素之和为 lower。
# 第 i 列（从 0 开始编号）的元素之和为 colsum[i]，colsum 是一个长度为 n 的整数数组。

# 如果不存在符合要求的答案，就请返回一个空的二维数组。
class Solution:
    def reconstructMatrix(self, upper: int, lower: int, colsum: List[int]) -> List[List[int]]:
        n = len(colsum)
        row1 = [0] * n
        row2 = [0] * n

        # 贪心 谁大就分给谁1
        for col, val in enumerate(colsum):
            if val == 1:
                if upper > lower:
                    row1[col] = 1
                    upper -= 1
                else:
                    row2[col] = 1
                    lower -= 1
            elif val == 2:
                row1[col] = row2[col] = 1
                upper -= 1
                lower -= 1
            elif val >= 3:
                return []

        return [row1, row2] if upper == lower == 0 else []


print(Solution().reconstructMatrix(upper=2, lower=1, colsum=[1, 1, 1]))
# 输出：[[1,1,0],[0,0,1]]
# 解释：[[1,0,1],[0,1,0]] 和 [[0,1,1],[1,0,0]] 也是正确答案。
