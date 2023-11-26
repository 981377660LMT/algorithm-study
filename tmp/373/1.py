from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个大小为 m x n 的整数矩阵 mat 和一个整数 k 。请你将矩阵中的 奇数 行循环 右 移 k 次，偶数 行循环 左 移 k 次。

# 如果初始矩阵和最终矩阵完全相同，则返回 true ，否则返回 false 。


class Solution:
    def areSimilar(self, mat: List[List[int]], k: int) -> bool:
        target = [list(row) for row in mat]
        for i in range(len(mat)):
            curRow = mat[i].copy()
            if i % 2 == 0:
                for j in range(len(curRow)):
                    target[i][j] = curRow[(j + k) % len(curRow)]
            else:
                for j in range(len(curRow)):
                    target[i][j] = curRow[(j - k) % len(curRow)]

        return all([target[i][j] == mat[i][j] for i in range(len(mat)) for j in range(len(mat[0]))])


# [[4,9,10,10],[9,3,8,4],[2,5,3,8],[6,1,10,4]]
# 5

print(Solution().areSimilar([[4, 9, 10, 10], [9, 3, 8, 4], [2, 5, 3, 8], [6, 1, 10, 4]], 5))
