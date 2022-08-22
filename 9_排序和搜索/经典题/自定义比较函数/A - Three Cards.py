# 拼接最大数
# n个数中选取三个数 任意顺序拼接
# !求能拼接的最大数
# n<=2e5 ai<=1e6

# !选择三个数使得拼接出来的整数最大
# !数字排序:nlargest保证大小
# !字典序排序:要循环找字典序最大
# 反例: 98 987 23 12

from heapq import nlargest
from itertools import permutations
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n = int(input())
nums = list(map(int, input().split()))

max3 = nlargest(3, nums)
res = 0
for perm in permutations(max3):
    cand = int("".join(map(str, perm)))
    res = cand if cand > res else res
print(res)
