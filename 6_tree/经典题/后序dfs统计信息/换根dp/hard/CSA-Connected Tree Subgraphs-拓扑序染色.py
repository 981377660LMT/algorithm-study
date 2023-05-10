# https://csacademy.com/contest/archive/task/connected-tree-subgraphs/statement/
# 给你一棵n个结点的树，求有多少种染色方案，
# 使得染色过程中染过色的结点始终连成一块。
# !蚁群/蚂蚁 那道题的换根dp版 (拓扑序染色)

# 或者说:
# 给每个顶点贴上标签(1-n) 问有多少种贴法(一共n!种可能)
# !使得贴标签的过程中,贴过标签的顶点始终连成一个连通块
# !每个结点作为根1,换根dp

import sys
from typing import Tuple
from Rerooting import Rerooting

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)


class Enumeration:
    __slots__ = ("_fac", "_ifac", "_inv", "_mod")

    def __init__(self, size: int, mod: int) -> None:
        self._mod = mod
        self._fac = [1]
        self._ifac = [1]
        self._inv = [1]
        self._expand(size)

    def fac(self, k: int) -> int:
        self._expand(k)
        return self._fac[k]

    def ifac(self, k: int) -> int:
        self._expand(k)
        return self._ifac[k]

    def inv(self, k: int) -> int:
        """模逆元"""
        self._expand(k)
        return self._inv[k]

    def C(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        mod = self._mod
        return self.fac(n) * self.ifac(k) % mod * self.ifac(n - k) % mod

    def P(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        mod = self._mod
        return self.fac(n) * self.ifac(n - k) % mod

    def H(self, n: int, k: int) -> int:
        """可重复选取元素的组合数"""
        if n == 0:
            return 1 if k == 0 else 0
        return self.C(n + k - 1, k)

    def _expand(self, size: int) -> None:
        size = min(size, self._mod - 1)
        if len(self._fac) < size + 1:
            mod = self._mod
            preSize = len(self._fac)
            diff = size + 1 - preSize
            self._fac += [1] * diff
            self._ifac += [1] * diff
            self._inv += [1] * diff
            for i in range(preSize, size + 1):
                self._fac[i] = self._fac[i - 1] * i % mod
            self._ifac[size] = pow(self._fac[size], mod - 2, mod)  # !modInv
            for i in range(size - 1, preSize - 1, -1):
                self._ifac[i] = self._ifac[i + 1] * (i + 1) % mod
            for i in range(preSize, size + 1):
                self._inv[i] = self._ifac[i] * self._fac[i - 1] % mod


MOD = int(1e9 + 7)
COMB = Enumeration(int(1e5 + 10), MOD)
if __name__ == "__main__":
    E = Tuple[int, int]  # (子树大小,方案数)

    def e(root: int) -> E:
        return (0, 1)

    def op(childRes1: E, childRes2: E) -> E:
        size1, value1 = childRes1
        size2, value2 = childRes2
        size3 = size1 + size2
        # !顺序安排空位 => 拓扑序+组合数
        value3 = value1 * value2 % MOD * COMB.C(size1 + size2, size1) % MOD
        return (size3, value3)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        size, value = fromRes
        return size + 1, value

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    # for _, v in dp:
    #     print(v % MOD)

    res = 0
    for _, v in dp:
        res += v
        res %= MOD
    print(res)
    # print(res * COMB.inv(2) % MOD)  // !不带方向时 https://atcoder.jp/contests/tdpc/tasks/tdpc_tree
