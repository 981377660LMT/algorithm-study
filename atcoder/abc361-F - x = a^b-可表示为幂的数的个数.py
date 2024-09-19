# abc361-F - x = a^b
# https://atcoder.jp/contests/abc361/tasks/abc361_f
# 给定n，求x∈[1,n]，满足存在a,b(b≥2)，使得x=a^b。
# n≤1e18。
# !可表示为幂的数的个数.
#
# b=2，不好做，先看b>=3的情况。这个时候a最多是n的1/3次方，可以枚举.
#
# !因此，分为 a<=n^(1/3) 和 a>n^(1/3) 两种情况。
# !第一种情况，可以枚举a，算出所有可能；
# !第二种情况，b=2，减去第一种情况的重复部分即可。


from math import isqrt
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    up = int(pow(N, 1 / 3)) + 1
    s = set([1])  # 特判1
    for a in range(2, up + 1):  # 2<=a<=up
        x = a * a
        while x <= N:
            s.add(x)
            x *= a

    res = len(s)
    half = isqrt(N)  # up+1<=a<=root，此时一定有b=2
    if half >= up:
        res += half - (up + 1) + 1
        res -= sum(up + 1 <= x <= half for x in s)

    print(res)
