# 高橋君は次の操作を好きなだけ行うことができます。
# !レベル n の赤い宝石 (n は 2 以上) を「レベル n−1 の赤い宝石 1 個と、レベル n の青い宝石 X 個」に変換する。
# !レベル n の青い宝石 (n は 2 以上) を「レベル n−1 の赤い宝石 1 個と、レベル n−1 の青い宝石 Y 個」に変換する。
# !高橋君はレベル N の赤い宝石を 1 個持っています。(他に宝石は持っていません。)
# !高橋君はレベル 1 の青い宝石ができるだけたくさんほしいです。操作によって高橋君はレベル 1 の青い宝石を最大何個入手できますか？
# 1≤N≤10
# 1≤X≤5
# 1≤Y≤5

# 数据量暗示直接dfs记忆化


from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    @lru_cache(None)
    def dfs(level: int, isRed: bool) -> int:
        """找1"""
        if level == 1:
            return 0 if isRed else 1
        if isRed:
            return dfs(level - 1, True) + dfs(level, False) * x
        else:
            return dfs(level - 1, True) + dfs(level - 1, False) * y

    n, x, y = map(int, input().split())
    print(dfs(n, True))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
