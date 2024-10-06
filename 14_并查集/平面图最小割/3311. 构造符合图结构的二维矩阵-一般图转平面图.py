# 3311. 构造符合图结构的二维矩阵(拼图)
# https://leetcode.cn/problems/construct-2d-grid-matching-graph-layout/description/
#
# 一般图转平面图(网格图)
# 给你一个二维整数数组 edges ，它表示一棵 n 个节点的 无向 图，其中 edges[i] = [ui, vi] 表示节点 ui 和 vi 之间有一条边。
#
# 请你构造一个二维矩阵，满足以下条件：
#
# 矩阵中每个格子 一一对应 图中 0 到 n - 1 的所有节点。
# 矩阵中两个格子相邻（横 的或者 竖 的）当且仅当 它们对应的节点在 edges 中有边连接。
# 题目保证 edges 可以构造一个满足上述条件的二维矩阵。
#
# 请你返回一个符合上述要求的二维整数数组，如果存在多种答案，返回任意一个。
#
# !构造出第一行后，剩下行唯一确定.
# !三种情况分类讨论.

from typing import List
from collections import defaultdict


class Solution:
    def constructGridLayout(self, n: int, edges: List[List[int]]) -> List[List[int]]:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
        degGroup = defaultdict(list)
        for i, nexts in enumerate(adjList):
            degGroup[len(nexts)].append(i)

        def find(cur: int, deg: int) -> int:
            """找到与cur相邻的节点中度数为deg的节点."""
            for v in adjList[cur]:
                if len(adjList[v]) == deg:
                    return v
            return -1

        if min(degGroup) == 1:  # 一列
            row = [degGroup[1][0]]
        elif max(degGroup) <= 3:  # 两列
            a = degGroup[2][0]
            b = find(a, deg=2)
            row = [a, b]
        else:  # 至少三列
            a = degGroup[2][0]
            row = [a]
            pre, cur = a, find(a, deg=3)
            while len(adjList[cur]) > 2:
                row.append(cur)
                for next in adjList[cur]:
                    if next != pre and len(adjList[next]) < 4:
                        pre, cur = cur, next
                        break
            row.append(cur)

        res = [row]
        visited = [False] * n
        for _ in range(n // len(row) - 1):  # !一共有n//len(row)行, len(row)列
            for v in row:
                visited[v] = True
            nextRow = []
            for cur in row:
                for next in adjList[cur]:
                    if not visited[next]:
                        nextRow.append(next)
                        break
            res.append(nextRow)
            row = nextRow
        return res
