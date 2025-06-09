# 3576. 数组元素相等转换
# https://leetcode.cn/problems/transform-array-to-all-equal-elements/description/
# 给你一个大小为 n 的整数数组 nums，其中只包含 1 和 -1，以及一个整数 k。
# 你可以最多进行 k 次以下操作：
# 选择一个下标 i（0 <= i < n - 1），然后将 nums[i] 和 nums[i + 1] 同时 乘以 -1。
# 注意：你可以在 不同 的操作中多次选择相同的下标 i。
# 如果在最多 k 次操作后可以使数组的所有元素相等，则返回 true；否则，返回 false。


from typing import List


class Solution:
    def canMakeEqual(self, nums: List[int], k: int) -> bool:
        def check(target: int) -> bool:
            n = len(nums)
            remain = k
            mul = 1
            for i, v in enumerate(nums):
                if v * mul == target:
                    mul = 1  # 不操作，下一个数不乘 -1
                    continue
                if remain == 0 or i == n - 1:
                    return False
                remain -= 1
                mul = -1
            return True

        return check(1) or check(-1)
