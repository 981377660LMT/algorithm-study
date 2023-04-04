# 哈希线性探测
# 线性探测 linear probing
# 初始时数组长度为2**20 全为-1
# 有set get两种操作
# 每次set操作都会从xi开始向右找到第一个-1的位置，然后将其置为xi
# 每次get操作都会返回xi索引处的值

# !区间并查集维护每个连通分量的端点 加速查找过程
# !带合并方向的并查集


class UnionFindArray:
    def __init__(self, n: int):
        self.part = n
        self._n = n
        self._parent = list(range(n))

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """union后x所在的root的parent指向y所在的root"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        self._parent[rootX] = rootY
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


if __name__ == "__main__":
    n = 2**20
    uf = UnionFindArray(n)
    nums = [-1] * n

    q = int(input())
    for _ in range(q):
        t, x = map(int, input().split())
        if t == 1:
            cand = uf.find(x % n)  # !从xi开始向右找到第一个-1的位置
            nums[cand] = x
            uf.union(cand, (cand + 1) % n)  # !向右连接
        else:
            print(nums[x % n])
