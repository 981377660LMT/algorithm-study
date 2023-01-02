# Dynamic Graph Vertex Add Component Sum
# 连接断开边/单点修改权值/查询连通分量和

# 0 u v 连接u v (保证u v不连接)
# 1 u v 断开u v  (保证u v连接)
# 2 u x 将u的值加上x
# 3 u 输出u所在连通块的值

# 离线查询
# n<=3e5

from typing import Tuple


class ATCRevocableUnionFindArray:
    """维护分量之和的可撤销并查集"""

    __slots__ = ("n", "parentSize", "sum", "history")

    def __init__(self, n: int):
        self.n = n
        self.parentSize = [-1] * n
        self.sum = [0] * n
        self.history = []

    def addSum(self, i: int, delta: int):
        """第i个元素的值加上v"""
        x = i
        while x >= 0:
            self.sum[x] += delta
            x = self.parentSize[x]

    def union(self, a: int, b: int) -> bool:
        x = self.find(a)
        y = self.find(b)
        if -self.parentSize[x] < -self.parentSize[y]:
            x, y = y, x
        self.history.append((x, self.parentSize[x]))
        self.history.append((y, self.parentSize[y]))
        if x == y:
            return False
        self.parentSize[x] += self.parentSize[y]
        self.parentSize[y] = x
        self.sum[x] += self.sum[y]
        return True

    def find(self, a: int) -> int:
        x = a
        while self.parentSize[x] >= 0:
            x = self.parentSize[x]
        return x

    def isConnected(self, a: int, b: int) -> bool:
        return self.find(a) == self.find(b)

    def revocate(self) -> bool:
        if not self.history:
            return False
        y, py = self.history.pop()
        x, px = self.history.pop()
        if self.parentSize[x] != px:
            self.sum[x] -= self.sum[y]
        self.parentSize[x] = px
        self.parentSize[y] = py
        return True

    def getComponentSum(self, i) -> int:
        return self.sum[self.find(i)]


class OfflineDynamicConnectivity:
    def __init__(self, n: int, q: int):
        """离线动态连通性查询

        Args:
            n (int): 顶点数
            q (int): 查询数
        """
        self.n = n
        self.q = q
        self.log = (q - 1).bit_length()
        self.size = 1 << self.log
        self.seg = [[] for _ in range(2 * self.size)]
        self.edges = []
        self.edgeId = {}
        self.remain = set()
        self.uf = ATCRevocableUnionFindArray(self.n)

    def addEdge(self, u: int, v: int, t: int):
        """時刻tに辺u-vを追加する。"""
        if u > v:
            u, v = v, u
        self.remain.add((u, v))
        self.edgeId[(u, v)] = t

    def removeEdge(self, u: int, v: int, t: int):
        """時刻tに辺u-vを削除する。"""
        if u > v:
            u, v = v, u
        self.remain.discard((u, v))
        self.edges.append((self.edgeId[(u, v)], t, (u, v)))

    def build(self):
        """SegmentTree上にinsert/eraseクエリを構築する。"""
        for e in self.remain:
            self.edges.append((self.edgeId[e], self.q, e))
        for l, r, e in self.edges:
            self._add(l, r, e)

    def run(self, func):
        """クエリを実行する。

        Args:
            func (function): 引数は時刻(kth query, 0-indexed)
        """
        stack = [1]
        while stack:
            k = stack.pop()
            if k >= 0:
                if self.size + self.q <= k:
                    continue
                stack.append(~k)
                for u, v in self.seg[k]:  # 行きがけ順
                    self.uf.union(u, v)
                if self.size <= k:
                    func(k - self.size)
                else:
                    stack.append(2 * k + 1)
                    stack.append(2 * k)
            else:
                for _ in self.seg[~k]:  # 帰りがけ順
                    self.uf.revocate()

    def _add(self, l: int, r: int, e: Tuple[int, int]):
        """SegmentTreeの[l, r)区間に辺u-vを追加する。"""
        l += self.size
        r += self.size
        while l < r:
            if l & 1:
                self.seg[l].append(e)
                l += 1
            if r & 1:
                r -= 1
                self.seg[r].append(e)
            l >>= 1
            r >>= 1


n, q = map(int, input().split())
res = []
query = []
dc = OfflineDynamicConnectivity(n, q)

for i, a in enumerate(map(int, input().split())):
    dc.uf.addSum(i, a)

for i in range(q):
    cur = tuple(map(int, input().split()))
    query.append(cur)
    if cur[0] == 0:
        _, u, v = cur
        dc.addEdge(u, v, i)
    elif cur[0] == 1:
        _, u, v = cur
        dc.removeEdge(u, v, i)


def func(k: int):
    cur = query[k]
    if cur[0] == 2:
        _, v, x = cur
        dc.uf.addSum(v, x)
    if cur[0] == 3:
        _, v = cur
        res.append(dc.uf.getComponentSum(v))


dc.build()
dc.run(func)

print("\n".join(map(str, res)))
