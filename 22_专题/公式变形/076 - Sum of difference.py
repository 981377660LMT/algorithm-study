# 所有(aj-ai)的绝对值的和
# 求一维平面上所有点对的(曼哈顿)距离之和
# n<=2e5
# !排序后和不加绝对值一样

from itertools import accumulate
from typing import List


n = int(input())
nums = list(map(int, input().split()))


def dist1d(nums: List[int]) -> int:
    n = len(nums)
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    res = 0
    for i in range(n):
        res += i * nums[i] - preSum[i]
    return res


print(dist1d(nums))

