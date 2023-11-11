# n个格子 k种颜色可选
# !1<=[i-j]<=2 不能同色
# 求涂色方案数
# n<=1e18
# k<=1e9
# 第一个格子k种 第二个格子 k-1种 第三个格子 k-2种 第四个格子 k-2种 ...

import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n, k = map(int, input().split())

if n == 1:
    print(k)
    exit(0)
if k == 1:
    print(0)
    exit(0)

res1 = k * (k - 1) % MOD
res2 = pow(k - 2, n - 2, MOD)
print(res1 * res2 % MOD)
