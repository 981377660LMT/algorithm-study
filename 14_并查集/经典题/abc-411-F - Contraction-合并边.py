# abc-411-F - Contraction-合并边
# https://atcoder.jp/contests/abc411/tasks/abc411_f
#
# 给定无向简单图 G₀（N≤3×10⁵, M≤3×10⁵），顶点 i 上放置棋子 i。
# 对每次操作给出一条初始边编号 X：若棋子 Uₓ 与 Vₓ 当前位于不同顶点且这两个顶点间仍存在边，则把此边收缩（两端点合并、删自环、去重多边），并把两点上的全部棋子移到新顶点。
# 每次操作后输出当前图的边数。

# 思路
# • 把“当前顶点”视为并查集的连通块；find(v) 得到棋子 v 所在顶点编号（连通块代表）。
# • 图始终保持简单。用邻接集合 adj[r] (set[int]) 存储代表顶点 r 的邻接顶点代表。
# • 维护全局边数 E。
# • 收缩一条现存边 (a,b)：
#
# E--（删掉 a-b）
# 令 ra, rb 为代表且按邻接集大小把小的并到大的（small→big，记 big=ra）。
# 枚举 small 的所有邻居 x：
# 从 adj[x] 删掉 small；
# 若 x==big：继续；（这正是 a-b）
# 若 x 已经与 big 相连，则去重多边 E--;
# 否则在两侧集合加入连接 big-x。
# 清空 adj[small]，并查集合并。
# 整个过程中，每条边至多在“被并入的小集合”里被遍历一次，因此总复杂度 O(M log N)（均摊近乎线性）。
# • adj[v] 仅存当前代表顶点的邻居代表，保持无重边。
# • 每条边在被所在较小邻接集遍历一次，整体近线性


from typing import Callable, Optional


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(
        self, key1: int, key2: int, beforeUnion: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        if beforeUnion is not None:
            beforeUnion(root1, root2)
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def unionTo(self, parent: int, child: int) -> bool:
        """定向合并, 将child合并到parent所在的连通分量中."""
        root1, root2 = self.find(parent), self.find(child)
        if root1 == root2:
            return False
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        root = key
        while self._data[root] >= 0:
            root = self._data[root]
        while key != root:
            parent = self._data[key]
            self._data[key] = root
            key = parent
        return root

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


if __name__ == "__main__":
    N, M = map(int, input().split())
    U, V = [0] * M, [0] * M
    for i in range(M):
        U[i], V[i] = map(int, input().split())
        U[i] -= 1
        V[i] -= 1
    Q = int(input())
    X = list(map(lambda x: int(x) - 1, input().split()))

    adjSet = [set() for _ in range(N)]
    for u, v in zip(U, V):
        adjSet[u].add(v)
        adjSet[v].add(u)
    uf = UnionFindArraySimple(N)
    edgeCount = M

    def contract(a: int, b: int):
        rootA, rootB = uf.find(a), uf.find(b)
        if rootA == rootB:
            return
        if rootB not in adjSet[rootA]:
            return

        global edgeCount
        if len(adjSet[rootA]) < len(adjSet[rootB]):  # 保证rootA是大的
            rootA, rootB = rootB, rootA

        for v in list(adjSet[rootB]):
            if v == rootA:
                edgeCount -= 1  # 删除重边
            elif v in adjSet[rootA]:
                edgeCount -= 1
                adjSet[v].remove(rootB)
            else:
                adjSet[rootA].add(v)
                adjSet[v].add(rootA)
                adjSet[v].remove(rootB)

        adjSet[rootA].remove(rootB)
        adjSet[rootB].clear()
        uf.unionTo(rootA, rootB)

    for i in X:
        contract(U[i], V[i])
        print(edgeCount)
