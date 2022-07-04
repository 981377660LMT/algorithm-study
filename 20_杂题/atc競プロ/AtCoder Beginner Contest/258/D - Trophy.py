# 每关都有长为s的故事时间和长为g的游戏时间

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, x = map(int, input().split())
    stories, games = [], []
    for _ in range(n):
        s, g = map(int, input().split())
        stories.append(s)
        games.append(g)

    # 枚举玩到哪一关
    res = int(9e18)
    base = 0
    for i in range(n):
        # 玩到第i关的时间
        base += stories[i] + games[i]
        # 注意这里的max(0)
        res = min(res, base + max(0, (x - (i + 1))) * games[i])
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
