"""
给你一个小写字母组成的字符串s
!取字符串s的子序列,不连续取字符
!问最后能形成多少种`不同的子序列`
len(s)<=2e5

设dp[i]表示前i个字符中能形成的不同子序列的个数
!dp[i] = dp[i-1] + dp[i-2] - dp[last[s[i]]] (当前选或不选)
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def distinctSubseq3(s: str) -> int:
    """不连续取字符,求最后能形成多少种`不同的子序列`"""
    n = len(s)
    dp = [0] * (n + 2)
    dp[0] = 1
    dp[1] = 1
    last = dict()
    for i, char in enumerate(s):
        dp[i + 2] = (dp[i] + dp[i + 1]) % MOD
        if char in last:
            dp[i + 2] -= dp[last[char]]
        last[char] = i
    return (dp[n + 1] - 1) % MOD


if __name__ == "__main__":
    s = input()
    print(distinctSubseq3(s))
