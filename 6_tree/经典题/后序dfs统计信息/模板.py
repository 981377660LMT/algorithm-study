from collections import defaultdict
from typing import List


class Solution:
    def getTreeInfo(self, n: int, edges: List[List[int]], values: List[int]) -> None:
        """获取树的信息"""

        def dfs(cur: int, parent_: int, depth_: int) -> None:
            parent[cur] = parent_
            depth[cur] = depth_
            for next in adjMap[cur]:
                if next == parent_:
                    continue
                dfs(next, cur, depth_ + 1)
                subTreeCount[cur] += subTreeCount[next]
                subTreeSum[cur] += subTreeSum[next]

            print('当前结点已处理完毕，准备向上回溯', subTreeCount[cur], subTreeSum[cur])

        depth = [-1] * n
        parent = [-1] * n
        subTreeSum = values[:]
        subTreeCount = [1] * n

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dfs(0, -1, 0)

