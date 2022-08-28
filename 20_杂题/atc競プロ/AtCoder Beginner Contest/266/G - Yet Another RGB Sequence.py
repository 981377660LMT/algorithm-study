"""
求RGB串的个数
R G B 的个数分别为R G B 个
RG子串的个数为K个


将RG视为A 则等价于
K个A (R-K)个R (G-K)个G B个B 组成的没有RG子串的字符串个数
"""


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# !RG 离散R、B 随意排列，然后剩下的R投到RG和B的桶上

R, G, B, K = map(int, input().split())
