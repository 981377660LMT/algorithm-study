from typing import List

# 数组 A 中符合下列属性的任意`子序列` B 称为 “山脉”：
# B.length >= 3
# 先绝对单增后绝对单减

# 返回最长 “山脉” 的长度。
# 数组中求一个最长的山脉的长度，满足山脉中没有两个相邻的相同高度的点，且只有一个峰

# n<=1000


from typing import List
from bisect import bisect_left, bisect_right


def caldp(nums: List[int], isStrict=True) -> List[int]:
    """求每个位置处的LIS长度(包括自身)"""
    if not nums:
        return []
    res = [1] * len(nums)
    LIS = [nums[0]]
    for i in range(1, len(nums)):
        if nums[i] > LIS[-1]:
            LIS.append(nums[i])
            res[i] = len(LIS)
        else:
            pos = bisect_left(LIS, nums[i]) if isStrict else bisect_right(LIS, nums[i])
            LIS[pos] = nums[i]
            res[i] = pos + 1
    return res


class Solution:
    def longestMountain(self, arr: List[int]) -> int:
        n = len(arr)
        up = caldp(arr)
        down = caldp(arr[::-1])[::-1]

        res = []
        for i in range(n):
            res.append(down[i] + up[i] - 1)
        return max(res, default=0)


print(Solution().longestMountain([2, 1, 4, 7, 3, 2, 5]))
print(Solution().longestMountain([1, 1, 1]))
# 输入：[2,1,4,7,3,2,5]
# 输出：5
# 解释：最长的 “山脉” 是 [1,4,7,3,2]，长度为 5。
