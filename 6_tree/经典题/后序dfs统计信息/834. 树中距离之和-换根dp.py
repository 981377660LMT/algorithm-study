from collections import defaultdict
from typing import List

# 834. 树中距离之和-换根dp
class Solution:
    def sumOfDistancesInTree(self, n: int, edges: List[List[int]]) -> List[int]:
        def dfs(cur: int, parent_: int, depth_: int) -> None:
            parent[cur] = parent_
            depth[cur] = depth_
            for next in adjMap[cur]:
                if next == parent_:
                    continue
                dfs(next, cur, depth_ + 1)
                subTreeCount[cur] += subTreeCount[next]

        def getRes(cur: int, parent: int) -> None:
            for next in adjMap[cur]:
                if next == parent:
                    continue
                # 注意这里都是 subTreeCount[next]
                res[next] = res[cur] - subTreeCount[next] + (n - subTreeCount[next])
                getRes(next, cur)

        depth = [-1] * n
        parent = [-1] * n
        subTreeCount = [1] * n

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dfs(0, -1, 0)

        res = [0] * n
        res[0] = sum(depth)
        getRes(0, -1)
        return res


print(Solution().sumOfDistancesInTree(n=6, edges=[[0, 1], [0, 2], [2, 3], [2, 4], [2, 5]]))

