# !求 Σ[k=0..10^100]floor(X／10^k)  的值
# !x<=1e(5e5)
# https://atcoder.jp/contests/abc233/editorial/3193

# 公式变形:
# !floor(a/b) = a/b - a%b/b
# !num%10^k 的意义是什么?  num的十进制下第k位的数值
# 这一部分，对于每个数位x，可以计算出它的贡献。

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

x = input()
digitSum = sum(int(char) for char in x)
n = int(x)
print((10 * n - digitSum) // 9)
