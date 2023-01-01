# https://judge.yosupo.jp/submission/102920
class Tree:
    def __init__(self, n):
        self.n = n
        self.tree = [[] for _ in range(n)]
        self.root = None

    def add_edge(self, u, v):
        self.tree[u].append(v)
        self.tree[v].append(u)

    def set_root(self, r=0):
        self.root = r
        self.par = [None] * self.n
        self.dep = [0] * self.n
        self.height = [0] * self.n
        self.size = [1] * self.n
        self.ord = [r]
        stack = [r]
        while stack:
            v = stack.pop()
            for adj in self.tree[v]:
                if self.par[v] == adj:
                    continue
                self.par[adj] = v
                self.dep[adj] = self.dep[v] + 1
                self.ord.append(adj)
                stack.append(adj)
        for v in self.ord[1:][::-1]:
            self.size[self.par[v]] += self.size[v]
            self.height[self.par[v]] = max(self.height[self.par[v]], self.height[v] + 1)

    def rerooting(self, op, e, merge, id):
        if self.root is None:
            self.set_root()
        dp = [e] * self.n
        lt = [id] * self.n
        rt = [id] * self.n
        inv = [id] * self.n
        for v in self.ord[::-1]:
            tl = tr = e
            for adj in self.tree[v]:
                if self.par[v] == adj:
                    continue
                lt[adj] = tl
                tl = op(tl, dp[adj])
            for adj in self.tree[v][::-1]:
                if self.par[v] == adj:
                    continue
                rt[adj] = tr
                tr = op(tr, dp[adj])
            dp[v] = tr
        for v in self.ord:
            if v == self.root:
                continue
            p = self.par[v]
            inv[v] = op(merge(lt[v], rt[v]), inv[p])
            dp[v] = op(dp[v], inv[v])
        return dp

    def euler_tour(self):
        if self.root is None:
            self.set_root()
        self.tour = []
        self.etin = [None for _ in range(self.n)]
        self.etout = [None for _ in range(self.n)]
        used = [0 for _ in range(self.n)]
        used[self.root] = 1
        stack = [self.root]
        while stack:
            v = stack.pop()
            if v >= 0:
                self.tour.append(v)
                stack.append(~v)
                if self.etin[v] is None:
                    self.etin[v] = len(self.tour) - 1
                for adj in self.tree[v]:
                    if used[adj]:
                        continue
                    used[adj] = 1
                    stack.append(adj)
            else:
                self.etout[~v] = len(self.tour)
                if ~v != self.root:
                    self.tour.append(self.par[~v])

    def heavylight_decomposition(self):
        if self.root is None:
            self.set_root()
        self.hldid = [None] * self.n
        self.hldtop = [None] * self.n
        self.hldtop[self.root] = self.root
        self.hldnxt = [None] * self.n
        self.hldrev = [None] * self.n
        stack = [self.root]
        cnt = 0
        while stack:
            v = stack.pop()
            self.hldid[v] = cnt
            self.hldrev[cnt] = v
            cnt += 1
            maxs = 0
            for adj in self.tree[v]:
                if self.par[v] == adj:
                    continue
                if maxs < self.size[adj]:
                    maxs = self.size[adj]
                    self.hldnxt[v] = adj
            for adj in self.tree[v]:
                if self.par[v] == adj or self.hldnxt[v] == adj:
                    continue
                self.hldtop[adj] = adj
                stack.append(adj)
            if self.hldnxt[v] is not None:
                self.hldtop[self.hldnxt[v]] = self.hldtop[v]
                stack.append(self.hldnxt[v])

    def lca(self, u, v):
        while True:
            if self.hldid[u] > self.hldid[v]:
                u, v = v, u
            if self.hldtop[u] != self.hldtop[v]:
                v = self.par[self.hldtop[v]]
            else:
                return u

    def dist(self, u, v):
        lca = self.lca(u, v)
        return self.dep[u] + self.dep[v] - 2 * self.dep[lca]

    def range_query(self, u, v, edge_query=False):
        while True:
            if self.hldid[u] > self.hldid[v]:
                u, v = v, u
            if self.hldtop[u] != self.hldtop[v]:
                yield self.hldid[self.hldtop[v]], self.hldid[v] + 1
                v = self.par[self.hldtop[v]]
            else:
                yield self.hldid[u] + edge_query, self.hldid[v] + 1
                return

    def subtree_query(self, u):
        return self.hldid[u], self.hldid[u] + self.size[u]

    def _get_centroid_(self, r):
        self._par_[r] = None
        self._size_[r] = 1
        ord = [r]
        stack = [r]
        while stack:
            v = stack.pop()
            for adj in self.tree[v]:
                if self._par_[v] == adj or self.cdused[adj]:
                    continue
                self._size_[adj] = 1
                self._par_[adj] = v
                ord.append(adj)
                stack.append(adj)
        if len(ord) <= 2:
            return r
        for v in ord[1:][::-1]:
            self._size_[self._par_[v]] += self._size_[v]
        sr = self._size_[r] // 2
        v = r
        while True:
            for adj in self.tree[v]:
                if self._par_[v] == adj or self.cdused[adj]:
                    continue
                if self._size_[adj] > sr:
                    v = adj
                    break
            else:
                return v

    def centroid_decomposition(self):
        self._par_ = [None] * self.n
        self._size_ = [1] * self.n
        self.cdpar = [None] * self.n
        self.cddep = [0] * self.n
        self.cdord = [None] * self.n
        self.cdused = [0] * self.n
        cnt = 0
        stack = [0]
        while stack:
            v = stack.pop()
            p = self.cdpar[v]
            c = self._get_centroid_(v)
            self.cdused[c] = True
            self.cdpar[c] = p
            self.cddep[c] = self.cddep[v]
            self.cdord[c] = cnt
            cnt += 1
            for adj in self.tree[c]:
                if self.cdused[adj]:
                    continue
                self.cdpar[adj] = c
                self.cddep[adj] = self.cddep[c] + 1
                stack.append(adj)

    def centroid(self):
        if self.root is None:
            self.set_root()
        sr = self.size[self.root] // 2
        v = self.root
        while True:
            for adj in self.tree[v]:
                if self.par[v] == adj:
                    continue
                if self.size[adj] > sr:
                    v = adj
                    break
            else:
                return v

    def diam(self):
        if self.root is None:
            self.set_root()
        u = self.dep.index(max(self.dep))
        self.set_root(u)
        v = self.dep.index(max(self.dep))
        return u, v

    def get_path(self, u, v):
        if self.root != u:
            self.set_root(u)
        path = []
        while v != None:
            path.append(v)
            v = self.par[v]
        return path

    def longest_path_decomposition(self, make_ladder=True):
        assert self.root is not None
        self.lpdnxt = [None] * self.n
        self.lpdtop = [None] * self.n
        self.lpdtop[self.root] = self.root
        stack = [self.root]
        while stack:
            v = stack.pop()
            for adj in self.tree[v]:
                if self.par[v] == adj:
                    continue
                if self.height[v] == self.height[adj] + 1:
                    self.lpdnxt[v] = adj
            for adj in self.tree[v]:
                if self.par[v] == adj or self.lpdnxt[v] == adj:
                    continue
                self.lpdtop[adj] = adj
                stack.append(adj)
            if self.lpdnxt[v] is not None:
                self.lpdtop[self.lpdnxt[v]] = self.lpdtop[v]
                stack.append(self.lpdnxt[v])
        if make_ladder:
            self._make_ladder_()

    def _make_ladder_(self):
        self.ladder = [[] for _ in range(self.n)]
        for v in range(self.n):
            if self.lpdtop[v] != v:
                continue
            to = v
            path = []
            while to is not None:
                path.append(to)
                to = self.lpdnxt[to]
            p = self.par[v]
            self.ladder[v] = path[::-1]
            for i in range(len(path)):
                self.ladder[v].append(p)
                if p is None:
                    break
                p = self.par[p]

    def level_ancestor(self, v, k):
        while v is not None:
            id = self.height[v]
            h = self.lpdtop[v]
            if len(self.ladder[h]) > k + id:
                return self.ladder[h][k + id]
            v = self.ladder[h][-1]
            k -= len(self.ladder[h]) - id - 1

    def jump(self, u, v, k):
        lca = self.lca(u, v)
        dist = self.dep[u] + self.dep[v] - 2 * self.dep[lca]
        if dist < k:
            return
        if k < self.dep[u] - self.dep[lca]:
            return self.level_ancestor(u, k)
        else:
            return self.level_ancestor(v, dist - k)

    def hash(self):
        import random

        mod = 2**61 - 1
        if self.root is None:
            self.set_root(0)
        dp = [1] * self.n
        rnd = [random.randint(0, mod - 1) for _ in range(self.n)]
        for v in self.ord[::-1]:
            h = self.height[v]
            r = rnd[h]
            for a in self.tree[v]:
                if self.par[v] == a:
                    continue
                dp[v] *= r + dp[a]
                dp[v] %= mod
        return dp


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    n = int(input())
    parents = list(map(int, input().split()))

    tree = Tree(n)
    for cur, pre in enumerate(parents):
        tree.add_edge(cur + 1, pre)
    tree.set_root(0)

    # 以每个根作为根节点求出树的哈希值
    dp = tree.hash()
    allNums = sorted(set(dp))
    mp = {v: k for k, v in enumerate(allNums)}
    res = [mp[h] for h in dp]

    print(len(allNums))  # 哈希值的种类数
    print(*res)  # 每个根的哈希值的编号
