from typing import List

# 一个序列的 宽度 定义为该序列中最大元素和最小元素的差值。
# 返回 nums 的所有非空 子序列 的 宽度之和 。
# 1 <= nums.length <= 10^5


# 对数组进行排序
# 统计具有最小值 A[i] 和最大值 A[j] 的子序列的数量
class Solution:
    def sumSubseqWidths(self, nums: List[int]) -> int:
        MOD = 1_000_000_007
        res = 0
        # 对每个数，统计他作为最大/最小值在多少个取他的子集里出现
        for i, x in enumerate(sorted(nums)):
            res += x * (pow(2, i, MOD) - pow(2, len(nums) - i - 1, MOD))
        return res % MOD


print(Solution().sumSubseqWidths(nums=[2, 1, 3]))
# 输出：6
# 解释：子序列为 [1], [2], [3], [2,1], [2,3], [1,3], [2,1,3] 。
# 相应的宽度是 0, 0, 0, 1, 1, 2, 2 。
# 宽度之和是 6 。
