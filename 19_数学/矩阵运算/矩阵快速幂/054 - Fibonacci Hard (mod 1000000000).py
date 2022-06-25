# 求第N项模1e9的值
# a1=1 a2=1 an=a(n-1)+a(n-2)
# N<=1e18

# [an  ]     =  [1 1]   * [an-1]
# [an-1]        [1 0]     [an-2]
from matqpow import matqpow1
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9)

n = int(input())

if n <= 2:
    print(1)
    exit(0)

T = [[1, 1], [1, 0]]  # 转移矩阵
resT = matqpow1(T, n - 2, MOD)  # 转移矩阵的幂次
a2, a1 = 1, 1
res = (resT[0][0] * a2 + resT[0][1] * a1) % MOD  # 第N项的值
print(res)

