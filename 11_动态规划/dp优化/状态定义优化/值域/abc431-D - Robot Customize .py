# D - Robot Customize
# https://atcoder.jp/contests/abc431/tasks/abc431_d
# 有一个头和身体组成的机器人。
# 该机器人可以同时安装重量为 wi 的零件，每个零件安装在头部和身体上时会产生不同的幸福值 hi 和 bi。
# 高桥希望将全部 N 种零件各安装一个到机器人上。请找出在不使机器人倾倒的前提下，所有零件美观度总和的最大可能值。
# N<=500, Wi<=500, Hi<=1e9，Bi<=1e9
#
# !dp[i][j] 表示前i个零件，头和为j时的最大美观度

INF = int(4e18)

if __name__ == "__main__":
    N = int(input())
    W, H, B = [0] * N, [0] * N, [0] * N
    for i in range(N):
        w, h, b = map(int, input().split())
        W[i], H[i], B[i] = w, h, b

    wSum = 0
    dp = [0]
    for w, h, b in zip(W, H, B):
        ndp = [-INF] * (wSum + w + 1)
        for j in range(wSum + 1):
            ndp[j + w] = max(ndp[j], dp[j] + h)  # head
            ndp[j] = max(ndp[j], dp[j] + b)  # body
        dp = ndp
        wSum += w

    res = 0
    for i in range(wSum + 1):
        if i <= wSum - i:
            res = max(res, dp[i])
    print(res)
