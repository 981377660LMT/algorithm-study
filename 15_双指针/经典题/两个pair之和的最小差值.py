# n ≤ 100
# 两个pair之和的最小差值
class Solution:
    def solve(self, nums):
        n = len(nums)
        nums.sort()

        res = int(1e20)
        # 两个指针固定首尾，中间两个移动
        for i in range(n - 3):
            for j in range(i + 3, n):
                target = nums[i] + nums[j]
                start = i + 1
                end = j - 1
                while start < end:
                    curSum = nums[start] + nums[end]
                    res = min(res, abs(curSum - target))
                    if curSum < target:
                        start += 1
                    else:
                        end -= 1
        return res
