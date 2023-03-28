# 把nums全部模M(任意选一个大于1的M)
# 求nums中种类的最小值

# !计算偏移量的gcd(按照模分组)
# gcd为0 说明所有的差都是0 说明所有的数都相等
# gcd为1 说明至少两种 取M=2
# gcd>=2 取M=gcd 最少一种

from math import gcd
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
nums = list(map(int, input().split()))
gcd_ = 0
for i in range(n):
    gcd_ = gcd(gcd_, nums[i] - nums[0])
print(2 if gcd_ == 1 else 1)
