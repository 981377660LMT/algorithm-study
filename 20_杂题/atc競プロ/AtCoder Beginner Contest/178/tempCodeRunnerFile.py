# 跳楼梯 每次跳>=3节
# 求跳到S节的方案数 (S<=2000)

from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    @lru_cache(None)
    def dfs(remain: int) -> int:
        if remain <= 0:
            return int(remain == 0)

        res = 0
        for cur in range(3, remain + 1):
            res += dfs(remain - cur) % MOD
        return res

    res = dfs(int(input()))
    dfs.cache_clear()
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
