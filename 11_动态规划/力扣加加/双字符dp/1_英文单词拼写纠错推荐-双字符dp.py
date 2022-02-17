# 插入、删除：3
# 修改：(q w e r t a s d f g z x c v ) (y u i o p h j k l b n m)
# 以上两个分组内的字符修改 1 分，两个分组间字符修改 2 分。


LEFT, RIGHT = set(list("qwertasdfgzxcv")), set(list("yuiophjklbnm"))
target, *arr = input().split()


def calDist(src: str, dist: str) -> int:
    n, m = len(src), len(dist)
    dp = [[int(1e20)] * (m) for _ in range(n)]
    if src[0] == dist[0]:
        dp[0][0] = 0
    else:
        if int(src[0] in LEFT) ^ int(dist[0] in LEFT):
            dp[0][0] = 2
        else:
            dp[0][0] = 1

    for i in range(1, n):
        for j in range(m):
            if src[i] == dist[j]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                if int(src[i] in LEFT) ^ int(dist[j] in LEFT):
                    dp[i][j] = min(dp[i][j], dp[i - 1][j - 1] + 2)
                else:
                    dp[i][j] = min(dp[i][j], dp[i - 1][j - 1] + 1)
            dp[i][j] = min(dp[i][j], dp[i - 1][j] + 3, dp[i][j - 1] + 3)

    return dp[-1][-1]


res = sorted(arr, key=lambda x: calDist(x, target))
print(' '.join(res[:3]))

