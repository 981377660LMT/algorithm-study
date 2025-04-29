# 766. 托普利茨矩阵
# https://leetcode.cn/problems/toeplitz-matrix/description/
# !托普利茨矩阵（Toeplitz Matrix）是一种每条从左上到右下的对角线元素都相同的矩阵。
# 托普利茨矩阵常用于需要高效处理具有平移不变性或卷积结构的问题。
# 进阶：
#
# 1. 如果矩阵存储在磁盘上，并且内存有限，以至于一次最多只能将矩阵的一行加载到内存中，该怎么办？
#    !我们将每一行复制到一个连续数组中，随后在读取下一行时，就与内存中此前保存的数组进行比较.
# 2. 如果矩阵太大，以至于一次只能将不完整的一行加载到内存中，该怎么办？
#    !将整个矩阵竖直切分成若干子矩阵，并保证两个相邻的矩阵至少有一列或一行是重合的，然后判断每个子矩阵是否符合要求


from typing import List


class Solution:
    def isToeplitzMatrix(self, matrix: List[List[int]]) -> bool:
        n, m = len(matrix), len(matrix[0])
        for i in range(1, n):
            for j in range(1, m):
                if matrix[i][j] != matrix[i - 1][j - 1]:
                    return False
        return True
