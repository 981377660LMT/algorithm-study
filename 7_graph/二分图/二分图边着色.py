# 二分图边着色(Edge Coloring of Bipartite Graph)
# L,R<=1e5 M<=1e5
# ai<L bi<R

from heapq import heapify, heappop, heappush
from collections import deque
from collections import defaultdict
from typing import DefaultDict, List


class _HopCroftKarp:
    def __init__(self, n, m):
        self.n = n
        self.m = m
        self.G = [[] for _ in range(n)]
        self.RG = [[] for _ in range(m)]
        self.match_l = [-1] * n
        self.match_r = [-1] * m
        self.used = [0] * n
        self.time_stamp = 0

    def add_edges(self, u, v):
        self.G[u].append(v)

    def _build_argument_path(self):
        queue = deque()
        self.dist = [-1] * self.n
        for i in range(self.n):
            if self.match_l[i] == -1:
                queue.append(i)
                self.dist[i] = 0
        while queue:
            a = queue.popleft()
            for b in self.G[a]:
                c = self.match_r[b]
                if c >= 0 and self.dist[c] == -1:
                    self.dist[c] = self.dist[a] + 1
                    queue.append(c)

    def _find_min_dist_argument_path(self, a):
        self.used[a] = self.time_stamp
        for b in self.G[a]:
            c = self.match_r[b]
            if c < 0 or (
                self.used[c] != self.time_stamp
                and self.dist[c] == self.dist[a] + 1
                and self._find_min_dist_argument_path(c)
            ):
                self.match_r[b] = a
                self.match_l[a] = b
                return True
        return False

    def max_matching(self):
        while 1:
            self._build_argument_path()
            self.time_stamp += 1
            flow = 0
            for i in range(self.n):
                if self.match_l[i] == -1:
                    flow += self._find_min_dist_argument_path(i)
            if flow == 0:
                break
        ret = []
        for i in range(self.n):
            if self.match_l[i] >= 0:
                ret.append((i, self.match_l[i]))
        return ret


class _UnionFindArray:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
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


def _contract(deg, k):
    hq = []
    for i, d in enumerate(deg):
        hq.append([d, i])
    heapify(hq)
    UF = _UnionFindArray(len(deg))
    while len(hq) >= 2:
        p = heappop(hq)
        q = heappop(hq)
        if p[0] + q[0] > k:
            continue
        p[0] += q[0]
        UF.union(p[1], q[1])
        heappush(hq, p)
    return UF


def _build_k_regular_graph(n, m, A, B):
    dega = [0] * n
    degb = [0] * m
    for a in A:
        dega[a] += 1
    for b in B:
        degb[b] += 1
    K = max(*dega, *degb)

    UFa = _contract(dega, K)
    ida = [-1] * n
    pa = 0
    for i in range(n):
        if UFa.find(i) == i:
            ida[i] = pa
            pa += 1

    UFb = _contract(degb, K)
    idb = [-1] * m
    pb = 0
    for i in range(m):
        if UFb.find(i) == i:
            idb[i] = pb
            pb += 1

    p = max(pa, pb)
    dega = [0] * p
    degb = [0] * p

    C = []
    D = []
    for i in range(len(A)):
        u = ida[UFa.find(A[i])]
        v = idb[UFb.find(B[i])]
        C.append(u)
        D.append(v)
        dega[u] += 1
        degb[v] += 1
    j = 0
    for i in range(p):
        while dega[i] < K:
            while degb[j] == K:
                j += 1
            C.append(i)
            D.append(j)
            dega[i] += 1
            degb[j] += 1

    return K, p, C, D


def edgeColoring(a, b, A, B):
    K, n, A, B = _build_k_regular_graph(a, b, A, B)

    ord = [i for i in range(len(A))]
    res = []

    def euler_trail(ord):
        V = 2 * n
        G = [[] for _ in range(V)]
        m = 0
        for i in ord:
            G[A[i]].append((B[i] + n, m))
            G[B[i] + n].append((A[i], m))
            m += 1
        used_v = [0] * V
        used_e = [0] * m
        res = []
        for i in range(V):
            if used_v[i]:
                continue
            st = []
            ord2 = []
            st.append((i, -1))
            while st:
                id_ = st[-1][0]
                used_v[id_] = True
                if len(G[id_]) == 0:
                    ord2.append(st[-1][1])
                    st.pop()
                else:
                    e = G[id_][-1]
                    G[id_].pop()
                    if used_e[e[1]]:
                        continue
                    used_e[e[1]] = True
                    st.append(e)
            ord2.pop()
            ord2 = ord2[::-1]
            res += ord2
        for i, a in enumerate(res):
            res[i] = ord[a]
        return res

    def rec(ord, K):
        if K == 0:
            return
        elif K == 1:
            res.append(ord)
            return
        elif K & 1:
            G = _HopCroftKarp(n, n)
            for i in ord:
                G.add_edges(A[i], B[i])
            G.max_matching()
            lst = []
            res.append([])
            for i in ord:
                if G.match_l[A[i]] == B[i]:
                    G.match_l[A[i]] = -1
                    res[-1].append(i)
                else:
                    lst.append(i)
            rec(lst, K - 1)
        else:
            path = euler_trail(ord)
            L = []
            R = []
            for i, p in enumerate(path):
                if i & 1:
                    L.append(p)
                else:
                    R.append(p)
            rec(L, K // 2)
            rec(R, K // 2)

    rec(ord, K)
    return K, res


L, R, m = map(int, input().split())
A = [-1] * m
B = [-1] * m
for i in range(m):
    A[i], B[i] = map(int, input().split())

count, res = edgeColoring(L, R, A, B)
color = [-1] * m
for i in range(len(res)):
    for j in res[i]:
        if j < m:
            color[j] = i

print(count)  # number of colors
print(*color, sep="\n")  # color of each edge
