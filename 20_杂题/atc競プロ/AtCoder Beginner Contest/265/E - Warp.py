# 你现在站在点(0,0)上，你有n次传送的机会，
# 且给出M个点(xi,yi) 这M个点是不可以传送到的,然后你有3种传送方式∶
# (1)从(x, y)- >(x+A,y+B)
# (1)从(x, y)- >(x+C,y+D)
# (1)从(x, y)- >(x+E,y+F)
# 现在问你在n次传送下可以形成多少条路径，答案对998244353取模.
# n<=300 m<=1e5

# !注意到在n的限制下 固定n时能走到的坐标(x,y)数量很少 => 组合数 O(n^2)
# !因此可以用dp[index][坐标] 用字典来dp
# O(n^3)

# 优化:
# !1.不用循环 用三条语句会快300ms
# !2.dp和ndp 滚动dp
# !3.字典存一维(哈希值或者元组)
# !4.不要重新创建元组 会快300ms
# !5.取模合成一条语句才不会TLE
# ndp[next] = (ndp[next] + preV) % MOD
# 比
# ndp[next] += preV
# ndp[next] %= MOD
# 快很多

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m = map(int, input().split())
A, B, C, D, E, F = map(int, input().split())
bad = set(tuple(map(int, input().split())) for _ in range(m))

dp = defaultdict(int, {(0, 0): 1})
for _ in range(n):
    ndp = defaultdict(int)
    for (preX, preY), preV in dp.items():
        next1 = (preX + A, preY + B)
        if next1 not in bad:
            ndp[next1] = (preV + ndp[next1]) % MOD
        next2 = (preX + C, preY + D)
        if next2 not in bad:
            ndp[next2] = (preV + ndp[next2]) % MOD
        next3 = (preX + E, preY + F)
        if next3 not in bad:
            ndp[next3] = (preV + ndp[next3]) % MOD
    dp = ndp


res = 0
for v in dp.values():
    res = (res + v) % MOD
print(res)
