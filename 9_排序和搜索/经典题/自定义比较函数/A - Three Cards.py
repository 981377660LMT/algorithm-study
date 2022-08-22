# !选择三个数使得拼接出来的整数最大
# !数字排序:nlargest保证大小
# !字符串排序:转字符串后sorted保证字典序

# !要循环找字典序最大

from heapq import nlargest
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n = int(input())
nums = list(map(int, input().split()))

max3 = nlargest(3, nums)

# 要循环求字典序
