# 拆成两半 左边和等于右边 且左边所有数小于右边
class Solution:
    def solve(self, nums):
        sum_ = sum(nums)
        if sum_ & 1:
            return False

        target = sum_ // 2
        nums.sort(reverse=True)

        index = 0
        while target - nums[index] >= 0:
            target -= nums[index]
            index += 1

        return not (nums[index - 1] == nums[index]) and target == 0

