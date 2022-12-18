# 双六2
# 从0出发 投掷一个m面骰子1-m 每次走出骰子点数的格子
# !如果走到bad[i]就会回到0 从0开始 求走到>=n的期望步数
# n,m<=1e5 len(bads)<10

# dp[i]表示从i走到n的期望步数
# !如果i能走,dp[i] = 1 + (dp[i+1] + dp[i + 2] + ... + dp[i + m] / m
# !如果不能走,dp[i]=dp[0]
from typing import List


def sugoroku2(n: int, bads: List[int], m: int) -> int:
    visited = [False] * (n + m + 10)
    for bad in bads:
        visited[bad] = True
    dp = [1] * (n + m + 10)
    dpSum = [0] * (n + m + 10)


n, m, k = map(int, input().split())
bads = list(map(int, input().split()))  # 走到bads[i]就会回到0
print(sugoroku2(n, bads, m))
