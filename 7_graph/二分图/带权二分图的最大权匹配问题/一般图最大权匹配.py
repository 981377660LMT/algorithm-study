# 任意无向图最大权匹配问题
# n<=500 m<=n*(n-1)/2
# 1<=wi<=1e6

# https://judge.yosupo.jp/submission/66790

INF = 1 << 60

from collections import deque
from typing import List, Tuple


class GeneralWeightedMatching:
    def __init__(self, n):
        self.n = n
        self.nx = n
        self.m = 2 * n + 1
        self.u = [0] * self.m * self.m
        self.v = [0] * self.m * self.m
        self.w = [0] * self.m * self.m
        self.match = [0] * self.m
        self.slack = [0] * self.m
        self.flower = [[] for _ in range(self.m)]
        self.flower_from = [0] * self.m * self.m
        self.label = [0] * self.m
        self.root = [0] * self.m
        self.par = [0] * self.m
        self.col = [0] * self.m
        self.vis = [0] * self.m
        self.que = deque()
        self.t = 0
        for u in range(1, self.m):
            for v in range(1, self.m):
                self.u[u * self.m + v] = u
                self.v[u * self.m + v] = v

    def add_edge(self, u: int, v: int, w: int):
        u += 1
        v += 1
        self.w[u * self.m + v] = max(self.w[u * self.m + v], w)
        self.w[v * self.m + u] = max(self.w[v * self.m + u], w)

    def work(self) -> Tuple[int, List[Tuple[int, int]]]:
        """返回(最大权，最大匹配)"""
        weight = 0
        count = 0
        for u in range(self.n + 1):
            self.root[u] = u
            self.flower[u].clear()
        w_max = 0
        for u in range(1, self.n + 1):
            for v in range(1, self.n + 1):
                self.flower_from[u * self.m + v] = u if u == v else 0
                w_max = max(w_max, self.w[u * self.m + v])
        for u in range(1, self.n + 1):
            self.label[u] = w_max
        while self._matching():
            count += 1
        for u in range(1, self.n + 1):
            if self.match[u] and self.match[u] < u:
                weight += self.w[u * self.m + self.match[u]]
        for i in range(self.n):
            self.match[i] = self.match[i + 1] - 1

        matching = [(i, self.match[i]) for i in range(self.n) if self.match[i] > i]
        return weight, matching

    def _dist(self, u, v):
        u, v = self.u[u * self.m + v], self.v[u * self.m + v]
        return self.label[u] + self.label[v] - self.w[u * self.m + v] * 2

    def _update_slack(self, u, x):
        if not self.slack[x] or self._dist(u, x) < self._dist(self.slack[x], x):
            self.slack[x] = u

    def _set_slack(self, x):
        self.slack[x] = 0
        for u in range(1, self.n + 1):
            if self.w[u * self.m + x] > 0 and self.root[u] != x and self.col[self.root[u]] == 0:
                self._update_slack(u, x)

    def _que_push(self, x):
        stack = [x]
        while stack:
            x = stack.pop()
            if x <= self.n:
                self.que.append(x)
                continue
            for _, fi in enumerate(self.flower[x]):
                stack.append(fi)

    def _set_root(self, x, b):
        stack = [x]
        while stack:
            x = stack.pop()
            self.root[x] = b
            if x <= self.n:
                continue
            for _, fi in enumerate(self.flower[x]):
                stack.append(fi)

    def _get_pr(self, b, xr):
        f = self.flower[b]
        pr = f.index(xr)
        if pr % 2 == 1:
            f = self.flower[b] = f[0:1] + f[1:][::-1]
            return len(f) - pr
        else:
            return pr

    def _set_match(self, u, v):
        self.match[u] = self.v[u * self.m + v]
        if u <= self.n:
            return
        xr = self.flower_from[u * self.m + self.u[u * self.m + v]]
        pr = self._get_pr(u, xr)
        f = self.flower[u]
        for i in range(pr):
            self._set_match(f[i], f[i ^ 1])
        self._set_match(xr, v)
        self.flower[u] = f[pr:] + f[:pr]

    def _augment(self, u, v):
        xnv = self.root[self.match[u]]
        self._set_match(u, v)
        while xnv:
            self._set_match(xnv, self.root[self.par[xnv]])
            u, v = self.root[self.par[xnv]], xnv
            xnv = self.root[self.match[u]]
            self._set_match(u, v)

    def _get_lca(self, u, v):
        self.t += 1
        while u or v:
            if not u:
                u, v = v, u
                continue
            if self.vis[u] == self.t:
                return u
            self.vis[u] = self.t
            u = self.root[self.match[u]]
            if u:
                u = self.root[self.par[u]]
            u, v = v, u
        return 0

    def _add_blossom(self, u, lca, v):
        b = self.n + 1
        while b <= self.nx and self.root[b]:
            b += 1
        if b > self.nx:
            self.nx += 1
        self.label[b] = 0
        self.col[b] = 0
        self.match[b] = self.match[lca]
        f = self.flower[b] = []
        f.append(lca)
        x = u
        while x != lca:
            f.append(x)
            y = self.root[self.match[x]]
            f.append(y)
            self._que_push(y)
            x = self.root[self.par[y]]
        f = self.flower[b] = f[0:1] + f[1:][::-1]
        x = v
        while x != lca:
            f.append(x)
            y = self.root[self.match[x]]
            f.append(y)
            self._que_push(y)
            x = self.root[self.par[y]]
        self._set_root(b, b)
        for x in range(1, self.nx + 1):
            self.w[b * self.m + x] = self.w[x * self.m + b] = 0
        for x in range(1, self.n + 1):
            self.flower_from[b * self.m + x] = 0
        for _, xs in enumerate(f):
            for x in range(1, self.nx + 1):
                if self.w[b * self.m + x] == 0 or self._dist(xs, x) < self._dist(b, x):
                    self.u[b * self.m + x] = self.u[xs * self.m + x]
                    self.u[x * self.m + b] = self.u[x * self.m + xs]
                    self.v[b * self.m + x] = self.v[xs * self.m + x]
                    self.v[x * self.m + b] = self.v[x * self.m + xs]
                    self.w[b * self.m + x] = self.w[xs * self.m + x]
                    self.w[x * self.m + b] = self.w[x * self.m + xs]
            for x in range(1, self.n + 1):
                if self.flower_from[xs * self.m + x]:
                    self.flower_from[b * self.m + x] = xs
        self._set_slack(b)

    def _expand_blossom(self, b):
        f = self.flower[b]
        for i, fi in enumerate(f):
            self._set_root(fi, fi)
        xr = self.flower_from[b * self.m + self.u[b * self.m + self.par[b]]]
        pr = self._get_pr(b, xr)
        f = self.flower[b]
        for i in range(0, pr, 2):
            xs = f[i]
            xns = f[i + 1]
            self.par[xs] = self.u[xns * self.m + xs]
            self.col[xs] = 1
            self.col[xns] = 0
            self.slack[xs] = 0
            self._set_slack(xns)
            self._que_push(xns)
        self.col[xr] = 1
        self.par[xr] = self.par[b]
        for i in range(pr + 1, len(f)):
            xs = f[i]
            self.col[xs] = -1
            self._set_slack(xs)
        self.root[b] = 0

    def _on_found_edge(self, u, v):
        eu = self.u[u * self.m + v]
        ev = self.v[u * self.m + v]
        u = self.root[eu]
        v = self.root[ev]
        if self.col[v] == -1:
            self.par[v] = eu
            self.col[v] = 1
            nu = self.root[self.match[v]]
            self.slack[v] = self.slack[nu] = 0
            self.col[nu] = 0
            self._que_push(nu)
        elif self.col[v] == 0:
            lca = self._get_lca(u, v)
            if not lca:
                self._augment(u, v)
                self._augment(v, u)
                return 1
            else:
                self._add_blossom(u, lca, v)
        return 0

    def _matching(self):
        for i in range(self.nx + 1):
            self.col[i] = -1
            self.slack[i] = 0
        self.que.clear()
        for x in range(1, self.nx + 1):
            if self.root[x] == x and not self.match[x]:
                self.par[x] = 0
                self.col[x] = 0
                self._que_push(x)
        if not self.que:
            return 0
        while True:
            while self.que:
                u = self.que.popleft()
                if self.col[self.root[u]] == 1:
                    continue
                for v in range(1, self.n + 1):
                    if self.w[u * self.m + v] and self.root[u] != self.root[v]:
                        if self._dist(u, v) == 0:
                            if self._on_found_edge(u, v):
                                return 1
                        else:
                            self._update_slack(u, self.root[v])
            d = INF
            for b in range(self.n + 1, self.nx + 1):
                if self.root[b] == b and self.col[b] == 1:
                    d = min(d, self.label[b] // 2)
            for x in range(1, self.nx + 1):
                if self.root[x] == x and self.slack[x]:
                    if self.col[x] == -1:
                        d = min(d, self._dist(self.slack[x], x))
                    elif self.col[x] == 0:
                        d = min(d, self._dist(self.slack[x], x) // 2)
            for u in range(1, self.n + 1):
                if self.col[self.root[u]] == 0:
                    if self.label[u] <= d:
                        return 0
                    self.label[u] -= d
                elif self.col[self.root[u]] == 1:
                    self.label[u] += d
            for b in range(self.n + 1, self.nx + 1):
                if self.root[b] == b:
                    if self.col[b] == 0:
                        self.label[b] += d * 2
                    elif self.col[b] == 1:
                        self.label[b] -= d * 2
            self.que.clear()
            for x in range(1, self.nx + 1):
                if (
                    self.root[x] == x
                    and self.slack[x]
                    and self.root[self.slack[x]] != x
                    and self._dist(self.slack[x], x) == 0
                ):
                    if self._on_found_edge(self.slack[x], x):
                        return 1
            for b in range(self.n + 1, self.nx + 1):
                if self.root[b] == b and self.col[b] == 1 and self.label[b] == 0:
                    self._expand_blossom(b)


n, m = map(int, input().split())
G = GeneralWeightedMatching(n)

for _ in range(m):
    u, v, w = map(int, input().split())
    G.add_edge(u, v, w)

weight, matching = G.work()
print(len(matching), weight)

for u, v in matching:
    print(u, v)
