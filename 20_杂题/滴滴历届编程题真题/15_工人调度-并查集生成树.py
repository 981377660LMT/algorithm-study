# 某工厂有 N 名工人和 M 台机器，每名工人都有一个能力值，
# 且只懂得操作其中某两台机器。
# !另外，每名工人最多只允许操作一台机器，
# 且每台机器最多只允许被一名工人操作。
# 那么如何分配每名工人操作哪台机器(或者不操作机器)，
# 才能让所有操作机器的工人的能力值总和最大?
# n,m<=1e5
# !二分图最大权匹配 => 最小费用最大流 会超时

# 链接：https://leetcode.cn/problems/sRI8mk/solution/by-freeyourmind-k0uy/
# 贪心:
# !由于每台机器最多只能有1个工人操作，那么在空闲的工人中，能力值最高的就应该分配给这台机器
# 1. 按照工人能力值从大到小处理，每次用并查集把a, b两台机器合并到一起。
# 表示这两台有一台被该工人操作
# 2. 处理到某个工人，检查他可以操作的a, b两台机器所在的集合，是否还能再加工人，
# 可以加的话，这个工人的能力值就计入结果
# 并查集需要在合并的时候处理是否能加工人的问题，只需要判断：
# 两个机器所在的集合，如果机器总数比人多，那么可以加工人，否则不能加这个工人。

# !相当于把工人当做边，工人能操作的两台机器是这个边的两个点。
# 该算法就是生成一片权值最大的森林，
# 其中每个连通块是树或者基环树。
# !权值最大的生成树森林

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# !union函数很特别的并查集
class UnionFind:
    def __init__(self, n: int):
        self.parent = list(range(n))
        self.rank = [1] * n  # 每个连通块中的机器数
        self.person = [0] * n  # 每个连通块中的人数

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)

        # !如果两个连通块中的机器数大于人数，那么可以加人
        if self.rank[rootX] + self.rank[rootY] <= self.person[rootX] + self.person[rootY]:
            return False

        if rootX == rootY:
            self.person[rootX] += 1
        else:
            if self.rank[rootX] < self.rank[rootY]:
                rootX, rootY = rootY, rootX
            self.parent[rootY] = rootX
            self.rank[rootX] += self.rank[rootY]
            self.person[rootX] += self.person[rootY] + 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


if __name__ == "__main__":
    n, m = map(int, input().split())
    uf = UnionFind(m)
    edges = []
    for _ in range(n):
        # !第 i 名工人只懂得操作第 a 台和第 b 台机器，且其能力值为 c 。
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((w, u, v))
    edges.sort(key=lambda x: x[0], reverse=True)
    res = 0
    for w, u, v in edges:
        if uf.union(u, v):
            res += w
    print(res)
