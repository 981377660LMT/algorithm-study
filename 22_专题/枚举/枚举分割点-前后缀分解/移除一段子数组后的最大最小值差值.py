# 给一个数组10^5长度，和一个整数K，小于数组长度
# 在数组中取出K个"连续"数后（滑窗）， 剩下的数字最大值和最小值之差最小是多少


# 维护前缀后缀的最小值、最大值，一共四个数组即可，然后枚举连续区间边界

from typing import List


def minDiff(nums: List[int], k: int) -> int:
    n = len(nums)

    preMin = [nums[0]] * n
    preMax = [nums[0]] * n
    sufMin = [nums[-1]] * n
    sufMax = [nums[-1]] * n
    for i in range(1, n):
        preMin[i] = min(preMin[i - 1], nums[i])
        preMax[i] = max(preMax[i - 1], nums[i])
        sufMin[~i] = min(sufMin[-i], nums[~i])
        sufMax[~i] = max(sufMax[-i], nums[~i])

    res = int(1e20)
    left = 0
    for right in range(n):
        if right >= k:
            left += 1
        if right >= k - 1:
            min1, max1 = (
                preMin[left - 1] if left - 1 >= 0 else int(1e20),
                preMax[left - 1] if left - 1 >= 0 else -int(1e20),
            )
            min2, max2 = (
                sufMin[right + 1] if right + 1 < n else int(1e20),
                sufMax[right + 1] if right + 1 < n else -int(1e20),
            )

            min_ = min(min1, min2)
            max_ = max(max1, max2)
            res = min(res, max_ - min_)

    return res


print(minDiff([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0, 0, 0], 3))

