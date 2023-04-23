# 最长摆动序列(zigzag sequence)：序列中相邻元素之间的差值正负交替出现

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


def longestZigZagSubsequence(nums: List[int]) -> int:
    n = len(nums)
    res = 1
    prev = -1
    for i in range(n):
        j = i + 1
        while j < n and nums[i] == nums[j]:
            j += 1
        if j < n:
            sign = nums[i] < nums[j]
            if prev != sign:
                res += 1
            prev = sign
    return res


print(Solution().wiggleMaxLength(nums=[1, 17, 5, 10, 13, 15, 10, 5, 16, 8]))
print(Solution().wiggleMaxLength(nums=[1, 2, 3, 4, 5, 6, 7, 8, 9]))
