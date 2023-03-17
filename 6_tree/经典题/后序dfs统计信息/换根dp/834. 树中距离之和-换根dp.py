from collections import defaultdict
from typing import List

# 834. 树中距离之和-换根dp


class Solution:
    def sumOfDistancesInTree2(self, n: int, edges: List[List[int]]) -> List[int]:
        """换根dp求每个节点到其他节点的距离之和"""
        from Rerooting import Rerooting

        def e(root: int) -> int:
            return 0

        def op(childRes1: int, childRes2: int) -> int:
            return childRes1 + childRes2

        def composition(fromRes: int, parent: int, cur: int, direction: int) -> int:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            if direction == 0:  # cur -> parent
                return fromRes + subSize[cur]
            return fromRes + (n - subSize[cur])  # parent -> cur

        def dfsForSubSize(cur: int, parent: int) -> None:
            for next in R.adjList[cur]:
                if next != parent:
                    dfsForSubSize(next, cur)
                    subSize[cur] += subSize[next]

        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        subSize = [1] * n
        dfsForSubSize(0, -1)
        dp = R.rerooting(e=e, op=op, composition=composition)
        return dp

    def sumOfDistancesInTree(self, n: int, edges: List[List[int]]) -> List[int]:
        # !子结点更新父结点向下的距离 求出根0的答案
        def dfs1(cur: int, parent_: int, depth_: int) -> None:

            parent[cur] = parent_
            depth[cur] = depth_
            for next in adjMap[cur]:
                if next == parent_:
                    continue
                dfs1(next, cur, depth_ + 1)
                subTreeCount[cur] += subTreeCount[next]

        # !父结点更新子结点向上的距离
        def dfs2(cur: int, parent: int) -> None:
            for next in adjMap[cur]:
                if next == parent:
                    continue
                # 注意这里都是 subTreeCount[next]
                res[next] = res[cur] - subTreeCount[next] + (n - subTreeCount[next])
                dfs2(next, cur)

        depth = [-1] * n
        parent = [-1] * n
        subTreeCount = [1] * n

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        dfs1(0, -1, 0)  # 求出根0的答案

        res = [0] * n
        res[0] = sum(depth)
        dfs2(0, -1)
        return res


print(Solution().sumOfDistancesInTree(n=6, edges=[[0, 1], [0, 2], [2, 3], [2, 4], [2, 5]]))
print(Solution().sumOfDistancesInTree2(n=6, edges=[[0, 1], [0, 2], [2, 3], [2, 4], [2, 5]]))
