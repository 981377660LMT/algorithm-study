# 给你一个包含若干 `互不相同` 整数的数组 nums ，你需要执行以下操作 直到数组为空 ：


# 如果数组中第一个元素是当前数组中的 最小值 ，则删除它。
# 否则，将第一个元素移动到数组的 末尾 。
# 请你返回需要多少个操作使 nums 为空。

from typing import List


class Solution:
    def countOperationsToEmptyArray(self, nums: List[int]) -> int:
        n = len(nums)
        order = sorted(range(n), key=lambda i: nums[i])  # 出队顺序
        res = 0
        for i in range(1, n):
            pos1, pos2 = order[i - 1], order[i]
            if pos2 < pos1:  # !如果逆序，则需要移动n-i个位置
                res += n - i
        return res + n


assert Solution().countOperationsToEmptyArray([1, 2, 3, 4, 5]) == 5
