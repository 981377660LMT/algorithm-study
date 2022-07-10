"""
无向带权连通图
1-k为休息点(基地)
q个查询 问能否从xi走到yi
高桥君不休息一次最多能行走距离为ti (ti单调递增;其实不递增也可以,并查集+离线查询按ti排序即可)
n,m,q<=2e5

解法:
1. 多源dijkstra算出点离他最近的基地 形成若干个组(每个点属于哪个基地)
那么从一个组走到另一个组后 最好先马上去休息 (最好先去离这个点最近的基地)
问题就转化为基地之间走了

2. 建立新图
3. 边排序；离线查询+并查集 判断两个组是否连通
"""

from collections import defaultdict
from heapq import heappop, heappush
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(1e18)


def main() -> None:
    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    n, m, k = map(int, input().split())
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u][v] = w
        adjMap[v][u] = w

    # !1.多源dijkstra
    dist = defaultdict(lambda: INF)
    group = defaultdict(lambda: -1)  # !点属于哪个基地组
    pq = []  # (距离, 当前点, 起始基地)
    for i in range(k):
        heappush(pq, (0, i, i))
        dist[i] = 0
        group[i] = i

    while pq:
        curDist, cur, curBase = heappop(pq)
        if dist[cur] < curDist:
            continue
        group[cur] = curBase
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next, curBase))

    # !2.建新图
    newAdjMap = defaultdict(lambda: defaultdict(lambda: INF))
    edges = set()
    for cur in range(n):
        for next in adjMap[cur]:
            u, v = group[cur], group[next]
            if u != v:
                w = dist[cur] + dist[next] + adjMap[cur][next]
                newAdjMap[u][v] = min(newAdjMap[u][v], w)
                newAdjMap[v][u] = min(newAdjMap[v][u], w)
                edges.add((*sorted((u, v)), w))

    edges = sorted(edges, key=lambda x: x[2])

    # !3.离线查询+并查集
    uf, ei = UnionFindMap(), 0
    q = int(input())
    for _ in range(q):
        x, y, t = map(int, input().split())
        x, y = x - 1, y - 1
        while ei < len(edges) and edges[ei][2] <= t:
            u, v, w = edges[ei]
            uf.union(u, v)
            ei += 1
        if uf.isConnected(x, y):
            print("Yes")
        else:
            print("No")


if __name__ == "__main__":
    from collections import defaultdict
    from typing import (
        DefaultDict,
        Generic,
        Iterable,
        List,
        Optional,
        TypeVar,
    )

    T = TypeVar("T")

    class UnionFindMap(Generic[T]):
        def __init__(self, iterable: Optional[Iterable[T]] = None):
            self.count = 0
            self.parent = dict()
            self.rank = defaultdict(lambda: 1)
            for item in iterable or []:
                self.add(item)

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
            self.count -= 1
            return True

        def find(self, key: T) -> T:
            if key not in self.parent:
                self.add(key)
                return key

            while self.parent.get(key, key) != key:
                self.parent[key] = self.parent[self.parent[key]]
                key = self.parent[key]
            return key

        def isConnected(self, key1: T, key2: T) -> bool:
            return self.find(key1) == self.find(key2)

        def getRoots(self) -> List[T]:
            return list(set(self.find(key) for key in self.parent))

        def getGroup(self) -> DefaultDict[T, List[T]]:
            groups = defaultdict(list)
            for key in self.parent:
                root = self.find(key)
                groups[root].append(key)
            return groups

        def add(self, key: T) -> bool:
            if key in self.parent:
                return False
            self.parent[key] = key
            self.rank[key] = 1
            self.count += 1
            return True

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
