"""
一条 好路径 需要满足以下条件：

开始节点和结束节点的值 相同 。
开始节点和结束节点中间的所有节点值都小于等于开始节点的值
（也就是说开始节点的值应该是路径上所有节点的最大值）。
请你返回不同好路径的数目。

注意，一条路径和它反向的路径算作 同一 路径。
比方说， 0 -> 1 与 1 -> 0 视为同一条路径。单个节点也视为一条合法路径。

n<=3e4 暗示nlogn
"""

# !启发式合并
# !后序dfs统计子树内的每种结点个数,同时删除不合法的点,再统计经过当前结点能产生多少条新路径
# !合并子树返回值时，小的dict合并到大的dict上去(启发式合并)
# !复杂度nlognlogn

from typing import List
from sortedcontainers import SortedDict


class Solution:
    def numberOfGoodPaths(self, vals: List[int], edges: List[List[int]]) -> int:
        def dfs(cur: int, pre: int) -> "SortedDict":
            """后序dfs返回子树内的每种结点个数 在当前结点处统计经过当前结点能产生多少条新路径"""
            self.res += 1
            curRes = SortedDict({vals[cur]: 1})

            for next in adjList[cur]:
                if next == pre:
                    continue
                nextRes = dfs(next, cur)
                # 我们枚举到的每个元素都会被删掉，所以总共最多只有 n 次删除，复杂度还是正确的。
                while nextRes and nextRes.peekitem(0)[0] < vals[cur]:  # type: ignore 启发式合并
                    nextRes.popitem(0)

                if len(curRes) < len(nextRes):
                    curRes, nextRes = nextRes, curRes
                for key in nextRes:
                    self.res += curRes.get(key, 0) * nextRes[key]
                    curRes[key] = curRes.get(key, 0) + nextRes[key]

            return curRes

        n = len(vals)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        self.res = 0
        dfs(0, -1)
        return self.res


print(Solution().numberOfGoodPaths([1, 3, 3, 5], [[0, 1], [0, 2], [3, 1]]))
