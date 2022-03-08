# 从西北角进去，东南角出来。
# 地里每个道路的交叉点上都有种着一株花生苗，上面有若干颗花生，经过一株花生苗就能摘走该它上面所有的花生。

# Hello Kitty只能向东或向南走，不能向西或向北走。

# 问Hello Kitty最多能够摘到多少颗花生。


n = int(input())

while n:
    n -= 1
    row, col = map(int, input().split())
    mat = []
    for _ in range(row):
        mat.append(list(map(int, input().split())))

    dp = [[0] * col for _ in range(row)]
    for r in range(row):
        for c in range(col):
            top, left = 0, 0
            if 0 <= r - 1 < row:
                top = dp[r - 1][c]
            if 0 <= c - 1 < col:
                left = dp[r][c - 1]
            dp[r][c] = max(top, left) + mat[r][c]
    print(dp[-1][-1])

