# 给你一个包含若干 `互不相同` 整数的数组 nums ，你需要执行以下操作 直到数组为空 ：


# 如果数组中第一个元素是当前数组中的 最小值 ，则删除它。
# 否则，将第一个元素移动到数组的 末尾 。
# 请你返回需要多少个操作使 nums 为空。
# !查找最小值的位置(类似于排队买票)
# !查找最小值的位置转化为预先处理每个数的索引，按顺序维护`最小值索引`而不是最小值

from typing import List


class Solution:
    def countOperationsToEmptyArray(self, nums: List[int]) -> int:
        n = len(nums)
        order = sorted(range(n), key=lambda i: nums[i])
        res = 0
        return res + n
