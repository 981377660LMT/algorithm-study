from typing import List


class Solution:
    def wiggleMaxLength(self, nums: List[int]) -> int:
        n = len(nums)
        up, down = 1, 1  # 结尾上升，结尾下降
        for i in range(1, n):
            if nums[i - 1] < nums[i]:
                up = max(up, down + 1)
            elif nums[i - 1] > nums[i]:
                down = max(down, up + 1)
        return max(up, down)


print(Solution().wiggleMaxLength(nums=[1, 17, 5, 10, 13, 15, 10, 5, 16, 8]))
