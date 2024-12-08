class UnionFind:
    def __init__(self, n):
        self.parent = list(range(n))
        self.rank = [0] * n
        self.A_count = [0] * n
        self.B_count = [0] * n

    def find(self, x):
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x, y, w, total_cost):
        x = self.find(x)
        y = self.find(y)
        if x == y:
            return total_cost
        if self.rank[x] < self.rank[y]:
            x, y = y, x
        before_pairs = min(self.A_count[x], self.B_count[x]) + min(self.A_count[y], self.B_count[y])
        self.parent[y] = x
        if self.rank[x] == self.rank[y]:
            self.rank[x] += 1
        self.A_count[x] += self.A_count[y]
        self.B_count[x] += self.B_count[y]
        after_pairs = min(self.A_count[x], self.B_count[x])
        new_pairs = after_pairs - before_pairs
        if new_pairs > 0:
            total_cost += w * new_pairs
        return total_cost


N, M, K = map(int, input().split())
edges = []
for _ in range(M):
    u, v, w = map(int, input().split())
    u -= 1
    v -= 1
    edges.append((w, u, v))
edges.sort(key=lambda x: x[0])

A_list = list(map(int, input().split()))
B_list = list(map(int, input().split()))
for i in range(K):
    A_list[i] -= 1
    B_list[i] -= 1

uf = UnionFind(N)
countA = [0] * N
countB = [0] * N
for a in A_list:
    countA[a] += 1
for b in B_list:
    countB[b] += 1

uf.A_count = countA[:]
uf.B_count = countB[:]

total_cost = 0
for w, u, v in edges:
    total_cost = uf.union(u, v, w, total_cost)

print(total_cost)
