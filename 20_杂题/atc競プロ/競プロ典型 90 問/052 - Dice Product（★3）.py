# n个骰子 一个骰子的六个面点数不同
# 求n个骰子出目的所有可能的点数乘积之和
# n<=100

# 数学公式变形
# ∑ai*bj*ck...=∑ai*∑bj*∑ck...
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())

res = 1
for _ in range(n):
    nums = list(map(int, input().split()))
    res *= sum(nums)
    res %= MOD
print(res)

