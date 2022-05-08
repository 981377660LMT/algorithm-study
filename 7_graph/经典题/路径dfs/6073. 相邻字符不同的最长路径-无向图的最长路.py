from functools import lru_cache
from heapq import nlargest
from typing import List
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= n <= 105
# parent[0] == -1


class Solution:
    def longestPath(self, parent: List[int], s: str) -> int:
        """请你找出路径上任意一对相邻节点都没有分配到相同字符的 最长路径 ，并返回该路径的长度。"""

        @lru_cache(None)
        def dfs(cur: int, pre: int) -> int:
            res = 1
            for next in adjMap[cur]:
                if next == pre:
                    continue
                if s[next] == s[cur]:
                    continue
                res = max(res, dfs(next, cur) + 1)
            return res

        n = len(parent)
        adjMap = defaultdict(set)
        for i in range(n):
            pre, cur = parent[i], i
            if pre == -1:
                continue
            adjMap[pre].add(cur)
            adjMap[cur].add(pre)

        res = 1
        for i in range(n):
            res = max(res, dfs(i, -1))
        dfs.cache_clear()
        return res


print(Solution().longestPath(parent=[-1, 0, 0, 1, 1, 2], s="abacbe"))
print(Solution().longestPath(parent=[-1, 0, 0, 0], s="aabc"))

####################################################################################
# 类似求树的直径/树的中心的解法
# 最后的路径一定可以表示为 : u1 <- target -> u2 , u1 , u2 均为 target 走向的一条路径 , 且 u1 , u2 可以为空。
# 那么我们枚举 target 即可 , 在 DFS 的过程中 , 把 当前点当成 target 节点即可 , 那么以 target 为 "中心" 的最长的路径一定是它的往下的合法的路径的 最长 和 次长 的路径和 + 1。
# 然后返回从 target 往下的最长的路径 + 1 即可。


class Solution2:
    def longestPath(self, parent: List[int], s: str) -> int:
        def dfs(cur: int, pre: int) -> int:
            """后序dfs求每个root处向下的次长路和最长路"""
            nonlocal res

            for next in adjMap[cur]:
                if next == pre:
                    continue

                nextRes = dfs(next, cur)

                if s[next] != s[cur]:
                    if nextRes > down1[cur]:
                        down1[cur], down2[cur] = nextRes, down1[cur]
                    elif nextRes > down2[cur]:
                        down2[cur] = nextRes

            res = max(res, down1[cur] + down2[cur] + 1)
            return 1 + down1[cur]

        n = len(parent)
        adjMap = defaultdict(set)
        for i in range(n):
            pre, cur = parent[i], i
            if pre == -1:
                continue
            adjMap[pre].add(cur)
            adjMap[cur].add(pre)

        # 分别记录向下的最大值和次大值
        down1, down2 = [0] * n, [0] * n
        res = 1
        dfs(0, -1)
        return res


class Solution3:
    def longestPath(self, parent: List[int], s: str) -> int:
        def dfs(cur: int, pre: int) -> int:
            """后序dfs求每个root处向下的次长路和最长路"""
            nonlocal res

            nexts = [0, 0]
            for next in adjMap[cur]:
                if next == pre:
                    continue
                nextRes = dfs(next, cur)
                if s[next] != s[cur]:
                    nexts.append(nextRes)

            max1, max2 = nlargest(2, nexts)
            # down1[cur], down2[cur] = max1, max2
            res = max(res, max1 + max2 + 1)
            return 1 + max1

        n = len(parent)
        adjMap = defaultdict(set)
        for i in range(n):
            pre, cur = parent[i], i
            if pre == -1:
                continue
            adjMap[pre].add(cur)
            adjMap[cur].add(pre)

        res = 1
        # down1, down2 = [0] * n, [0] * n
        dfs(0, -1)
        return res
