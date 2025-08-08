# E - Hungry Takahashi
# https://atcoder.jp/contests/abc415/tasks/abc415_e
#
# 在 H×W 网格中，从 (1,1) 到 (H,W) 的任何路径，每走一步先收集 A[i][j] 枚硬币，再消费对应天数 Pk 枚硬币，手上硬币不能负数。
# 求使得存在一条路径一直不饿死所需最小初始硬币 x。
#
# 类似 `地下城游戏`

INF = int(4e18)


def solve():
    H, W = map(int, input().split())
    A = [list(map(int, input().split())) for _ in range(H)]
    P = list(map(int, input().split()))

    dp = [[INF] * W for _ in range(H)]
    dp[H - 1][W - 1] = 0

    for i in range(H - 1, -1, -1):
        for j in range(W - 1, -1, -1):
            if i + 1 < H:
                dp[i][j] = min(dp[i][j], dp[i + 1][j])
            if j + 1 < W:
                dp[i][j] = min(dp[i][j], dp[i][j + 1])
            dp[i][j] += P[i + j] - A[i][j]
            dp[i][j] = max(dp[i][j], 0)

    print(dp[0][0])


if __name__ == "__main__":
    solve()
