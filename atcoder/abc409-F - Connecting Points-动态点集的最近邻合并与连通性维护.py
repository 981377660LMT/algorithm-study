# F - Connecting Points (动态点集的最近邻合并与连通性维护)
# https://atcoder.jp/contests/abc409/tasks/abc409_f
#
# 有一个二维平面上的无边图，初始有 $N$ 个点，每个点有编号 $1\sim N$，坐标为 $(x_i, y_i)$。
#
# **曼哈顿距离** $d(u,v) = |x_u - x_v| + |y_u - y_v|$。
#
# 有 $Q$ 个操作，操作有三种：
#
# 1. `1 a b`：添加一个新点，坐标为 $(a, b)$，编号为当前最大编号+1。
# 2. `2`：
#    - 若所有点已连通，输出 $-1$。
#    - 否则，找到所有不同连通块之间的最小曼哈顿距离 $k$，对于所有距离为 $k$ 的点对 $(u,v)$（属于不同连通块），连一条边，把这些连通块合并，然后输出 $k$。
# 3. `3 u v`：查询点 $u$ 和 $v$ 是否在同一个连通块，是则输出 `Yes`，否则输出 `No`。
#
# N，Q <=1500
#
# “最近点合并”，类似于**层次聚类（Hierarchical Clustering）**中的最近邻合并法（single linkage），常用于数据挖掘、图像处理、地理信息聚类等场景。
# !数据的范围只有1500，所以我们可以预处理除所有点对的距离放入优先队列中，并且当新增一个点的时候，我们也将这个点所产生的所有新点对全部加入优先队列中。
# 这样对于第二个操作，我们只需要从优先对待中取出最小的点对即可。


from heapq import heappop, heappush
from typing import Tuple


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


n, q = map(int, input().split())
xs, ys = [0] * n, [0] * n
for i in range(n):
    xs[i], ys[i] = map(int, input().split())


def dist(i: int, j: int) -> int:
    return abs(xs[i] - xs[j]) + abs(ys[i] - ys[j])


pq = []


def add(i: int, j: int):
    if i > j:
        i, j = j, i
    heappush(pq, (dist(i, j), (i << 16) | j))


def get() -> Tuple[int, int, int]:
    x, ij = pq[0]
    return x, ij >> 16, ij & 0xFFFF


uf = UnionFindArraySimple(n + q)


for i in range(n):
    for j in range(i + 1, n):
        add(i, j)


def op1(a: int, b: int):
    global n
    xs.append(a)
    ys.append(b)
    n += 1
    cur = n - 1
    for i in range(n - 1):
        add(cur, i)


def op2():
    while pq:
        _, u, v = get()
        if uf.find(u) == uf.find(v):
            heappop(pq)
            continue
        break
    if not pq:
        print(-1)
        return
    d = pq[0][0]
    while pq:
        x, u, v = get()
        if d < x:
            break
        heappop(pq)
        uf.union(u, v)
    print(d)


def op3(u: int, v: int):
    if uf.find(u) == uf.find(v):
        print("Yes")
    else:
        print("No")


for _ in range(q):
    args = list(map(int, input().split()))
    if args[0] == 1:
        a, b = args[1], args[2]
        op1(a, b)
    elif args[0] == 2:
        op2()
    elif args[0] == 3:
        u, v = args[1] - 1, args[2] - 1
        op3(u, v)
