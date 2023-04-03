# 哈希线性探测
# 线性探测 linear probing
# 初始时数组长度为2**20 全为-1
# 有set get两种操作
# 每次set操作都会从xi开始向右找到第一个-1的位置，然后将其置为xi
# 每次get操作都会返回xi索引处的值

# !区间并查集维护每个连通分量的端点 加速查找过程
# !带合并方向的并查集


class UnionFindWithDirected:
    """带有合并方向的并查集(向右合并)"""

    __slots__ = "part", "_n", "_parent", "_rank", "_direction"

    def __init__(self, n: int, direction: int):
        """direction: 合并方向, 1: 向右合并, -1: 向左合并"""
        assert direction in (1, -1), "direction must be 1 or -1"
        self.part = n
        self._n = n
        self._parent = list(range(n))
        self._rank = [1] * n
        self._direction = direction

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """union后x所在的root的parent指向y所在的root"""
        if x < y and self._direction == -1:
            x, y = y, x
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionRange(self, left: int, right: int) -> int:
        """合并[left,right]区间, 返回合并次数."""
        if left > right:
            left, right = right, left
        leftRoot = self.find(left)
        rightRoot = self.find(right)
        unionCount = 0
        if self._direction == 1:
            while leftRoot != rightRoot:
                unionCount += 1
                self.union(leftRoot, leftRoot + 1)
                leftRoot = self.find(leftRoot + 1)
        else:
            while leftRoot != rightRoot:
                unionCount += 1
                self.union(rightRoot, rightRoot - 1)
                rightRoot = self.find(rightRoot - 1)
        return unionCount

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


if __name__ == "__main__":
    n = 2**20
    uf = UnionFindWithDirected(n, direction=1)
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
