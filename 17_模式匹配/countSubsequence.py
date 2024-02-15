# 非空的本质不同的子序列个数

from typing import Any, Sequence


def countSubSequence(seq: Sequence[Any], mod=int(1e9 + 7)) -> int:
    n = len(seq)
    dp = [0] * (n + 1)
    dp[0] = 1
    last = dict()
    for i, c in enumerate(seq):
        dp[i + 1] = 2 * dp[i] % mod
        if c in last:
            dp[i + 1] -= dp[last[c]]
        last[c] = i
    return (dp[n] - 1) % mod


if __name__ == "__main__":
    # https://yukicoder.me/problems/no/1493
    # 给定一个长度为n的数组，每次可以将相邻的两个数换成xor
    # 问可以得到的数组的个数模1e9+7
    MOD = int(1e9 + 7)
    n = int(input())
    nums = list(map(int, input().split()))
    for i in range(n - 1):
        nums[i + 1] ^= nums[i]
    nums = nums[:-1]
    res = countSubSequence(nums, MOD)
    res = (res + 1) % MOD  # 空集
    print(res)
