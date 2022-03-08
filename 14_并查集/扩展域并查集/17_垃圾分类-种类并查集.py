# 最近很多城市都搞起了垃圾分类，已知有一个小区有 n 堆垃圾需要丢弃，
# 这些垃圾都打包了，我们并不知道这是什么垃圾，要知道有些垃圾如果丢在一起是会影响环境的。
# 这个小区一共只有`两辆垃圾车`，我们希望在不影响环境的情况下，每次让垃圾车载走最多的垃圾，但是因为两位司机师傅有矛盾，所以两辆车装的垃圾数目一定要相同，不然其中一位司机师傅就会不开心。
# 我们为垃圾袋标了号，分别是 1-n，有 m 个约束，每个约束表示为“a b”，意思是第a堆垃圾与第b堆垃圾不能装在一辆垃圾车上。(每堆垃圾最多有两个约束条件)
# 请问两辆垃圾车一次最多可以带走多少堆垃圾呢?


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


n, m = map(int, input().split())
uf = UnionFindArray(n * 2 + 2)
res = n
for _ in range(m):
    a, b = map(int, input().split())
    if uf.isConnected(a, b):
        res -= 1
    uf.union(a, b + n)
    uf.union(a + n, b)
print(2 * (res // 2))
