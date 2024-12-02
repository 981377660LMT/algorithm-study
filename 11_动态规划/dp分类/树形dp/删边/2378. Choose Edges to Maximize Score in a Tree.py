"""树形dp选边(删边、删除树边、树上删边)

11_动态规划/dp分类/树上dp/F - Select Edges.py 的特例
"""

from typing import List, Tuple

INF = int(4e18)


class Solution:
    def maxScore(self, edges: List[List[int]]) -> int:
        """在树中选择一些边使得选出的边不相连，最大化边权之和."""

        def dfs(cur: int, pre: int) -> Tuple[int, int]:
            """
            返回[选择连接父亲的边时子树最大价值, 不选择连接父亲的边时子树最大价值].

            如果选到父亲的边,所有儿子都不能选.
            如果不选到父亲的边,所有儿子中选一个diff最大的.
            """
            res1, res2 = 0, 0
            diff = [0]
            for next, weight in adjList[cur]:
                if next == pre:
                    continue
                select, skip = dfs(next, cur)
                res1, res2 = res1 + skip, res2 + skip
                diff.append(select + weight - skip)

            res2 += max(diff)
            return res1, res2

        n = len(edges) + 1
        adjList = [[] for _ in range(n)]
        for cur, (parent, weight) in enumerate(edges):
            if parent == -1:
                continue
            adjList[parent].append((cur, weight))
            adjList[cur].append((parent, weight))

        return dfs(0, -1)[1]  # 不连接虚拟根节点
