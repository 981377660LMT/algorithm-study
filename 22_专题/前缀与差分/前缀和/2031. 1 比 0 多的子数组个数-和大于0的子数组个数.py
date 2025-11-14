# 1 <= nums.length <= 105
# 1比0多的子数组个数/和大于0的子数组个数/和为正数的子数组数

from typing import List
from collections import defaultdict
from sortedcontainers import SortedList

MOD = int(1e9 + 7)


def subarraysWithSumMoreThanZero(nums: List[int]) -> int:
    """和大于0的子数组个数"""
    res, curSum, sl = 0, 0, SortedList([0])
    for num in nums:
        curSum += num
        res += sl.bisect_left(curSum)
        sl.add(curSum)
    return res


class Solution:
    def subarraysWithMoreZerosThanOnes(self, nums: List[int]) -> int:
        nums = [1 if num == 1 else -1 for num in nums]
        return subarraysWithSumMoreThanZero(nums) % MOD

    def subarraysWithMoreZerosThanOnes2(self, nums: List[int]) -> int:
        """和为正数的子数组有多少个

        O(n) dp解法
        每次需要查询小于当前值的个数，但是查询值每次变化都是+1或者-1,
        所以可以使用一个额外的变量来记录,查询复杂度O(logn)变为O(1)
        """

        counter = defaultdict(
            int, {0: 1}
        )  # 以前一个元素结尾的前缀和为key的子数组个数(TODO: 换成长为2n+1的数组会更快)

        res, curSum, leftSmaller = 0, 0, 0
        for num in nums:
            if num == 1:
                curSum += 1
                leftSmaller += counter[curSum - 1]  # 之前前缀和等于curSum的子数组都可以累加
            else:
                curSum -= 1
                leftSmaller -= counter[curSum]  # 之前前缀和等于curSum的子数组不能要了

            counter[curSum] += 1
            res = (res + leftSmaller) % MOD

        return res


print(Solution().subarraysWithMoreZerosThanOnes2(nums=[0, 1, 1, 0, 1]))
# 输入: nums = [0,1,1,0,1]
# 输出: 9
# 解释:
# 长度为 1 的、1 的数量大于 0 的数量的子数组有: [1], [1], [1]
# 长度为 2 的、1 的数量大于 0 的数量的子数组有: [1,1]
# 长度为 3 的、1 的数量大于 0 的数量的子数组有: [0,1,1], [1,1,0], [1,0,1]
# 长度为 4 的、1 的数量大于 0 的数量的子数组有: [1,1,0,1]
# 长度为 5 的、1 的数量大于 0 的数量的子数组有: [0,1,1,0,1]
