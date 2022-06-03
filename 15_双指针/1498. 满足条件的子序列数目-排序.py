from typing import List
from bisect import bisect_right

# 请你统计并返回 nums 中能满足其最小元素与最大元素的 和 小于或等于 target 的 非空 子序列的数目。

# 排序+双指针查找+计算贡献
# 两数之和

MOD = int(1e9 + 7)


class Solution:
    def numSubseq(self, nums: List[int], target: int) -> int:
        """
        子序列：排序
        双指针查找+计算贡献
        固定取左端点 看右端点可以伸到哪里
        """
        nums = sorted(nums)
        res, right = 0, len(nums) - 1
        for left in range(len(nums)):
            while right >= 0 and nums[right] + nums[left] > target:
                right -= 1
            res += pow(2, right - left, MOD) if right >= left else 0
            res %= MOD
        return res


print(Solution().numSubseq(nums=[3, 5, 6, 7], target=9))
# 输出：4
# 解释：有 4 个子序列满足该条件。
# [3] -> 最小元素 + 最大元素 <= target (3 + 3 <= 9)
# [3,5] -> (3 + 5 <= 9)
# [3,5,6] -> (3 + 6 <= 9)
# [3,6] -> (3 + 6 <= 9)
