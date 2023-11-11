# 两点间的距离为曼哈顿距离
# !二分 + 引爆所有的炸弹

from collections import defaultdict, deque
from itertools import combinations, product
from math import ceil
import sys
from typing import DefaultDict, Set, Tuple

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

Point = Tuple[int, int, int]


def main1() -> None:
    def bfs(start: Point, adjMap: DefaultDict[Point, Set[Point]]) -> bool:
        count = 0
        queue, visited = deque([start]), set([start])
        while queue:
            cur = queue.popleft()
            count += 1
            for next in adjMap[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)

        return count == n

    def check(mid: int) -> bool:
        adjMap = defaultdict(set)
        for u, v in product(points, repeat=2):
            if dist[u][v] <= mid * u[2]:
                adjMap[u].add(v)
        return any(bfs(start, adjMap) for start in points)

    n = int(input())
    dist = defaultdict(lambda: defaultdict(lambda: int(1e18)))
    points = []
    for _ in range(n):
        x, y, p = map(int, input().split())
        points.append((x, y, p))
    for u, v in combinations(points, 2):
        cur = abs(u[0] - v[0]) + abs(u[1] - v[1])
        dist[u][v] = dist[v][u] = cur

    left, right = 1, int(1e10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1

    print(left)


def main2() -> None:
    """floyd 每个点到其他点里距离的最小值

    字典的key不要存point元组 容易TLE
    """

    n = int(input())
    points: list[Point] = []
    for _ in range(n):
        x, y, p = map(int, input().split())
        points.append((x, y, p))

    dist = defaultdict(lambda: defaultdict(lambda: int(1e10)))
    for i, j in product(range(n), repeat=2):
        cur = abs(points[i][0] - points[j][0]) + abs(points[i][1] - points[j][1])
        dist[i][j] = ceil(cur / points[i][2])

    for k, i, j in product(range(n), repeat=3):
        dist[i][j] = min(dist[i][j], max(dist[i][k], dist[k][j]))  # 注意这里距离的定义

    res = int(1e10)
    for i in range(n):
        res = min(res, max(dist[i].values()))
    print(res)


while True:
    try:
        main2()
    except (EOFError, ValueError):
        break
