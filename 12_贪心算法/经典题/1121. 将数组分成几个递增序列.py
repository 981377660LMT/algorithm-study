from typing import List
from collections import Counter

# 给你一个 非递减 的正整数数组 nums 和整数 K，判断该数组是否可以被分成一个或几个 长度至少 为 K 的 不相交的递增子序列。
class Solution:
    def canDivideIntoSubsequences2(self, nums: List[int], k: int) -> bool:
        # 能不能把一个数组分成一个或多个长度为K的不相交的递增子序列，仅取决于出现次数最多的那个元素的出现次数
        # 即每次分组时先从出现次数最多的元素取
        counter = Counter(nums)
        maxFreq = max(counter.values())
        return maxFreq * k <= len(nums)

    # 使用递增特性计数
    def canDivideIntoSubsequences(self, nums: List[int], k: int) -> bool:
        # 能不能把一个数组分成一个或多个长度为K的不相交的递增子序列，仅取决于出现次数最多的那个元素的出现次数
        # 即每次分组时先从出现次数最多的元素取
        pre = nums[0]
        maxDup = 0
        for num in nums:
            if num == pre:
                maxDup += 1
            else:
                pre = num
                maxDup = 1
            if maxDup * k > len(nums):
                return False
        return True


print(Solution().canDivideIntoSubsequences(nums=[1, 2, 2, 3, 3, 4, 4], k=3))
# 输出：true
# 解释：
# 该数组可以分成两个子序列 [1,2,3,4] 和 [2,3,4]，每个子序列的长度都至少是 3。
