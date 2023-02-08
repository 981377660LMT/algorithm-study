# 有向连通图最小生成树
# Directed MST
# n,m<=2e5


from typing import List, Tuple
from collections import defaultdict
from typing import DefaultDict, List


def directedMST(n: int, edges: List[Tuple[int, int, int]], root=0) -> Tuple[int, List[int]]:
    """返回生成树的权值和和生成树的父亲节点列表 根节点父亲节点为-1"""
    m = len(edges)
    froms = [0] * n
    from_cost = [0] * n
    from_heap = [SkewHeap() for _ in range(n)]
    UF = UnionFindArray(n)
    par = [-1] * m
    stem = [-1] * n
    used = [0] * n
    used[root] = 2
    inds = []
    for i, (u, v, c) in enumerate(edges):
        from_heap[v].heappush(c * m + i)
    res = 0
    for v in range(n):
        if used[v] != 0:
            continue
        proc = []
        chi = []
        cycle = 0
        while used[v] != 2:
            used[v] = 1
            proc.append(v)
            if from_heap[v].root is None:
                return -1, [-1] * n
            tmp = from_heap[v].heappop()
            from_cost[v], ind = tmp // m, tmp % m
            froms[v] = UF.find(edges[ind][0])
            if stem[v] == -1:
                stem[v] = ind
            if froms[v] == v:
                continue
            res += from_cost[v]
            inds.append(ind)
            while cycle:
                par[chi.pop()] = ind
                cycle -= 1
            chi.append(ind)
            if used[froms[v]] == 1:
                p = v
                while True:
                    if not from_heap[p].root is None:
                        from_heap[p].heapadd(-from_cost[p] * m)
                    if p != v:
                        UF.union(v, p)
                        from_heap[v].root = from_heap[v].heapmeld(
                            from_heap[v].root, from_heap[p].root
                        )
                    p = UF.find(froms[p])
                    new_v = UF.find(v)
                    from_heap[new_v] = from_heap[v]
                    v = new_v
                    cycle += 1
                    if p == v:
                        break
            else:
                v = froms[v]
        for v in proc:
            used[v] = 2
    visited = [0] * m
    tree = [-1] * n
    for i in inds[::-1]:
        if visited[i]:
            continue
        u, v, c = edges[i]
        tree[v] = u
        x = stem[v]
        while x != i:
            visited[x] = 1
            x = par[x]
    return res, tree


class UnionFindArray:

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


class SHNode:
    __slots__ = "left", "right", "val", "add"

    def __init__(self, val):
        self.left = None
        self.right = None
        self.val = val
        self.add = 0

    def lazy(self):
        if self.left is not None:
            self.left.add += self.add
        if self.right is not None:
            self.right.add += self.add
        self.val += self.add
        self.add = 0


class SkewHeap:
    __slots__ = "root"

    def __init__(self):
        self.root = None

    def heapmeld(self, h1, h2):
        if h1 is None:
            return h2
        if h2 is None:
            return h1
        if h1.val + h1.add > h2.val + h2.add:
            h1, h2 = h2, h1
        h1.lazy()
        h1.right = self.heapmeld(h2, h1.right)
        h1.left, h1.right = h1.right, h1.left
        return h1

    def heappop(self):
        res = self.root
        res.lazy()
        self.root = self.heapmeld(res.left, res.right)
        return res.val

    def heappush(self, x):
        nh = SHNode(x)
        self.root = self.heapmeld(self.root, nh)

    def heaptop(self):
        if self.root is None:
            return None
        return self.root.val

    def heapadd(self, val):
        self.root.add += val


n, m, root = map(int, input().split())
edges = [tuple(map(int, input().split())) for _ in range(m)]
w, parents = directedMST(n, edges, root)

rootIndex = parents.index(-1)
parents[rootIndex] = rootIndex
print(w)
print(*parents)
