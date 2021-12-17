# 我们可以从初始列表中`删除一个节点`，并完全移除该节点以及从该节点到任何其他节点的任何连接。

# 并查集本身只适合用作集合的合并，并不适合用作集合的拆分
# 当碰到`拆分` 并查集的题干，应该想到 逆向思维 地利用并查集
# 将问题由删除一个感染节点能减少多个节点受到感染转换成添加一个感染节点会增加多少个被感染节点，
# 即添加一个感染节点，使得该节点能感染的节点最多。


# 此题与924区别是断开了点的连接,影响了连通分量的个数
# 并查集并不适合删除点,因此考虑反向添加点

# https://leetcode.com/problems/minimize-malware-spread-ii/discuss/1561512/Python3-union-find
# https://leetcode-cn.com/problems/minimize-malware-spread-ii/solution/bing-cha-ji-mo-ban-by-yexiso-w7t9/
from typing import List
from collections import Counter, defaultdict


class UnionFind:
    def __init__(self, n):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x):
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        root_x = self.find(x)
        root_y = self.find(y)
        if root_x == root_y:
            return False
        if self.size[root_x] > self.size[root_y]:
            root_x, root_y = root_y, root_x
        self.parent[root_x] = root_y
        self.size[root_y] += self.size[root_x]
        self.part -= 1
        return True

    def isconnected(self, p, q):
        return self.find(p) == self.find(q)


class Solution:
    def minMalwareSpread(self, graph: List[List[int]], initial: List[int]) -> int:
        n = len(graph)
        uf = UnionFind(n)
        evil = set(initial)

        # 1.忽略所有感染节点，只考虑正常节点。
        for i in range(len(graph)):
            if i in evil:
                continue
            for j in range(i + 1, len(graph)):
                if j in evil:
                    continue
                if graph[i][j] == 1:
                    uf.union(i, j)

        # 看每个感染源感染了哪些门派
        infectRange = defaultdict(set)
        for u in initial:
            for v in range(n):
                if v in evil:
                    continue
                if graph[u][v] == 1:
                    infectRange[u].add(uf.find(v))

        # 看每个门派的感染次数(被几个感染源感染)
        freq = sum((Counter(v) for v in infectRange.values()), Counter())

        res = min(initial)
        best = -1
        for u in initial:
            validCount = 0
            for v in infectRange[u]:
                # v门派只被u感染,有效的感染了多少人
                if freq[v] == 1:
                    validCount += uf.size[v]
            if validCount > best or validCount == best and u < res:
                res, best = u, validCount
        return res


print(Solution().minMalwareSpread(graph=[[1, 1, 0], [1, 1, 0], [0, 0, 1]], initial=[0, 1]))
