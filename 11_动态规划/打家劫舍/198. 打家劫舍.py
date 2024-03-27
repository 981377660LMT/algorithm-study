# 打家劫舍
# 不能偷窃相邻的房间
# 给定一个代表每个房屋存放金额的非负整数数组，
# 计算你 不触动警报装置的情况下 ，一夜之内能够偷窃到的最高金额。
# 1 <= nums.length <= 100
# 0 <= nums[i] <= 400

# 一般打家劫舍系列都可以写dp0/dp1

from typing import List


def max(x, y):
    if x > y:
        return x
    return y


def rob(nums: List[int]) -> int:
    if not nums:
        return 0

    dp0, dp1 = 0, nums[0]
    for i in range(1, len(nums)):
        dp0, dp1 = max(dp0, dp1), dp0 + nums[i]
    return max(dp0, dp1)


# 数轴上有n个位置,每个位置有分数scores[i]
# 不能偷相邻的位置(x-1,x+1),求最大分数
# !positions[i]<=1e9 且单调递增
def rob2(positions: List[int], scores: List[int]) -> int:
    if not positions:
        return 0

    n = len(positions)
    dp0, dp1 = 0, scores[0]  # 不偷当前/偷当前 的方案数
    for i in range(1, n):
        dist = positions[i] - positions[i - 1]
        if dist <= 1:
            dp0, dp1 = max(dp0, dp1), max(dp0 + scores[i], dp1)
        else:
            dp0, dp1 = max(dp0, dp1), max(dp0, dp1) + scores[i]
    return max(dp0, dp1)


assert rob([1, 2, 3, 1]) == 4
assert rob2(list(range(4)), [1, 2, 3, 1]) == 4
