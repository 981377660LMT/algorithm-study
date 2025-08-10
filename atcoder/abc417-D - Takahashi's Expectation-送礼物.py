# abc417-D - Takahashi's Expectation-送礼物
# https://atcoder.jp/contests/abc417/tasks/abc417_d
#
# Takahashi 将会收到 N 个礼物。
# 他有一个称为心情（mood）的参数，该参数是一个非负整数，并且每次收到礼物时心情都会发生变化。
# 每个礼物都有三个参数：价值 P、心情增加量 A 和心情减少量 B。
# 他心情的变化取决于这些参数，具体如下：如果收到的礼物的价值 P 大于等于他当前的心情值，他会感到开心，心情将增加 A。
# 如果收到的礼物的价值 P 小于他当前的心情值，他会感到失望，心情将减少 B。但如果原本的心情值小于 B，那么心情将变为 0。
# 他收到的第 i个礼物具有价值 P_i、心情增加量 A_i 和心情减少量 B_i。
# 你将需要回答 Q 个问题。在第 i 个问题中，给定一个非负整数 X_i，
# 你需要回答以下问题： 当初始心情为 X_i 时，求 Takahashi 收到所有 N 个礼物后的心情值。
#
# !N<=1e4, 1<=P, A, B <= 500, Q<=1e5
#
# !当前心情 >500 则必然会感到失望，只有心情值 <=500 才可能出现开心的状况
# dp[i][j] -> 接受第i个礼物的情绪值为j时，最终的情绪值
# j<=500 => 预处理
# j>500 => 二分+前缀和，使用二分查找来找到失望的终点

from bisect import bisect_right
from itertools import accumulate


def main():
    M = 510 * 2

    N = int(input())
    P, A, B = [0] * N, [0] * N, [0] * N
    for i in range(N):
        P[i], A[i], B[i] = map(int, input().split())

    presumB = [0] + list(accumulate(B))
    dp = [[0] * M for _ in range(N + 1)]
    for i in range(M):
        dp[-1][i] = i
    for i in range(N - 1, -1, -1):
        for j in range(M):
            if j <= P[i]:
                dp[i][j] = dp[i + 1][j + A[i]]
            else:
                dp[i][j] = dp[i + 1][max(0, j - B[i])]

    Q = int(input())
    for _ in range(Q):
        x = int(input())
        if x < M:
            print(dp[0][x])
        else:
            pos = bisect_right(presumB, x - M)
            if pos == N + 1:
                print(x - presumB[-1])
            else:
                print(dp[pos][x - presumB[pos]])


if __name__ == "__main__":
    main()
