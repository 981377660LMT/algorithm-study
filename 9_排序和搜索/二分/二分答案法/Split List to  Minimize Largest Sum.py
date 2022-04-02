# 分成k个子数组 并最小化最大和
# 求出最小的最大值
class Solution:
    def solve(self, nums, k):
        def check(mid):
            count = 0
            i = 0
            while i < n:
                tempX = mid
                while i < n and tempX - nums[i] >= 0:
                    tempX -= nums[i]
                    i += 1
                count += 1
            return count <= k

        left = max(nums)
        right = sum(nums)
        n = len(nums)

        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1

        return left
