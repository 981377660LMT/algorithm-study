# acwing4084. 号码牌
# 有n个小朋友，编号1～n。
# !每个小朋友都拿着一个号码牌，初始时，每个小朋友拿的号码牌上的号码都等于其编号。
# 每个小朋友都有一个幸运数字，第i个小朋友的幸运数字为d;。
# 对于第i个小朋友，他可以向第j个小朋友发起交换号码牌的请求，当且仅当|i一j|=d成立。
# 注意，请求一旦发出，对方无法拒绝，只能立刻进行交换。
# 每个小朋友都可以在任意时刻发起任意多次交换请求。给定一个1～ n的排列a1, a2, .. . , an o
# !请问，通过小朋友相互之间交换号码牌，能否使得第i个小朋友拿的号码牌上的号码恰好为ai，对i ∈[1, n]均成立。
# !n<=1e5

# 结论:
# !i和a[i]如果在同一个分量内 那么就可以交换

import sys
from collections import defaultdict
from typing import DefaultDict, List, Set

sys.setrecursionlimit(int(1e9))
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

    def getGroups(self) -> DefaultDict[int, Set[int]]:
        groups = defaultdict(set)
        for key in range(self.n):
            root = self.find(key)
            groups[root].add(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


if __name__ == "__main__":
    n = int(input())
    target = [int(num) - 1 for num in input().split()]
    lucky = list(map(int, input().split()))

    uf = UnionFindArray(n)
    for i in range(n):
        offset = lucky[i]
        left, right = i - offset, i + offset
        if left >= 0:
            uf.union(i, left)
        if right < n:
            uf.union(i, right)

    for i in range(n):
        if not uf.isConnected(i, target[i]):
            print("NO")
            exit(0)

    print("YES")
