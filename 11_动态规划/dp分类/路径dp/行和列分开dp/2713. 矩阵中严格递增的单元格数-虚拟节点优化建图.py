# https://leetcode.cn/problems/maximum-strictly-increasing-cells-in-a-matrix/solution/you-xiang-wu-huan-tu-de-zui-chang-lu-jia-ssrq/
# 6456. 矩阵中严格递增的单元格数
# https://atcoder.jp/contests/abc224/tasks/abc224_e
# E - Integers on Grid (爬山)
# 有n个方格上有正整数，其余的方格上的数字为0
# 在方格上走，每次只能走同一行或同一列，
# 且到达的格点上的值必须严格大于(真に大きい)当前值，
# 问从给定的n个点出发 最多走多少步
# !ROW,COL<=2e5 n<=2e5


# 1. 同一行列只连接相邻的点
# 2. 虚拟结点优化多个结点的相互连接(2*ROW*COL个虚拟结点)


from functools import lru_cache
from itertools import groupby
from typing import List


class Solution:
    def maxIncreasingCells(self, mat: List[List[int]]) -> int:
        ROW, COL = len(mat), len(mat[0])
        all_ = ROW * COL
        adjList = [[] for _ in range(3 * all_)]
        dummy = all_
        for r, row in enumerate(mat):
            pair = [(v, r * COL + c) for c, v in enumerate(row)]
            pair.sort()
            group = [[id for _, id in group] for _, group in groupby(pair, key=lambda x: x[0])]
            for pre, cur in zip(group, group[1:]):
                for id in pre:
                    adjList[id].append(dummy)
                for id in cur:
                    adjList[dummy].append(id)
                dummy += 1

        for c in range(COL):
            pair = [(mat[r][c], r * COL + c) for r in range(ROW)]
            pair.sort()
            group = [[id for _, id in group] for _, group in groupby(pair, key=lambda x: x[0])]
            for pre, cur in zip(group, group[1:]):
                for id in pre:
                    adjList[id].append(dummy)
                for id in cur:
                    adjList[dummy].append(id)
                dummy += 1

        @lru_cache(None)
        def dfs(cur: int) -> int:
            """有向图最长路."""
            res = 1
            for next in adjList[cur]:
                res = max(res, dfs(next) + 1)
            return res

        res = 0
        for r in range(ROW):
            for c in range(COL):
                res = max(res, dfs(r * COL + c))
        dfs.cache_clear()
        return (res + 1) // 2
