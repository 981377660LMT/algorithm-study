# 次短路计数
# https://www.acwing.com/problem/content/385/

from heapq import heappop, heappush
from typing import Tuple, List


INF = int(1e18)


def solve(
    n: int, directedEdges: List[Tuple[int, int, int]], start: int
) -> Tuple[List[List[int]], List[List[int]]]:
    """求最短路、次短路的距离与路径数."""
    adjList = [[] for _ in range(n)]
    for u, v, w in directedEdges:
        adjList[u].append((v, w))
    dist = [[INF] * 2 for _ in range(n)]
    count = [[0] * 2 for _ in range(n)]
    dist[start][0] = 0
    count[start][0] = 1
    pq = [(0, start, 0)]  # (dist, node, type)
    while pq:
        curDist, cur, curType = heappop(pq)
        if dist[cur][curType] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur][curType] + weight
            if dist[next][0] > cand:
                dist[next][1], count[next][1] = dist[next][0], count[next][0]
                heappush(pq, (dist[next][1], next, 1))
                dist[next][0], count[next][0] = cand, count[cur][curType]
                heappush(pq, (dist[next][0], next, 0))
            elif dist[next][0] == cand:
                count[next][0] += count[cur][curType]
            elif dist[next][1] > cand:
                dist[next][1], count[next][1] = cand, count[cur][curType]
                heappush(pq, (dist[next][1], next, 1))
            elif dist[next][1] == cand:
                count[next][1] += count[cur][curType]
    return dist, count


# G. Counting Shortcuts
# 给定一张无向图，求出与 (s,t) 之间的与最短路长度差不超 1 的路径条数模 1e9+7。
# https://codeforces.com/contest/1650/problem/G
def cf776(n: int, edges: List[Tuple[int, int, int]], start: int, end: int) -> int:
    MOD = int(1e9 + 7)
    dist, count = solve(n, edges, start)
    count1, count2 = count[end][0], count[end][1]
    return (count1 + count2) % MOD if dist[end][0] + 1 == dist[end][1] else count1


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    def solve1():
        T = int(input())
        for _ in range(T):
            n, m = map(int, input().split())
            edges = []
            for _ in range(m):
                u, v, w = map(int, input().split())
                edges.append((u - 1, v - 1, w))
            start, end = map(int, input().split())
            start, end = start - 1, end - 1
            dist, count = solve(n, edges, start)
            count1, count2 = count[end][0], count[end][1]
            if dist[end][0] + 1 == dist[end][1]:
                count1 += count2
            print(count1)

    def solve2():
        T = int(input())

        for _ in range(T):
            input()
            n, m = map(int, input().split())
            start, end = map(int, input().split())
            start, end = start - 1, end - 1
            edges = []
            for _ in range(m):
                u, v = map(int, input().split())
                edges.append((u - 1, v - 1, 1))
                edges.append((v - 1, u - 1, 1))
            print(cf776(n, edges, start, end))

    solve2()
