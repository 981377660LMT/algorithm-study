from typing import List

# 你可以设计一个对数时间复杂度（即，O(log(n))）的解决方案吗？

# 现有一个按 `升序` 排列的整数数组 nums ，其中每个数字都 `互不相同` 。
# 给你一个整数 k ，请你找出并返回从数组最左边开始的第 k 个缺失数字。


class Solution:
    def missingElement2(self, nums: List[int], k: int) -> int:
        n = len(nums)

        for i in range(1, n):
            loss = nums[i] - nums[i - 1] - 1
            if loss >= k:
                return nums[i - 1] + k
            k -= loss

        return nums[n - 1] + k

    # 在idx时缺失的元素数目为 nums[idx]-nums[0]-idx
    # 最右二分
    def missingElement(self, nums: List[int], k: int) -> int:
        l, r = 0, len(nums) - 1
        while l <= r:
            mid = (l + r) >> 1
            missed = nums[mid] - nums[0] - mid
            if missed >= k:
                r = mid - 1
            else:
                l = mid + 1
        print(nums[l], nums[r], l, r)

        # 疑惑
        return nums[0] + k + r


print(Solution().missingElement(nums=[4, 7, 9, 10], k=3))
# 解释：缺失数字有 [5,6,8,...]，因此第三个缺失数字为 8 。
