from typing import List

# 2，对于所有 0 <= j < i 且 i < k <= nums.length - 1 ，满足 nums[j] < nums[i] < nums[k]
# 1，如果满足 nums[i - 1] < nums[i] < nums[i + 1] ，且不满足前面的条件
# 0，如果上述条件全部不满足


class Solution:
    def sumOfBeauties(self, nums: List[int]) -> int:
        n = len(nums)
        res = 0
        for i in range(1, n - 1):
            if nums[i - 1] < nums[i] < nums[i + 1]:
                res += 1

        # ----左侧最大值
        l_max = [0] * n
        l_max[0] = nums[0]
        for i in range(1, n):
            l_max[i] = max(l_max[i - 1], nums[i])

        # ----右侧最小值
        r_min = [0] * n
        r_min[n - 1] = nums[n - 1]
        for i in range(n - 2, -1, -1):
            r_min[i] = min(r_min[i + 1], nums[i])

        # ----判断是否符合2的情况
        for i in range(1, n - 1):
            if l_max[i - 1] < nums[i] < r_min[i + 1]:
                res += 1

        return res


print(Solution().sumOfBeauties(nums=[2, 4, 6, 4]))
# 输出：1
# 解释：对于每个符合范围 1 <= i <= 2 的下标 i :
# - nums[1] 的美丽值等于 1
# - nums[2] 的美丽值等于 0
