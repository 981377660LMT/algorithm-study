# 求十进制下由0,2组成的第k小的数
# 1<=k<=1e18
# O(logk)
# 二进制分解


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


k = int(input())
res = []
while k:
    if k % 2 == 1:
        res.append(2)
    else:
        res.append(0)
    k //= 2

res.reverse()
print(int("".join(map(str, res))))
