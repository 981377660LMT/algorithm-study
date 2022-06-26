# 高橋君は料理 1 から N の N 品の料理を作ろうとしています。
# 料理 i はオーブンを連続した Ti分間使うことで作れます。
# 2 つのオーブンを使えるとき、N 品の料理を全て作るまでに最短で何分かかりますか？
# n<=100
# Ti<=1000

# !注意有两个烤箱
# 因此最佳的策略是将料理分为相近的两半 取两者较大值
# 问题转换为01背包

import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)
n = int(input())
costs = list(map(int, input().split()))
V = 100 * 1000
dp = [False] * (V + 10)
dp[0] = True

for i in range(n):
    for j in range(V, costs[i] - 1, -1):
        dp[j] = dp[j] or dp[j - costs[i]]

start = (sum(costs) + 1) // 2
print(dp.index(True, start))

