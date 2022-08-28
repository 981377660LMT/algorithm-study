# !掷骰子的点数期望值
# If it is the N-th turn now, your score is X, and the game ends.
# Otherwise, choose whether to continue or end the game.
# If you end the game, your score is X, and there is no more turn.
# Find the expected value of your score when you play the game to maximize this expected value.

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())


# !原来这样算的

dp = [0, 3.5, 4.25]
for i in range(n):
    sum_ = 0
    for cur in range(1, 7):
        sum_ += max(cur, dp[-1])  # !如果本次更高，则在本轮终止
    dp.append(sum_ / 6)
print(dp[n])
