from collections import defaultdict
from typing import List


class Solution:
    def getTreeInfo(self, n: int, edges: List[List[int]], values: List[int]) -> None:
        """获取树的信息 模板示例"""

        def dfs(cur: int, pre: int, dep: int) -> None:
            parent[cur] = pre
            depth[cur] = dep
            for next in adjMap[cur]:
                if next == pre:
                    continue
                ancestors[next] |= ancestors[cur] | {cur}  # !先序遍历传递祖先节点
                dfs(next, cur, dep + 1)
                subCount[cur] += subCount[next]
                subSum[cur] += subSum[next]

            print('当前结点已处理完毕，准备向上回溯', subCount[cur], subSum[cur])

        depth = [-1] * n  # 到根节点的距离/深度
        parent = [-1] * n  # 直接的父节点
        ancestors = [set() for _ in range(n)]  # 祖先节点
        subSum = values[:]
        subCount = [1] * n

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dfs(0, -1, 0)

