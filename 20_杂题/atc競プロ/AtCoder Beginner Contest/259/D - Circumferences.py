# 给定二维平面上的n 个圆，以及某个圆上的起点和某个圆上的终点。
# 只能走圆的边界(可以通过两圆交点更换所在的圆)，问能否从起点走到终点?
# n<=3000

# 注意这题与引爆炸弹的区别:
# !这里圆之间的关系是无向的,引爆炸弹的圆之间的关系是有向的

# 并查集判连通性，数据范围只需要n2暴力枚举判断是否有交点即可。
# 判断两圆有交点:不相离也不包含。
# 相离:圆心距大于两圆半径之和。
# 包含:圆心距小于两圆半径之差的绝对值。

from itertools import combinations
import sys
import os
from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, Tuple, TypeVar

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

Circle = Tuple[int, int, int]


def 外离(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) > (r1 + r2) * (r1 + r2)


def 外切(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) == (r1 + r2) * (r1 + r2)


def 相交(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) < (r1 + r2) * (r1 + r2)


def 内切(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) == (r1 - r2) * (r1 - r2)


def 内含(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) < (r1 - r2) * (r1 - r2)


def main() -> None:
    n = int(input())
    sx, sy, tx, ty = map(int, input().split())
    circles = [tuple(map(int, input().split())) for _ in range(n)]

    uf = UnionFindMap[int]()
    for i, j in combinations(range(n), 2):
        if not 外离(circles[i], circles[j]) and not 内含(circles[i], circles[j]):
            uf.union(i, j)

    g1, g2 = -1, -2
    for i in range(n):
        cx, cy, cr = circles[i]
        if (sx - cx) * (sx - cx) + (sy - cy) * (sy - cy) == cr * cr:
            g1 = i
        if (tx - cx) * (tx - cx) + (ty - cy) * (ty - cy) == cr * cr:
            g2 = i

    if uf.isConnected(g1, g2):
        print("Yes")
    else:
        print("No")


if __name__ == "__main__":
    T = TypeVar("T", bound=Hashable)

    class UnionFindMap(Generic[T]):
        """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

        def __init__(self, iterable: Optional[Iterable[T]] = None):
            self.part = 0
            self.parent = dict()
            self.rank = defaultdict(lambda: 1)
            for item in iterable or []:
                self._add(item)

        def union(self, key1: T, key2: T) -> bool:
            """rank一样时 默认key2作为key1的父节点"""
            root1 = self.find(key1)
            root2 = self.find(key2)
            if root1 == root2:
                return False
            if self.rank[root1] > self.rank[root2]:
                root1, root2 = root2, root1
            self.parent[root1] = root2
            self.rank[root2] += self.rank[root1]
            self.part -= 1
            return True

        def find(self, key: T) -> T:
            if key not in self.parent:
                self._add(key)
                return key

            while self.parent.get(key, key) != key:
                self.parent[key] = self.parent[self.parent[key]]
                key = self.parent[key]
            return key

        def isConnected(self, key1: T, key2: T) -> bool:
            return self.find(key1) == self.find(key2)

        def getRoots(self) -> List[T]:
            return list(set(self.find(key) for key in self.parent))

        def getGroups(self) -> DefaultDict[T, List[T]]:
            groups = defaultdict(list)
            for key in self.parent:
                root = self.find(key)
                groups[root].append(key)
            return groups

        def _add(self, key: T) -> bool:
            if key in self.parent:
                return False
            self.parent[key] = key
            self.rank[key] = 1
            self.part += 1
            return True

        def __str__(self) -> str:
            return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

        def __len__(self) -> int:
            return self.part

        def __contains__(self, key: T) -> bool:
            return key in self.parent

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
