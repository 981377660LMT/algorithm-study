# 可以在字符串头部添加任意多个a
# 问最后字符串能否成为回文

# !看两端a的个数
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

s = input()
la = len(s) - len(s.lstrip("a"))
ra = len(s) - len(s.rstrip("a"))
if la > ra:
    print("No")
    exit(0)
res = "a" * (ra - la) + s
print("Yes" if res == res[::-1] else "No")
