from typing import List
from heapq import heappop, heappush

# 矩阵中坐标 (a, b) 的 值 可由对所有满足 0 <= i <= a < m 且 0 <= j <= b < n 的元素 matrix[i][j]（下标从 0 开始计数）执行异或运算得到。
# 请你找出 matrix 的所有坐标中第 k 大的值（k 的值从 1 开始计数）。

# 前缀
#  xor[i][j] = xor[i-1][j] ^ xor[i][j-1] ^ xor[i-1][j-1] ^ matrix[i][j]
# O(MNlogK)
class Solution:
    def kthLargestValue(self, matrix: List[List[int]], k: int) -> int:
        m, n = len(matrix), len(matrix[0])  # dimensions
        pq = []

        for i in range(m):
            for j in range(n):
                if i > 0:
                    matrix[i][j] ^= matrix[i - 1][j]
                if j > 0:
                    matrix[i][j] ^= matrix[i][j - 1]
                if i > 0 and j > 0:
                    matrix[i][j] ^= matrix[i - 1][j - 1]
                heappush(pq, matrix[i][j])
                if len(pq) > k:
                    heappop(pq)
        return pq[0]


print(Solution().kthLargestValue(matrix=[[5, 2], [1, 6]], k=1))
# 输入：matrix = [[5,2],[1,6]], k = 1
# 输出：7
# 解释：坐标 (0,1) 的值是 5 XOR 2 = 7 ，为最大的值。
# 输入：matrix = [[5,2],[1,6]], k = 1
# 输出：7
# 解释：坐标 (0,1) 的值是 5 XOR 2 = 7 ，为最大的值。
