# 黑板上有n类整数
# 第i类整数为Ai 有Bi个
# 每次选择两个数x和y 如果(x+y)为素数 那么可以把x和y消去
# 问最多能够操作多少次

# n<=100
# Ai<=1e7 Bi<=1e9

# !注意到Ai->Aj连边 这个图是二分图
from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
counter = defaultdict(int)
for _ in range(n):
    a, b = map(int, input().split())
    counter[a] += b

# TODO
# 好难...
