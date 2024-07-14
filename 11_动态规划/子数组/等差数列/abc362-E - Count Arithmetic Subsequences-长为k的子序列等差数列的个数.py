# abc362-E - Count Arithmetic Subsequences
# https://atcoder.jp/contests/abc362/tasks/abc362_e
# 给定一个长度为 n 的整数数组,
# 对k=1,2,...,n，求长为k的子序列中等差数列的个数。
# n<=80
#
# O(n^3)
# dp[i][k][diff] 表示以 nums[i] 结尾的公差为 diff 的等差数列个数，长度为 k.

from collections import defaultdict
from typing import List


def countArithmeticSubsequences(nums: List[int]) -> List[int]:
    n = len(nums)
    dp = [[defaultdict(int) for _ in range(n + 1)] for _ in range(n)]
    for i in range(n):
        for j in range(i):
            diff = nums[i] - nums[j]
            for k in range(n):
                dp[i][k + 1][diff] += dp[j][k][diff]
            dp[i][2][diff] += 1

    res = [0] * (n + 1)
    res[1] = n
    for row in dp:
        for len_, d in enumerate(row):
            for count in d.values():
                res[len_] += count
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353

    n = int(input())
    nums = list(map(int, input().split()))
    res = countArithmeticSubsequences(nums)

    for i in range(1, n + 1):
        print(res[i] % MOD)
