from typing import List
from bisect import bisect_right

# 请你统计并返回 nums 中能满足其最小元素与最大元素的 和 小于或等于 target 的 非空 子序列的数目。

# 排序+双指针查找+计算贡献
# 两数之和

M = int(1e9 + 7)


class Solution:
    def numSubseq(self, nums: List[int], target: int) -> int:
        nums = sorted(nums)
        res, left, right = 0, 0, len(nums) - 1
        while left <= right:
            if nums[left] + nums[right] > target:
                right -= 1
            else:
                res += pow(2, right - left, M)
                left += 1
        return res % M


print(Solution().numSubseq(nums=[3, 5, 6, 7], target=9))
# 输出：4
# 解释：有 4 个子序列满足该条件。
# [3] -> 最小元素 + 最大元素 <= target (3 + 3 <= 9)
# [3,5] -> (3 + 5 <= 9)
# [3,5,6] -> (3 + 6 <= 9)
# [3,6] -> (3 + 6 <= 9)
