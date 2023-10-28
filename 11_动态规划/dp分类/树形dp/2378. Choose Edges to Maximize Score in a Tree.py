"""树形dp选边

11_动态规划/dp分类/树上dp/F - Select Edges.py 的特例
"""


from collections import defaultdict
from typing import List, Tuple

INF = int(4e18)


class Solution:
    def maxScore(self, edges: List[List[int]]) -> int:
        """在树中选择一些边使得选出的边不相连，最大化边权之和"""

        def dfs(cur: int, pre: int) -> Tuple[int, int]:
            """返回[选择连接父亲的边时子树最大价值, 不选择连接父亲的边时子树最大价值]

            如果选到父亲的边,所有儿子都不能选
            如果不选到父亲的边,所有儿子中选一个diff最大的
            """
            res1, res2 = 0, 0
            diff = [0]
            for next in adjMap[cur]:
                if next == pre:
                    continue
                select, jump = dfs(next, cur)
                res1, res2 = res1 + jump, res2 + jump
                diff.append(select + adjMap[cur][next] - jump)

            res2 += max(diff)
            return res1, res2

        adjMap = defaultdict(lambda: defaultdict(lambda: -INF))
        for cur, (parent, weight) in enumerate(edges):
            if parent == -1:
                continue
            adjMap[parent][cur] = weight
            adjMap[cur][parent] = weight

        return dfs(0, -1)[1]  # 不连接虚拟根节点


# edges[i] = [parenti, weighti]
print(Solution().maxScore(edges=[[-1, -1], [0, 5], [0, 10], [2, 6], [2, 4]]))
# 11
