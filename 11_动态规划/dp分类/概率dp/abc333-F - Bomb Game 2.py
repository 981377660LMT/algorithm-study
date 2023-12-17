# https://atcoder.jp/contests/abc333/tasks/abc333_f
# 概率dp，状态转移消环技巧
# n 个人排成一排，进行操作：
# 对于队首的人，每轮操作中有1/2的概率移出队伍，1/2的概率移动到队尾
# 对i=1,2,...,n，问第i个人最后获胜的概率(队伍中仅剩一人时获胜)，对998244353取模

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")

from typing import List


def modInv(a: int, mod: int) -> int:
    return pow(a, mod - 2, mod)


MOD = 998244353
INV2 = modInv(2, MOD)


def solve(n: int, win: int, shift: int, out: int) -> List[int]:
    dp = [0] * (n + 1)
    a, tmp = [0] * (n + 1), [0] * (n + 1)
    dp[1] = 1
    for i in range(2, n + 1):
        a[1] = win
        for j in range(2, i + 1):
            a[j] = out * dp[j - 1] % MOD
        sum_, q = 0, 1
        for j in range(1, i + 1):
            sum_ = (sum_ * shift + a[j]) % MOD
            q = q * shift % MOD
        tmp[i] = sum_ * modInv(1 - q, MOD) % MOD
        tmp[0] = tmp[i]
        for j in range(1, i):
            tmp[j] = (a[j] + tmp[j - 1] * shift) % MOD
        for j in range(1, i + 1):
            dp[j] = tmp[j]
    return dp


if __name__ == "__main__":
    N = int(input())
    dp = solve(N, 0, INV2, INV2)  # 直接获胜、移动到队尾、移出队伍的概率
    for i in range(1, N + 1):
        print(dp[i], end=" ")
