from collections import defaultdict, deque
from typing import DefaultDict, List, Set

# LCA求树两点间的路径


class Solution:
    def closestNode(self, n: int, edges: List[List[int]], query: List[List[int]]) -> List[int]:
        def dfs(cur: int, pre: int, dep: int) -> None:
            """处理高度、父节点信息"""
            depth[cur], parent[cur] = dep, pre
            for next in adjMap[cur]:
                if next == pre:
                    continue
                dfs(next, cur, dep + 1)

        def getPath(
            root1: int, root2: int, level: DefaultDict[int, int], parent: DefaultDict[int, int]
        ) -> Set[int]:
            """求两个结点间的路径，不断上跳到LCA并记录经过的结点"""
            res = {root1, root2}
            if level[root1] < level[root2]:
                root1, root2 = root2, root1
            diff = level[root1] - level[root2]
            for _ in range(diff):
                root1 = parent[root1]
                res |= {root1}
            while root1 != root2:
                root1 = parent[root1]
                root2 = parent[root2]
                res |= {root1, root2}
            return res

        def bfs(start: int, hit: Set[int]) -> int:
            """求到目标路径的最近交点"""
            visited, queue = set([start]), deque([start])
            while queue:
                cur = queue.popleft()
                if cur in hit:
                    return cur
                for next in adjMap[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            raise Exception("impossible")

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        depth, parent = defaultdict(int), defaultdict(lambda: -1)
        dfs(0, -1, 0)
        return [bfs(root3, getPath(root1, root2, depth, parent)) for root1, root2, root3 in query]
