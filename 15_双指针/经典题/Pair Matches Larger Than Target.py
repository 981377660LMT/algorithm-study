# 最多有几个pair, 使得abs(nums[i] - nums[j]) ≥ target
# i,j不可重复
# 差值大于target的二元组个数


class Solution:
    def solve(self, nums, target):
        def check(mid):
            return all(nums[-mid + i] - nums[i] >= target for i in range(mid))

        nums.sort()
        left, right = 0, len(nums) // 2
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right

    def solve2(self, nums, target):
        """双指针,因为单调性;对每一个left,找最近的right"""

        n = len(nums)
        nums.sort()
        res = 0

        right = n // 2
        for left in range(n // 2):
            while right < n and nums[right] - nums[left] < target:
                right += 1
            if right < n:
                res += 1
                right += 1
        return res


print(Solution().solve([1, 3, 5, 9, 10], 3))

