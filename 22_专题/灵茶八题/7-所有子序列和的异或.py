# 7-所有子序列和的异或
# 01背包
# n<=1000
# sum(nums[i])<2**16
# !dp[i][j]表示前i个数，选出元素和为j的方案数的奇偶性


from typing import List


def xorOfSubsetSum(nums: List[int]) -> int:
    sum_ = sum(nums)
    dp = [0] * (sum_ + 1)
    dp[0] = 1
    for num in nums:
        ndp = dp[:]
        for i in range(sum_ + 1):
            if i + num >= len(ndp):
                break
            ndp[i + num] ^= dp[i]
        dp = ndp

    res = 0
    for i, v in enumerate(dp):
        if v:
            res ^= i
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(xorOfSubsetSum(nums))
