from typing import List

# 每次大数都变成比自己小一号的数字；逐步降级，每大一次后面的层数就加一(nestedInteger那题?)
# 最后相当于计算各个数字的层数之和
class Solution:
    def reductionOperations(self, nums: List[int]) -> int:
        nums = sorted(nums)

        res = 0
        depth = 0
        for i in range(1, len(nums)):
            if nums[i - 1] != nums[i]:
                depth += 1
            res += depth

        return res


print(Solution().reductionOperations(nums=[5, 1, 3]))
# 输出：3
# 解释：需要 3 次操作使 nums 中的所有元素相等：
# 1. largest = 5 下标为 0 。nextLargest = 3 。将 nums[0] 减少到 3 。nums = [3,1,3] 。
# 2. largest = 3 下标为 0 。nextLargest = 1 。将 nums[0] 减少到 1 。nums = [1,1,3] 。
# 3. largest = 3 下标为 2 。nextLargest = 1 。将 nums[2] 减少到 1 。nums = [1,1,1] 。
