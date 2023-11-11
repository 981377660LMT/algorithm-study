# 给定一个长度为N(N <1e3)的排列p，给出 M(M <1e5)组关系opt(i,j)，
# 你可以交换p[i], p[j]。
# !问能否通过不超过 5e5次(n*(n-1)//2) 次交换，使得排列p 变成升序的，
# !并输出方案。
# 输出形式为
# k  (操作次数)
# c1 c2 ... ck  (每次执行哪个操作,opt[ci])


# !我们将关系看过图的一条边，如果两个点在同一个联通量中，
# !说明这两个点可以通过一系列的交换关系，来相互交换。
# 如果我们暴力去跑的话，可以最大的交换次数为1e6
# !我们可以模仿拓扑排序（利用并查集），每次处理度数为1的点(把依赖少的先交换好)，
# 这样最大的操作数为
# 999 +998+...＋1 = 499500
# !注意我们必须使用并查集来建立树 否则会出现环

import sys
from collections import defaultdict, deque
from typing import DefaultDict, List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


if __name__ == "__main__":
    n = int(input())
    perm = [int(num) - 1 for num in input().split()]  # 0-n-1 的全排列

    uf = UnionFindArray(n)
    adjList = [[] for _ in range(n)]
    deg = [0] * n

    m = int(input())
    for i in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        if not uf.isConnected(u, v):  # !并查集构造生成树 注意在边里附加边的编号
            uf.union(u, v)
            adjList[u].append((v, i))
            adjList[v].append((u, i))
            deg[u] += 1
            deg[v] += 1

    def dfs(cur: int, pre: int, target: int) -> bool:
        """
        从cur出发,找到target,递归过程中在res中保存经过的边的编号

        如果找到了,返回True,否则返回False
        """
        if cur == target:
            return True
        for next, edgeId in adjList[cur]:
            if next == pre:
                continue
            if dfs(next, cur, target):
                res.append(edgeId)
                perm[cur], perm[next] = perm[next], perm[cur]
                return True
        return False

    res = []
    queue = deque([i for i in range(n) if deg[i] == 1])  # !拓扑排序
    while queue:
        cur = queue.popleft()
        target = perm.index(cur)  # cur需要换到哪个位置
        if not dfs(cur, -1, target):
            print(-1)
            exit(0)
        for next, _ in adjList[cur]:
            deg[next] -= 1
            if deg[next] == 1:
                queue.append(next)

    print(len(res))
    print(*[num + 1 for num in res])
