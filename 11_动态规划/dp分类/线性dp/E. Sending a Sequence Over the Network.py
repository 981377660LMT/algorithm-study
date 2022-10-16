# https://zhuanlan.zhihu.com/p/572692304
# 给定一个数组，将该数组分成连续的若干部分。
# !每一部分的长度大小是该部分最左端或者最右端的大小-1，
# 求能否将这个数组完全分成合法的若干部分。
# n<=1e5
# dp[i]表示前i个数能否分成合法的若干部分
# 每次传入一个a[i]，a[i]可能是某一个序列的最左端也可能是最右端
from typing import List


def split(nums: List[int]) -> bool:
    n = len(nums)
    dp = [False] * (n + 1)
    dp[0] = True
    for i in range(1, n + 1):
        cur = nums[i - 1]
        if i + cur <= n:
            dp[i + cur] |= dp[i - 1]
        if i - cur - 1 >= 0:
            dp[i] |= dp[i - cur - 1]
    return dp[-1]


assert split([1, 1, 3, 4, 1, 3, 2, 2, 3])
