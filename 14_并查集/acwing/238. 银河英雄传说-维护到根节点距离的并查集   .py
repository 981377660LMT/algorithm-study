class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n
        self.distToRoot = [0] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            root = self.find(self.parent[x])
            # 类似线段树的懒更新，把x到根节点的距离更新了
            self.distToRoot[x] += self.distToRoot[self.parent[x]]
            self.parent[x] = root
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """x接到y上"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        # if self.rank[rootX] > self.rank[rootY]:
        #     rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        # 注意这里
        self.distToRoot[rootX] += self.rank[rootY]
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 并查集+懒更新
# N≤30000,T≤500000


def main():
    uf = UnionFindArray(30000 + 10)

    T = int(input())
    for _ in range(T):
        line = list(input().split())
        op = line[0]
        x, y = int(line[1]), int(line[2])

        if op == 'M':
            uf.union(x, y)
        else:
            if not uf.isConnected(x, y):
                print(-1)
            else:
                cur = abs(uf.distToRoot[x] - uf.distToRoot[y])
                print(max(0, cur - 1))


if __name__ == "__main__":
    main()

