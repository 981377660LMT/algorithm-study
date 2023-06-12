# https://leetcode.cn/circle/discuss/yIRkxI/
# 长度大于1的子数组和为非负数的最少修改次数
# 最少进行多少次运算（其实就是替代多少次值了），能够使得arr的任意一个长度大于1的子数组和都是非负数
# !=> 只需要考察长度为2的子数组和长度为3的子数组
# 更长的区间可以拆分为长度为 2和3的子数组, 每个组内只要一个数修改就可以了

from itertools import accumulate
from typing import List
from 用最少数量的箭引爆气球 import findMinArrowShots


def minModify(nums: List[int]) -> int:
    n = len(nums)
    intervals = []  # 长为2和3,和为负的子数组起点和终点(闭区间)
    preSum = [0] + list(accumulate(nums))

    for i in range(n - 1):
        if preSum[i + 2] - preSum[i] < 0:
            intervals.append((i, i + 1))
        if i + 3 <= n and preSum[i + 3] - preSum[i] < 0:
            intervals.append((i, i + 2))

    return findMinArrowShots(intervals)


if __name__ == "__main__":
    assert minModify([-1, 1, -1]) == 1
    assert minModify([-1, -1, -1, -1]) == 2
