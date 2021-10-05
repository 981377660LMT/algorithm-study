# 给你一个数组 nums ，数组中有 2n 个元素，按 [x1,x2,...,xn,y1,y2,...,yn] 的格式排列。

# 请你将数组按 [x1,y1,x2,y2,...,xn,yn] 格式重新排列


class Solution:
    def shuffle(self, nums, n):
        nums[::2], nums[1::2] = nums[:n], nums[n:]
        return nums

