from typing import List

# 对于每个整数 A[i]，可以选择 x = -K 或是 x = K （K 总是非负整数），并将 x 加到 A[i] 中。
# 在此过程之后，得到数组 B。 返回 B 的最大值和 B 的最小值之间可能存在的最小差值。
# 1 <= A.length <= 10000

# -k+k 等价于每个数+0 或 +2k
# 关键一点就是，排序后，在某一个节点，之前的数全部加2k


class Solution:
    def smallestRangeII(self, nums: List[int], k: int) -> int:
        nums = sorted(nums)
        res = nums[-1] - nums[0]

        for i in range(1, len(nums)):
            minVal = min(nums[0] + 2 * k, nums[i])
            maxVal = max(nums[i - 1] + 2 * k, nums[-1])
            res = min(res, maxVal - minVal)

        return res


print(Solution().smallestRangeII([1, 3, 6], 3))
# 输出：3
# 解释：B = [4,6,3]
