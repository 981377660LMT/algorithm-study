from collections import deque
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)
# 2 次元平面上に N 個の街と M 個の宝箱があります。街 i は座標 (X
# i
# ​
#  ,Y
# i
# ​
#  ) に、宝箱 i は座標 (P
# i
# ​
#  ,Q
# i
# ​
#  ) にあります。

# 高橋君は原点を出発し、N 個の街全てを訪れたのち原点に戻る旅行をしようと考えています。
# 宝箱を訪れる必要はありませんが、宝箱の中にはそれぞれブースターが 1 つあり、ブースターを拾うごとに移動速度が 2 倍になります。

# 高橋君の最初の移動速度が単位時間あたり 1 であるとき、旅行にかかる時間の最小値を求めてください。


def long_long_bit_count(n: int) -> int:
    c = (n & 0x5555555555555555) + ((n >> 1) & 0x5555555555555555)
    c = (c & 0x3333333333333333) + ((c >> 2) & 0x3333333333333333)
    c = (c & 0x0F0F0F0F0F0F0F0F) + ((c >> 4) & 0x0F0F0F0F0F0F0F0F)
    c = (c & 0x00FF00FF00FF00FF) + ((c >> 8) & 0x00FF00FF00FF00FF)
    c = (c & 0x0000FFFF0000FFFF) + ((c >> 16) & 0x0000FFFF0000FFFF)
    c = (c & 0x00000000FFFFFFFF) + ((c >> 32) & 0x00000000FFFFFFFF)
    return c


if __name__ == "__main__":
    n, m = map(int, input().split())
    street = [tuple(map(int, input().split())) for _ in range(n)]
    box = [tuple(map(int, input().split())) for _ in range(m)]
    target = (1 << n) - 1
    points = street + box + [(0, 0)]
    pairCost = [[0] * (n + m + 1) for _ in range(n + m + 1)]
    for i in range(n + m + 1):
        for j in range(i + 1, n + m + 1):
            x1, y1, x2, y2 = *points[i], *points[j]
            cur = ((x1 - x2) ** 2 + (y1 - y2) ** 2) ** 0.5
            pairCost[i][j] = pairCost[j][i] = cur

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

    res = dfs(n + m, 0)
    # dfs.cache_clear()
    print(res)

    res = INF
    dist = [[INF] * (1 << (n + m)) for _ in range(n + m + 1)]
    dist[n + m][0] = 0
    queue = deque([(0, n + m, 0)])
    while queue:
        curDist, cur, visited = queue.popleft()
        if visited & target == target:
            speed = 2 ** long_long_bit_count(visited >> n)
            cand = pairCost[cur][n + m] / speed + curDist
            res = cand if cand < res else res
            continue

        if curDist > dist[cur][visited]:
            continue

        speed = 2 ** long_long_bit_count(visited >> n)
        for next in range(n + m):
            if visited & (1 << next):
                continue
            cost = pairCost[cur][next] / speed
            cand = cost + curDist
            if dist[next][visited | (1 << next)] > cand:
                dist[next][visited | (1 << next)] = cand
                queue.append((cand, next, visited | (1 << next)))

    print(res)
