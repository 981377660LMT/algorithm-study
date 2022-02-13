class UnionFind:
    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n
        self.partSum = [0] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.size[rootX] > self.size[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.size[rootY] += self.size[rootX]
        self.partSum[rootY] += self.partSum[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getPartSum(self, x: int) -> int:
        root = self.find(x)
        return self.partSum[root]


# 反向并查集：删除元素变成添加元素
n = int(input())
# 读取数组
nums = [int(v) for v in input().split()]
queries = [int(v) - 1 for v in input().split()]


uf = UnionFind(n)
visited = [False] * n
res = [0] * n


for i in range(n - 1, 0, -1):
    curIndex = queries[i]
    visited[curIndex] = True
    uf.partSum[curIndex] = nums[curIndex]

    # 合并左右邻居
    if curIndex - 1 >= 0 and visited[curIndex - 1]:
        uf.union(curIndex - 1, curIndex)
    if curIndex + 1 < n and visited[curIndex + 1]:
        uf.union(curIndex + 1, curIndex)
    # 反向更新
    res[i - 1] = max(res[i], uf.getPartSum(curIndex))


# 输出答案
for v in res:
    print(v)

