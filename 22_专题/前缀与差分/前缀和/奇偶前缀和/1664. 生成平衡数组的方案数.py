# 如果一个数组满足奇数下标元素的和与偶数下标元素的和相等，该数组就是一个 平衡数组 。
# 你需要选择 恰好 一个下标（下标从 0 开始）并删除对应的元素。
# 请你返回删除操作后，剩下的数组 nums 是 平衡数组 的 方案数 。

# 1 <= nums.length <= 1e5

# !技巧:奇偶前缀和统计奇数和偶数前缀和的差值

from typing import List


class Solution:
    def waysToMakeFair(self, nums: List[int]) -> int:
        n = len(nums)
        preSum = [0] * (n + 1)  # 偶数下标之和-奇数下标之和 的前缀和
        for i in range(n):
            preSum[i + 1] = preSum[i] + (nums[i] if i % 2 == 0 else -nums[i])

        res = 0
        for i in range(n):
            left = preSum[i]  # 左边的差值
            right = (preSum[n] - preSum[i + 1]) * -1  # 右边的差值
            res += left == -right
        return res


print(Solution().waysToMakeFair(nums=[2, 1, 6, 4]))
# 输出：1
# 解释：
# 删除下标 0 ：[1,6,4] -> 偶数元素下标为：1 + 4 = 5 。奇数元素下标为：6 。不平衡。
# 删除下标 1 ：[2,6,4] -> 偶数元素下标为：2 + 4 = 6 。奇数元素下标为：6 。平衡。
# 删除下标 2 ：[2,1,4] -> 偶数元素下标为：2 + 4 = 6 。奇数元素下标为：1 。不平衡。
# 删除下标 3 ：[2,1,6] -> 偶数元素下标为：2 + 6 = 8 。奇数元素下标为：1 。不平衡。
# 只有一种让剩余数组成为平衡数组的方案。
