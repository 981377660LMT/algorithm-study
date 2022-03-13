from typing import List, Optional, Tuple

MOD = int(1e9 + 7)

# 思路： 贪心取出前 k - 1 个元素，第 k 次操作可以
# 将已经移除的最大的元素放回栈顶，或
# 取出第 k 个元素，此时栈顶为第 k + 1 个元素

# 要分类讨论的题，先把边界去除


class Solution:
    def maximumTop(self, nums: List[int], k: int) -> int:
        n = len(nums)

        if n == 1:
            return nums[0] if k % 2 == 0 else -1
        if k == 0:
            return nums[0]
        if k == 1:
            return nums[1]

        # 取出前 k - 1个元素，放回已取出的最大数或将第 k 个元素取出
        return max(max(nums[: k - 1]), (nums[k] if k < n else -1))
        if k > n:
            return max(nums)
        elif k < n:
            # 取出前 k - 1个元素，放回已取出的最大数或将第 k 个元素取出
            return max(max(nums[: k - 1]), (nums[k] if k < n else -1))
        else:
            # 最后一个不能去
            return max(nums[:-1])


print(Solution().maximumTop(nums=[5, 2, 2, 4, 0, 6], k=4))
