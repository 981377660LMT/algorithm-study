# 树的最长路(有向/无向均可)


from heapq import nlargest
from typing import List


# 类似求树的直径/树的中心的解法
# 最后的路径一定可以表示为 : u1 <- target -> u2 , u1 , u2 均为 target 走向的一条路径 ,
# 且 u1 , u2 可以为空。
# 那么我们枚举 target 即可 , 在 DFS 的过程中 , 把 当前点当成 target 节点即可 ,
# 那么以 target 为 "中心" 的最长的路径一定是它的往下的合法的路径的 最长 和 次长 的路径和 + 1。
# 然后返回从 target 往下的最长的路径 + 1 即可。


class Solution:
    def longestPath(self, parent: List[int], s: str) -> int:
        def dfs(cur: int, pre: int) -> int:
            """后序dfs求每个root处向下的次长路和最长路"""
            nonlocal res
            nexts = [0, 0]
            for next in adjList[cur]:
                if next == pre:
                    continue
                nextRes = dfs(next, cur)
                if s[next] != s[cur]:  # 只有不同的时候才能加入
                    nexts.append(nextRes)
            max1, max2 = nlargest(2, nexts)
            res = max(res, max1 + max2 + 1)
            return 1 + max1

        n = len(parent)
        adjList = [[] for _ in range(n)]
        for cur, pre in enumerate(parent):
            if pre == -1:
                continue
            adjList[pre].append(cur)
            adjList[cur].append(pre)

        res = 1
        dfs(0, -1)
        return res


assert Solution().longestPath(parent=[-1, 0, 0, 1, 1, 2], s="abacbe") == 3
