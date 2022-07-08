# 爬n节(n<=1e5)楼梯 每次1节或者2节
# 有的楼梯不能爬
# 到达楼顶的爬行方案数
from collections import defaultdict
import sys
import os


def main() -> None:
    # @lru_cache(None)
    # def dfs(index: int) -> int:
    #     if index in bad:
    #         return 0
    #     if index >= n:
    #         return int(index == n)

    #     return (dfs(index + 1) % MOD + dfs(index + 2) % MOD) % MOD

    # n, m = map(int, input().split())
    # bad = set()
    # for _ in range(m):
    #     bad.add(int(input()))
    # res = dfs(0)
    # dfs.cache_clear()
    # print(res)

    n, m = map(int, input().split())
    bad = set()
    for _ in range(m):
        bad.add(int(input()))

    # 如果有坏楼梯 那么dp中途把坏楼梯赋为0
    dp = defaultdict(int, {0: 1})
    for i in range(1, n + 1):
        if i in bad:
            dp[i] = 0
        else:
            dp[i] = (dp[i - 1] + dp[i - 2]) % MOD
    print(dp[n])


if __name__ == "__main__":

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = int(1e9 + 7)
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
