# 给你一个大小为 m x n 的矩阵 mat ，请以对角线遍历的顺序，用一个数组返回这个矩阵中的所有元素。
from collections import defaultdict
from typing import List


class Solution:
    def findDiagonalOrder(self, mat: List[List[int]]) -> List[int]:
        ROW, COL = len(mat), len(mat[0])
        adjMap = defaultdict(list)
        for r in range(ROW):
            for c in range(COL):
                adjMap[r + c].append(mat[r][c])
        res = []
        for key, vals in adjMap.items():
            res.extend(vals if key & 1 else vals[::-1])
        return res

