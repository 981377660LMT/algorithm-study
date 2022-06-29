# f(x)=(x の各位の数字の積)
# 整数 N と B が与えられるので、 1 <=m<=N, m−f(m)=B となるものの個数を求めてください。
# N<1e11
# B<1e11

# !典型０２５　あるものだけを数えていく
# 注意到乘积的个数是有限的 因此直接枚举m质因子个数/枚举各位的乘积
from functools import reduce
from math import ceil, log, log2
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def cal(num: int) -> int:
    return reduce(lambda x, y: x * y, [int(i) for i in str(num)], initial=1)


N, B = map(int, input().split())
two = ceil(log2(B))
three = ceil(log(B, 3))
five = ceil(log(B, 5))
seven = ceil(log(B, 7))

res = 0
for f2 in range(two + 1):
    for f3 in range(three + 1):
        for f5 in range(five + 1):
            for f7 in range(seven + 1):
                mul = (2 ** f2) * (3 ** f3) * (5 ** f5) * (7 ** f7)  # 枚举各位的乘积
                if mul + B > N:
                    break
                if cal(mul + B) == mul:
                    res += 1

# 特判
if cal(B) == 0 and B <= N:
    res += 1
print(res)
