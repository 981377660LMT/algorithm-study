from functools import lru_cache
from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 1 开始、大小为 m x n 的整数矩阵 mat，你可以选择任一单元格作为 起始单元格 。

# 从起始单元格出发，你可以移动到 同一行或同一列 中的任何其他单元格，但前提是目标单元格的值 严格大于 当前单元格的值。

# 你可以多次重复这一过程，从一个单元格移动到另一个单元格，直到无法再进行任何移动。

# 请你找出从某个单元开始访问矩阵所能访问的 单元格的最大数量 。

# 返回一个表示可访问单元格最大数量的整数。


class Solution:
    def maxIncreasingCells(self, mat: List[List[int]]) -> int:
        # 减少移动的边数
        ROW, COL = len(mat), len(mat[0])
        all_ = ROW * COL
        adjList = [[] for _ in range(3 * all_)]
        dummy = all_
        for r, row in enumerate(mat):
            curRow = sorted([(v, r * COL + j) for j, v in enumerate(row)])
            groups = [
                (char, [id for _, id in group])
                for char, group in groupby(curRow, key=lambda x: x[0])
            ]
            for (_, id1), (_, id2) in zip(groups, groups[1:]):
                for id in id1:
                    adjList[id].append(dummy)
                for id in id2:
                    adjList[dummy].append(id)
                dummy += 1
        for c in range(COL):
            colInfo = [(mat[r][c], r * COL + c) for r in range(ROW)]
            colInfo.sort()
            groups = [
                (char, [id for _, id in group])
                for char, group in groupby(colInfo, key=lambda x: x[0])
            ]  # !这里需要引入虚拟节点减少边数

            for (_, id1), (_, id2) in zip(groups, groups[1:]):
                for id in id1:
                    adjList[id].append(dummy)
                for id in id2:
                    adjList[dummy].append(id)
                dummy += 1

        @lru_cache(None)
        def dfs(cur: int) -> int:
            res = 1
            for next in adjList[cur]:
                res = max(res, dfs(next) + 1)
            return res

        res = max(dfs(r * COL + c) for r in range(ROW) for c in range(COL))

        return (res + 1) // 2


# mat = [[3,1,6],[-9,5,7]]
# print(Solution().maxIncreasingCells([[3, 1, 6], [-9, 5, 7]]))
# [[2,-4,2,-2,6]]
# print(Solution().maxIncreasingCells([[2, -4, 2, -2, 6]]))
# # [[3,1],[3,4]]
# print(Solution().maxIncreasingCells([[3, 1], [3, 4]]))
# # [[1,1],[1,1]]
# print(Solution().maxIncreasingCells([[1, 1], [1, 1]]))
# [[3,1,6],[-9,5,7]]
print(Solution().maxIncreasingCells([[3, 1, 6], [-9, 5, 7]]))
