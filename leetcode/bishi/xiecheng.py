# 将k个数变为平均数后 最小化最大值与最小值的差
# !前后缀分解 枚举前缀变几个数
from itertools import accumulate
from typing import List

INF = int(1e18)


def minimizeDiff(nums: List[int], k: int) -> float:
    nums.sort()
    preSum = [0] + list(accumulate(nums))
    sufSum = ([0] + list(accumulate(nums[::-1])))[::-1]
    res = INF
    for i in range(k):
        min_, max_, avg = nums[i], nums[-k + i], (preSum[i] + sufSum[-k + i]) / k
        res = min(res, max(min_, max_, avg) - min(min_, max_, avg))
    return res


print(minimizeDiff([1, 2, 3, 4, 5], 2))
