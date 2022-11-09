# https://atcoder.jp/contests/abc162/tasks/abc162_f

# 本题是力扣 198. 打家劫舍 的变形题。
# !输入 n (2≤n≤2e5) 和长为 n 的数组 a (-1e9≤a[i]≤1e9)。
# 数组 a 就是 198 题的房屋存放金额。
# 在本题中，你必须恰好偷 floor(n/2) 个房子，不能偷相邻的房子。
# 输出你能偷窃到的最高金额。

# https://atcoder.jp/contests/abc162/submissions/36315908
# !定义 dp[i] 表示从前 i+1 个数中选 floor(i+1/2) 个数的最大收益。
# dp[0]=0,dp[1]=max(nums[0],nums[1])。
# 对于 i≥2，
# !如果选，那么 dp[i] = dp[i-2] + nums[i]。
# !如果不选：
# !- 如果 i 为奇数，那么前面的选法是固定的，也就是所有偶数下标的和，记作 s。
# !- 如果 i 为偶数，那么需要选 (i+1)//2 个数，这恰好是 dp[i-1]。
# 选或不选取 max。
# 答案为 f[n]。
# TODO O(1) 空间写法

from typing import List

INF = int(1e18)


def selectHalf(n: int, nums: List[int]) -> int:
    """恰好选择n//2 (n<=1e5) 个数,使得和最大"""
    if n <= 1:
        return 0

    preEvenSum = [0] * (n + 1)  # 每个前缀偶数下标的和
    for i in range(1, n + 1):
        preEvenSum[i] = (preEvenSum[i - 1] + nums[i - 1]) if i & 1 else preEvenSum[i - 1]

    dp = [-INF] * n
    dp[0] = 0
    dp[1] = max(nums[0], nums[1])
    for i in range(2, n):  # 从前i+1个数中选 floor(i+1/2) 个数的最大收益
        # not jump
        cand1 = dp[i - 2] + nums[i]
        # jump
        if i % 2 == 0:
            cand2 = dp[i - 1]
        else:
            cand2 = preEvenSum[i + 1]
        dp[i] = max(cand1, cand2)

    return dp[n - 1]


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(selectHalf(n, nums))
