# Dynamic Tree Vertex Add Path Sum
# !添加边删除边, 查询路径和, 修改点权
# 0 root1 root2 root3 root4 => 删除边root1-root2, 添加边root3-root4 (每次保证root1-root2在树中)
# 1 root add => 将root的权值加上add
# 2 root1 root2 => 查询root1到root2的路径和
# n,q<=2e5

# https://oi.men.ci/link-cut-tree-notes/


from typing import List


class LinkCutTree:
    def __init__(self, values: List[int]):
        n = len(values)
        self.ptr = [-1] * (n << 2)  # l, r, p, rev
        for i in range(n):
            self.ptr[i << 2 | 3] = 0
        self.val = values[:] + [0]
        self.sum = values[:] + [0]

    def add_edge(self, u: int, v: int):
        self._evert(u)
        self._link(u, v)

    def remove_edge(self, u: int, v: int):
        self._evert(u)
        self._cut(v)

    def add(self, u: int, delta: int):
        self._access(u)
        self.val[u] += delta

    def path_sum(self, u: int, v: int) -> int:
        self._evert(u)
        self._access(v)
        return self.sum[v]

    def _evert(self, u):
        """将某个节点 u 置为其所在树的根节点，该操作等价于把该节点到根节点所经过的所有边方向取反；"""
        self._access(u)
        self._toggle(u)
        self._push(u)

    def _link(self, u, v):
        """将某两个节点 u 和 v 连接，执行操作后 u 成为 v 的父节点；"""
        self._access(u)
        self._access(v)
        self.ptr[u << 2 | 2] = v
        self.ptr[u << 2 | 1] = u
        self._update(v)

    def _cut(self, u):
        self._access(u)
        self.ptr[self.ptr[u << 2 | 0] << 2 | 2] = -1
        self.ptr[u << 2 | 0] = -1
        self._update(u)

    def _toggle(self, u):
        if u != -1:
            self.ptr[u << 2 | 0], self.ptr[u << 2 | 1] = self.ptr[u << 2 | 1], self.ptr[u << 2 | 0]
            self.ptr[u << 2 | 3] ^= 1

    def _push(self, u):
        if self.ptr[u << 2 | 3]:
            self._toggle(self.ptr[u << 2 | 0])
            self._toggle(self.ptr[u << 2 | 1])
            self.ptr[u << 2 | 3] = 0

    def _update(self, u):
        self.sum[u] = self.val[u] + self.sum[self.ptr[u << 2 | 0]] + self.sum[self.ptr[u << 2 | 1]]

    def _state(self, u):
        par = self.ptr[u << 2 | 2]
        if par == -1 or (self.ptr[par << 2 | 0] != u and self.ptr[par << 2 | 1] != u):
            return 0
        elif self.ptr[par << 2 | 0] == u:
            return 1
        return -1

    def _rotate(self, u):
        ptr = self.ptr
        par = ptr[u << 2 | 2]
        parpar = ptr[par << 2 | 2]
        if ptr[par << 2 | 0] == u:
            c = ptr[u << 2 | 1]
            ptr[u << 2 | 1] = par
            ptr[par << 2 | 0] = c
        else:
            c = ptr[u << 2 | 0]
            ptr[u << 2 | 0] = par
            ptr[par << 2 | 1] = c
        if parpar != -1:
            if ptr[parpar << 2 | 0] == par:
                ptr[parpar << 2 | 0] = u
            if ptr[parpar << 2 | 1] == par:
                ptr[parpar << 2 | 1] = u
        ptr[u << 2 | 2] = parpar
        ptr[par << 2 | 2] = u
        if c != -1:
            ptr[c << 2 | 2] = par
        self._update(par)
        self._update(u)
        return u

    def _splay(self, u):
        self._push(u)
        while self._state(u) != 0:
            par = self.ptr[u << 2 | 2]
            if self._state(par) == 0:
                self._push(par)
                self._push(u)
                self._rotate(u)
            elif self._state(u) == self._state(par):
                self._push(self.ptr[par << 2 | 2])
                self._push(par)
                self._push(u)
                self._rotate(par)
                self._rotate(u)
            else:
                self._push(self.ptr[par << 2 | 2])
                self._push(par)
                self._push(u)
                self._rotate(u)
                self._rotate(u)

    def _access(self, u):
        """访问”某个节点 u，被“访问”过的节点会与根节点之间以路径相连，并且该节点为路径头部（最下端）；"""
        cur = u
        r_cur = -1
        while cur != -1:
            self._splay(cur)
            self.ptr[cur << 2 | 1] = r_cur
            self._update(cur)
            r_cur = cur
            cur = self.ptr[cur << 2 | 2]
        self._splay(u)


import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.buffer.readline


if __name__ == "__main__":
    n, q = map(int, input().split())
    values = list(map(int, input().split()))
    edges = [list(map(int, input().split())) for _ in range(n - 1)]
    queries = [list(map(int, input().split())) for _ in range(q)]

    lct = LinkCutTree(values)
    for u, v in edges:
        lct.add_edge(u, v)

    res = []
    for op, *args in queries:
        if op == 0:
            root1, root2, root3, root4 = args
            lct.remove_edge(root1, root2)
            lct.add_edge(root3, root4)
        elif op == 1:
            root, x = args
            lct.add(root, x)
        elif op == 2:
            root1, root2 = args
            res.append(lct.path_sum(root1, root2))

    print("\n".join(map(str, res)))
