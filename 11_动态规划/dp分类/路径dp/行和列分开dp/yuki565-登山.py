# No.595 登山(爬山)
# https://yukicoder.me/problems/no/595
# 给定一些山，高度为h[i].
# 每次移动可以有两种方式：
# - 移动到相邻的山j，花费max(0, h[i]-h[j])的代价
# - 花费c的代价, 瞬间移动到任意山(Warp).
# 问最少花费多少代价，可以访问所有山.


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, C = map(int, input().split())
    H = list(map(int, input().split()))

    dp = [INF, 0, INF]  # (jump, right, left)
    for h1, h2 in zip(H, H[1:]):
        ndp = [INF, INF, INF]
        ndp[0] = min(dp) + C
        ndp[1] = min(dp[0], dp[1]) + max(0, h2 - h1)
        ndp[2] = min(dp[0], dp[2]) + max(0, h1 - h2)
        dp = ndp
    print(min(dp))
