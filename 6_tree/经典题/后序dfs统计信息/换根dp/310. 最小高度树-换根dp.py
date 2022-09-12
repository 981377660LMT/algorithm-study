from collections import defaultdict
from typing import List

# 310. 最小高度树-换根dp
# https://luyuhuang.tech/2022/05/06/tree-dp.html


class Solution:
    def findMinHeightTrees(self, n: int, edges: List[List[int]]) -> List[int]:
        """
        每个点作为根结点,看具有最小高度的树。
        找到所有的 最小高度树 并按 任意顺序 返回它们的根节点标签列表。
        """

        # !子结点更新父结点向下的最远距离 求出根0的答案
        def dfs1(cur: int, parent: int) -> int:
            """后序dfs,down1/down2数组记录一下每个结点往下的最远距离、次远距离,返回每个root处的最大路径长度"""
            for next, weight in adjMap[cur]:
                if next == parent:
                    continue

                maxCand = dfs1(next, cur) + weight
                if maxCand > down1[cur]:
                    down2[cur], down1[cur] = down1[cur], maxCand
                    downMaxNeedRoot[cur] = next
                elif maxCand > down2[cur]:
                    down2[cur] = maxCand

            return down1[cur]

        # !父结点更新子结点向上的最远距离
        def dfs2(cur: int, parent: int) -> None:
            """前序dfs,利用父结点来更新子结点"""
            # 每个点处最长和次长
            for next, weight in adjMap[cur]:
                if next == parent:
                    continue
                if downMaxNeedRoot[cur] == next:
                    # 另一条次长路
                    up[next] = max(up[cur], down2[cur]) + weight
                else:
                    up[next] = max(up[cur], down1[cur]) + weight
                dfs2(next, cur)

        # 分别记录`向下(子节点)`的最大值和次大值dp1
        down1, down2 = [0] * n, [0] * n
        # `向下`取最大值时必须经过的结点
        downMaxNeedRoot = [0] * n
        # 记录`节点向上(父结点)`的最大距离dp2
        up = [0] * n

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add((v, 1))
            adjMap[v].add((u, 1))

        dfs1(0, -1)
        dfs2(0, -1)
        minDists = [max(up, down) for down, up in zip(down1, up)]
        min_ = min(minDists)
        return [i for i, dist in enumerate(minDists) if dist == min_]


print(Solution().findMinHeightTrees(n=6, edges=[[3, 0], [3, 1], [3, 2], [3, 4], [5, 4]]))
