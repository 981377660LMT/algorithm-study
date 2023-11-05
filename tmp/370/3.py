from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 有一棵 n 个节点的无向树，节点编号为 0 到 n - 1 ，根节点编号为 0 。给你一个长度为 n - 1 的二维整数数组 edges 表示这棵树，其中 edges[i] = [ai, bi] 表示树中节点 ai 和 bi 有一条边。

# 同时给你一个长度为 n 下标从 0 开始的整数数组 values ，其中 values[i] 表示第 i 个节点的值。

# 一开始你的分数为 0 ，每次操作中，你将执行：

# 选择节点 i 。
# 将 values[i] 加入你的分数。
# 将 values[i] 变为 0 。
# 如果从根节点出发，到任意叶子节点经过的路径上的节点值之和都不等于 0 ，那么我们称这棵树是 健康的 。

# 你可以对这棵树执行任意次操作，但要求执行完所有操作以后树是 健康的 ，请你返回你可以获得的 最大分数 。


# TODO: 去首都、和这一题，总结
# 用needValue 也可以记忆化, 去首都能否记忆化


def max(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximumScoreAfterOperations(self, edges: List[List[int]], values: List[int]) -> int:
        n = len(values)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        def dfs2(cur: int, pre: int) -> Tuple[int, int]:
            """所有子树ok时的最大值,子树中不ok时的最大值."""
            if len(adjList[cur]) == 1 and pre != -1:
                return 0, values[cur]

            oks = []
            notOks = []
            for next in adjList[cur]:
                if next == pre:
                    continue
                ok, notOk = dfs2(next, cur)
                oks.append(ok)
                notOks.append(notOk)
            res1 = max(subValue[cur], sum(oks) + values[cur])
            res2 = sum(notOks) + values[cur]
            return res1, res2

        subValue = [0] * n
        return dfs2(0, -1)[0]


# edges = [[0,1],[0,2],[1,3],[1,4],[2,5],[2,6]], values = [20,10,9,7,4,3,5]

print(
    Solution().maximumScoreAfterOperations(
        [[0, 1], [0, 2], [1, 3], [1, 4], [2, 5], [2, 6]], [20, 10, 9, 7, 4, 3, 5]
    )
)
