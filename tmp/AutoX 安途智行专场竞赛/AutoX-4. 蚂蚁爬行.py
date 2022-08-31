# 蚂蚁爬行
# 在爬行过程中，对于任意两个线条，只要有接触（公共点），小蚂蚁就能从一个爬到另一个。
# 请判断这些小蚂蚁能否到达各自的目的地。

from typing import List
from typing import Tuple


Segment = Tuple[int, int, int, int]
Circle = Tuple[int, int, int]


def cross(x1: int, y1: int, x2: int, y2: int) -> int:
    """内积"""
    return x1 * y2 - y1 * x2


def isSegCross(segment1: Segment, segment2: Segment) -> bool:
    """线段 (x1,y1,x2,y2) 与 (x3,y3,x4,y4) 是否相交"""
    x1, y1, x2, y2 = segment1
    x3, y3, x4, y4 = segment2
    res1 = cross(x2 - x1, y2 - y1, x3 - x1, y3 - y1)  # 2 1 3
    res2 = cross(x2 - x1, y2 - y1, x4 - x1, y4 - y1)  # 2 1 4
    res3 = cross(x4 - x3, y4 - y3, x1 - x3, y1 - y3)  # 4 3 1
    res4 = cross(x4 - x3, y4 - y3, x2 - x3, y2 - y3)  # 4 3 2

    # 线段共线
    if res1 == 0 and res2 == 0 and res3 == 0 and res4 == 0:
        A, B, C, D = (x1, y1), (x2, y2), (x3, y3), (x4, y4)
        A, B = sorted((A, B))
        C, D = sorted((C, D))
        return max(A, C) <= min(B, D)

    # 不共线
    canAB = (res1 >= 0 and res2 <= 0) or (res1 <= 0 and res2 >= 0)  # 線分 AB が点 C, D を分けるか？
    canCD = (res3 >= 0 and res4 <= 0) or (res3 <= 0 and res4 >= 0)  # 線分 CD が点 A, B を分けるか？
    return canAB and canCD


##########################################################################


def isSegCircleCross(segment: Segment, circle: Circle) -> bool:
    """
    线段与圆是否相交

    https://blog.csdn.net/SongBai1997/article/details/86599879
    """
    sx1, sy1, sx2, sy2 = segment
    cx, cy, r = circle
    flag1 = (sx1 - cx) * (sx1 - cx) + (sy1 - cy) * (sy1 - cy) <= r * r
    flag2 = (sx2 - cx) * (sx2 - cx) + (sy2 - cy) * (sy2 - cy) <= r * r
    if flag1 and flag2:  # 两点都在圆内 不相交
        return False
    if flag1 or flag2:  # 一点在圆内一点在圆外 相交
        return True

    # 两点都在圆外

    # 将直线p1p2化为直线方程的一般式:Ax+By+C=0的形式。先化为两点式，然后由两点式得出一般式
    A = sy1 - sy2
    B = sx2 - sx1
    C = sx1 * sy2 - sx2 * sy1
    # 使用距离公式判断圆心到直线ax+by+c=0的距离是否大于半径
    dist1 = A * cx + B * cy + C
    dist1 *= dist1
    dist2 = (A * A + B * B) * r * r
    if dist1 > dist2:  # 圆心到直线距离大于半径,不相交
        return False

    # 需要判断角op1p2和角op2p1是否都为锐角,都为锐角则相交,否则不相交
    angle1 = (cx - sx1) * (sx2 - sx1) + (cy - sy1) * (sy2 - sy1)
    angle2 = (cx - sx2) * (sx1 - sx2) + (cy - sy2) * (sy1 - sy2)
    if angle1 > 0 and angle2 > 0:
        return True
    return False


##########################################################################
# !两圆之间的关系


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
    return (
        (r1 - r2) * (r1 - r2)
        < (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)
        < (r1 + r2) * (r1 + r2)
    )


def 内切(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) == (r1 - r2) * (r1 - r2)


def 内含(circle1: Circle, circle2: Circle) -> bool:
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) < (r1 - r2) * (r1 - r2)


def isCircleCross1(circle1: Circle, circle2: Circle) -> bool:
    """圆与圆是否有交点"""
    x1, y1, r1 = circle1
    x2, y2, r2 = circle2
    a = (r1 - r2) * (r1 - r2)
    b = (r1 + r2) * (r1 + r2)
    dist = (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)
    return a <= dist <= b


# def isCircleCross2(circle1: Circle, circle2: Circle) -> bool:
#     """圆与圆是否有交点"""
#     return 外切(circle1, circle2) or 相交(circle1, circle2) or 内切(circle1, circle2)


#############################################################################################
from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


class Solution:
    def antPass(self, geometry: List[List[int]], path: List[List[int]]) -> List[bool]:
        n = len(geometry)
        shapes = []
        for g in geometry:
            if len(g) == 4:
                shapes.append((g[0], g[1], g[2], g[3]))
            elif len(g) == 3:
                shapes.append((g[0], g[1], g[2]))

        uf = UnionFindArray(n)
        for i in range(n):
            for j in range(i + 1, n):
                s1, s2 = shapes[i], shapes[j]
                len1, len2 = len(s1), len(s2)
                if len1 == 4 and len2 == 4:
                    if isSegCross(s1, s2):
                        uf.union(i, j)
                elif len1 == 4 and len2 == 3:
                    if isSegCircleCross(s1, s2):
                        uf.union(i, j)
                elif len1 == 3 and len2 == 4:
                    if isSegCircleCross(s2, s1):
                        uf.union(i, j)
                elif len1 == 3 and len2 == 3:
                    if isCircleCross1(s1, s2):
                        uf.union(i, j)

        return [uf.isConnected(u, v) for u, v in path]


if __name__ == "__main__":
    print(
        Solution().antPass(
            geometry=[[2, 5, 7, 3], [1, 1, 4, 2], [4, 3, 2]], path=[[0, 1], [1, 2], [0, 2]]
        )
    )
