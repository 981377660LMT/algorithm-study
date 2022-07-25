# 给定N,在可用位为N中有多少个数满各位差的绝对值<=1且为1-9?
# !字典dp比数组dp慢太多

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    n = int(input())
    dp = [1] * 10
    for _ in range(1, n):
        ndp = [0] * 10
        for j in range(1, 10):
            for pre in (j - 1, j, j + 1):
                if pre < 1 or pre > 9:
                    continue
                ndp[j] += dp[pre]
                ndp[j] %= MOD
        dp = ndp

    res = 0
    for count in dp[1:]:
        res += count
        res %= MOD
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
