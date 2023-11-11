# 判断比赛结果是否存在矛盾
# W L D 分别表示 赢 输 平局

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

OK = (("W", "L"), ("L", "W"), ("D", "D"))


def main() -> None:
    n = int(input())
    grid = []
    for _ in range(n):
        grid.append(list(input()))

    for i in range(n):
        for j in range(i + 1, n):
            if (grid[i][j], grid[j][i]) not in OK:
                print("incorrect")
                exit(0)
    print("correct")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
