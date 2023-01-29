from itertools import product


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
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

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 1≤m,n≤1000
row, col = map(int, input().split())
uf = UnionFindArray(row * col + 10)
while True:
    try:
        x1, y1, x2, y2 = map(int, input().split())
        # 注意减1
        x1, y1, x2, y2 = x1 - 1, y1 - 1, x2 - 1, y2 - 1
        p1, p2 = x1 * col + y1, x2 * col + y2
        uf.union(p1, p2)
    except EOFError:
        break

edges = set()

for r, c in product(range(row), range(col)):
    cur = r * col + c
    if r + 1 < row:
        next = (r + 1) * col + c
        edges.add((cur, next, 1))

for r, c in product(range(row), range(col)):
    cur = r * col + c
    if c + 1 < col:
        next = r * col + (c + 1)
        edges.add((cur, next, 2))

res = 0
for u, v, w in edges:
    if uf.isConnected(u, v):
        continue
    uf.union(u, v)
    res += w
print(res)
