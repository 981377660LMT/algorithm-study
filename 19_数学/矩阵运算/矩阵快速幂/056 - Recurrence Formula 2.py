# 求第N项模1e9+7的值
# a1=1 a2=1 a3=2 an=a(n-1)+a(n-2)+a(n-3)
# N<=1e18

# [an  ]     =  [1 1 1]   * [an-1]
# [an-1]        [1 0 0]     [an-2]
# [an-2]        [0 1 0]     [an-3]
from matqpow import matqpow2
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())

if n <= 3:
    print(2 if n == 3 else 1)
    exit(0)

T = [[1, 1, 1], [1, 0, 0], [0, 1, 0]]
resT = matqpow2(T, n - 3, MOD)
a3, a2, a1 = 2, 1, 1
res = (resT[0][0] * a3 + resT[0][1] * a2 + resT[0][2] * a1) % MOD
print(res)
