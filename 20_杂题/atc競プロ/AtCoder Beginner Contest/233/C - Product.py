# 每个袋子里取出一个数
# 问最后取出的数乘积为mul的方案数
# mul<=1e18 n<=16
# !注意到1e18以下正整数约数个数最多1e5个 可以考虑dp

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, target = map(int, input().split())
bags = []
for _ in range(n):
    _, *rest = list(map(int, input().split()))
    bags.append(rest)


@lru_cache(None)
def dfs(gi: int, remain: int) -> int:
    if gi == n:
        return 1 if remain == 1 else 0
    res = 0
    for cur in bags[gi]:
        if remain % cur == 0:
            res += dfs(gi + 1, remain // cur)
    return res


print(dfs(0, target))
