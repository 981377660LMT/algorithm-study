from bisect import bisect_right


class Solution:
    def solve(self, nums):
        # 这里使用的是直径，因此最终返回需要除以 2
        def check(diameter):
            start = nums[0]
            end = start + diameter
            for _ in range(LIGHTS):
                idx = bisect_right(nums, end)
                if idx >= N:
                    return True
                start = nums[idx]
                end = start + diameter
            return False

        nums.sort()
        N = len(nums)
        if N <= 3:
            return 0
        LIGHTS = 3

        l, r = 0, nums[-1] - nums[0]
        while l <= r:
            mid = (l + r) // 2
            if check(mid):
                r = mid - 1
            else:
                l = mid + 1
        return l / 2
