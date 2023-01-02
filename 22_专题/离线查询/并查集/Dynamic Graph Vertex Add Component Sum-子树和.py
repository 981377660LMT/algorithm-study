# Dynamic Graph Vertex Add Component Sum
# 连接断开边/单点修改权值/查询子树和

# 0 root1 root2 root3 root4 断开(root1-root2) 连接(root3-root4)
# 1 root x 将root的值加上x
# 2 root1 root2 输出root1所在子树的和,其中root2是root1的父亲节点

# 离线查询
# n<=2e5
# !技巧:查询子树和=断开父亲边+查询连通块和+连回父亲边

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
                for u, v in self.seg[k]:
                    self.uf.union(u, v)
                if self.size <= k:
                    func(k - self.size)
                else:
                    stack.append(2 * k + 1)
                    stack.append(2 * k)
            else:
                for _ in self.seg[~k]:
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

for _ in range(n - 1):
    u, v = map(int, input().split())
    dc.addEdge(u, v, 0)  # 最开始存在的边

for i in range(q):
    cur = tuple(map(int, input().split()))
    query.append(cur)
    if cur[0] == 0:
        _, root1, root2, root3, root4 = cur
        dc.removeEdge(root1, root2, i)
        dc.addEdge(root3, root4, i)
    elif cur[0] == 2:
        _, root, parent = cur
        dc.removeEdge(root, parent, i)
        dc.addEdge(root, parent, i + 1)  # 下个时刻回溯


def func(k: int):
    cur = query[k]
    if cur[0] == 1:
        _, root, delta = cur
        dc.uf.addSum(root, delta)
    if cur[0] == 2:
        _, root, _ = cur
        res.append(dc.uf.getComponentSum(root))


dc.build()
dc.run(func)

print("\n".join(map(str, res)))
