from typing import List


# 叉乘模为0则三点共线
class Solution:
    def isBoomerang(self, points: List[List[int]]) -> bool:
        def cal_cross_product(A, B, C):
            AB = [B[0] - A[0], B[1] - A[1]]
            AC = [C[0] - A[0], C[1] - A[1]]
            return AB[0] * AC[1] - AB[1] * AC[0]

        return cal_cross_product(*points) != 0


print(Solution().isBoomerang([[1, 1], [2, 2], [3, 3]]))
print(Solution().isBoomerang([[1, 1], [2, 3], [3, 2]]))
