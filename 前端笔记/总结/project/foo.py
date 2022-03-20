class UnionFind:
    def __init__(self, n):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x):
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True


def check(mid):
    uf = UnionFind(n)
    for u, v, w in edges:
        if w < mid:
            continue
        uf.union(u, v)
    return uf.count == 1


n, m = map(int, input().split())
edges = []
for _ in range(m):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    edges.append((u, v, w))

left, right = 1, int(1e4 + 10)
while left <= right:
    mid = (left + right) >> 1
    if check(mid):
        left = mid + 1
    else:
        right = mid - 1

print(right)

