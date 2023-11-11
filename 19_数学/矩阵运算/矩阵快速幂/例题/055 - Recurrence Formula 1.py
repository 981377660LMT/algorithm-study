# 求第N项模1e9+7的值
# a1=1 a2=1 an=2*a(n-1)+a(n-2)
# N<=1e18

# [an  ]     =  [2 1]   * [an-1]
# [an-1]        [1 0]     [an-2]
from matqpow import matqpow1, matmul
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())

if n <= 2:
    print(1)
    exit(0)

res = [[1], [1]]
trans = [[2, 1], [1, 0]]
tmp = matqpow1(trans, n - 2, MOD)
res = matmul(tmp, res, MOD)  # 第N项的值
print(res[0][0])
