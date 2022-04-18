from typing import List

# 如果一个数组满足奇数下标元素的和与偶数下标元素的和相等，该数组就是一个 平衡数组 。
# 请你返回删除操作后，剩下的数组 nums 是 平衡数组 的 方案数 。

# 1 <= nums.length <= 105


# 借助前缀和----->奇偶差的和


class Solution:
    def waysToMakeFair(self, nums: List[int]) -> int:
        n = len(nums)
        preSum = [0] * (n + 1)
        for i in range(n):
            preSum[i + 1] = preSum[i] + (nums[i] if not i & 1 else -nums[i])

        res = 0
        for i in range(n):
            leftDiff = preSum[i]
            rightDiff = (preSum[-1] - preSum[i + 1]) * -1
            if leftDiff + rightDiff == 0:
                res += 1

        return res


print(Solution().waysToMakeFair(nums=[2, 1, 6, 4]))
# 输出：1
# 解释：
# 删除下标 0 ：[1,6,4] -> 偶数元素下标为：1 + 4 = 5 。奇数元素下标为：6 。不平衡。
# 删除下标 1 ：[2,6,4] -> 偶数元素下标为：2 + 4 = 6 。奇数元素下标为：6 。平衡。
# 删除下标 2 ：[2,1,4] -> 偶数元素下标为：2 + 4 = 6 。奇数元素下标为：1 。不平衡。
# 删除下标 3 ：[2,1,6] -> 偶数元素下标为：2 + 6 = 8 。奇数元素下标为：1 。不平衡。
# 只有一种让剩余数组成为平衡数组的方案。

