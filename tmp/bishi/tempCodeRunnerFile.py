from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

N, M = map(int, input().split())
A, B, C, D, E, F = map(int, input().split())
bad = set()
for _ in range(M):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    bad.add((x, y))

DIR = ((A, B), (C, D), (E, F))


@lru_cache(None)
def dfs(row: int, col: int, remain: int) -> int:
    if remain == 0:
        return 1
    res = 0
    for dr, dc in DIR:
        nr, nc = row + dr, col + dc
        if (nr, nc) not in bad:
            res += dfs(nr, nc, remain - 1)
            res %= MOD
    return res


print(dfs(0, 0, N))
