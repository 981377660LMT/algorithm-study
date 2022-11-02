# 访问所有城市的最短时间
# 从原点出发，起始速度为1，有n个城市，
# 有m个盒子，经过盒子会让速度加倍，
# 问拜访这n个城市最终回到原点的最短时间

# !dfs(index,visited)

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)


def long_long_bit_count(n: int) -> int:
    c = (n & 0x5555555555555555) + ((n >> 1) & 0x5555555555555555)
    c = (c & 0x3333333333333333) + ((c >> 2) & 0x3333333333333333)
    c = (c & 0x0F0F0F0F0F0F0F0F) + ((c >> 4) & 0x0F0F0F0F0F0F0F0F)
    c = (c & 0x00FF00FF00FF00FF) + ((c >> 8) & 0x00FF00FF00FF00FF)
    c = (c & 0x0000FFFF0000FFFF) + ((c >> 16) & 0x0000FFFF0000FFFF)
    c = (c & 0x00000000FFFFFFFF) + ((c >> 32) & 0x00000000FFFFFFFF)
    return c


if __name__ == "__main__":

    @lru_cache(None)
    def dfs(cur: int, visited: int) -> float:
        if visited & target == target:
            speed = 2 ** long_long_bit_count(visited >> n)
            return dist[cur][n + m] / speed
        speed = 2 ** long_long_bit_count(visited >> n)
        res = INF
        for next in range(n + m):
            if visited & (1 << next):
                continue
            cost = dist[cur][next] / speed
            cand = cost + dfs(next, visited | (1 << next))
            res = cand if cand < res else res
        return res

    n, m = map(int, input().split())
    street = [tuple(map(int, input().split())) for _ in range(n)]
    box = [tuple(map(int, input().split())) for _ in range(m)]
    target = (1 << n) - 1
    points = street + box + [(0, 0)]
    dist = [[0] * (n + m + 1) for _ in range(n + m + 1)]
    for i in range(n + m + 1):
        for j in range(i + 1, n + m + 1):
            x1, y1, x2, y2 = *points[i], *points[j]
            cur = ((x1 - x2) ** 2 + (y1 - y2) ** 2) ** 0.5
            dist[i][j] = dist[j][i] = cur

    res = dfs(n + m, 0)
    dfs.cache_clear()
    print(res)
