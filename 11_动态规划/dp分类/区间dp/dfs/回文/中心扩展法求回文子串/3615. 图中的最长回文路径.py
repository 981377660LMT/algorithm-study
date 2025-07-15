# 3615. 图中的最长回文路径(中心扩展法 + 状压 DP)
# https://leetcode.cn/problems/longest-palindromic-path-in-graph/solutions/3722469/zhong-xin-kuo-zhan-fa-zhuang-ya-dp-by-en-ai9s/
# 给你一个整数 n 和一个包含 n 个节点的 无向图 ，节点编号从 0 到 n - 1，以及一个二维数组 edges，其中 edges[i] = [ui, vi] 表示节点 ui 和节点 vi 之间有一条边。
# 同时给你一个长度为 n 的字符串 label，其中 label[i] 是与节点 i 关联的字符。
# 你可以从任意节点开始，移动到任意相邻节点，每个节点 最多 访问一次。
# 返回通过访问一条路径，路径中 不包含重复 节点，所能形成的 最长回文串 的长度。
# n <= 14


from functools import lru_cache
from typing import List


class Solution:
    def maxLen(self, n: int, edges: List[List[int]], label: str) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        @lru_cache(None)
        def dfs(left: int, right: int, visited: int) -> int:
            """从left和right向两侧扩展(不含left、right)，visited记录已访问的节点，返回回文长度."""
            res = 0
            for a in adjList[left]:
                if visited & (1 << a):
                    continue
                for b in adjList[right]:
                    if a == b or visited & (1 << b):
                        continue
                    if label[a] == label[b]:
                        na, nb = a, b
                        if na < nb:
                            na, nb = nb, na
                        res = max(res, dfs(na, nb, visited | (1 << a) | (1 << b)) + 2)
            return res

        res = 0
        for i, nexts in enumerate(adjList):
            res = max(res, dfs(i, i, 1 << i) + 1)  # 奇
            for j in nexts:
                if i < j and label[i] == label[j]:
                    res = max(res, dfs(i, j, (1 << i) | (1 << j)) + 2)  # 偶

        dfs.cache_clear()
        return res
