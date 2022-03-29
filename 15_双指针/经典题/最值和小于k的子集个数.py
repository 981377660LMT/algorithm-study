class Solution:
    def solve(self, nums, k):
        """
        return the number of non-empty subsets S such that min(S) + max(S) ≤ k
        固定最小值，看最大值能取多少
        """
        n = len(nums)
        nums.sort()

        res = 0
        right = n - 1
        for left in range(n):
            while right >= 0 and nums[right] + nums[left] > k:
                right -= 1
            if left <= right and nums[right] + nums[left] <= k:
                res += 2 ** (right - left)

        return res

