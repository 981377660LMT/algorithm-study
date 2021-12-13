from typing import List

# nums 中的整数 互不相同
# 一次 `删除` 操作定义为从数组的 前面 移除一个元素或从数组的 后面 移除一个元素。
# 返回将数组中最小值和最大值 都 移除需要的最小`删除`次数。
class Solution:
    def minimumDeletions(self, nums: List[int]) -> int:
        i, j, n = nums.index(min(nums)), nums.index(max(nums)), len(nums)
        i, j = sorted((i, j))
        return min(j + 1, n - i, i + 1 + n - j)


print(Solution().minimumDeletions(nums=[2, 10, 7, 5, 4, 1, 8, 6]))
# 输出：5
# 解释：
# 数组中的最小元素是 nums[5] ，值为 1 。
# 数组中的最大元素是 nums[1] ，值为 10 。
# 将最大值和最小值都移除需要从数组前面移除 2 个元素，从数组后面移除 3 个元素。
# 结果是 2 + 3 = 5 ，这是所有可能情况中的最小删除次数。

