# 插入、删除：3
# 修改：(q w e r t a s d f g z x c v ) (y u i o p h j k l b n m)
# 以上两个分组内的字符修改 1 分，两个分组间字符修改 2 分。
# 编辑距离

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


LEFT, RIGHT = set(list("qwertasdfgzxcv")), set(list("yuiophjklbnm"))
target, *arr = input().split()


def calDist(src: str, dist: str) -> int:
    n, m = len(src), len(dist)
    dp = [[INF] * (m + 1) for _ in range(n + 1)]
    dp[0][0] = 0

    for i in range(1, n + 1):
        for j in range(1, m + 1):
            if src[i - 1] == dist[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                if (src[i - 1] in LEFT) ^ (dist[j - 1] in LEFT):
                    dp[i][j] = min(dp[i][j], dp[i - 1][j - 1] + 2)  # 两个分组间字符修改
                else:
                    dp[i][j] = min(dp[i][j], dp[i - 1][j - 1] + 1)  # 两个分组内字符修改
                dp[i][j] = min(dp[i][j], dp[i - 1][j] + 3, dp[i][j - 1] + 3)  # 插入、删除

    return dp[-1][-1]


res = sorted(arr, key=lambda x: calDist(x, target))
print(" ".join(res[:3]))
