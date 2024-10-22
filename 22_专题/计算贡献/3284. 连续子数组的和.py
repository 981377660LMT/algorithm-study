# 3284. 连续子数组的和
# https://leetcode.cn/problems/sum-of-consecutive-subarrays/description/
# 如果一个长度为 n 的数组 arr 符合下面其中一个条件，可以称它 连续：
#
# 对于所有的 1 <= i < n，arr[i] - arr[i - 1] == 1。
# 对于所有的 1 <= i < n，arr[i] - arr[i - 1] == -1。
# 给定一个整数数组 nums，返回所有 连续 子数组的 值 之和。
# 由于答案可能很大，返回它对 1e9 + 7 取模 的值。
# 注意 长度为 1 的数组也被认为是连续的。


from itertools import chain
from typing import Generator, List, Tuple


MOD = int(1e9 + 7)


def enumetateConsecutiveSubarrays(
    nums: List[int], diff=1
) -> Generator[Tuple[int, int], None, None]:
    """遍历连续的子数组.

    >>> list(enumetateConsecutiveSubarrays([1, 2, 3, 5, 6, 7, 9]))
    [(0, 3), (3, 6), (6, 7)]
    """
    i, n = 0, len(nums)
    while i < n:
        start = i
        while i < n - 1 and nums[i] + diff == nums[i + 1]:
            i += 1
        i += 1
        yield start, i


class Solution:
    def getSum(self, nums: List[int]) -> int:
        group1, group2 = [], []
        for start, end in enumetateConsecutiveSubarrays(nums, diff=1):
            if end - start > 1:
                group1.append((start, end))
        for start, end in enumetateConsecutiveSubarrays(nums, diff=-1):
            if end - start > 1:
                group2.append((start, end))

        res = 0
        for start, end in chain(group1, group2):
            for i in range(start, end):
                v = nums[i]
                # !每个元素在子数组中出现的次数，注意减去自身
                c = (end - i) * (i - start + 1) - 1
                res += v * c
                res %= MOD

        res += sum(nums)
        return res % MOD


print(Solution().getSum([7, 6, 1, 2]))  # [(1, 3), (5, 7)] [(3, 1), (7, 5)]
# [1,2,3]
print(Solution().getSum([1, 2, 3]))  # 10

print(Solution().getSum([1, 3, 5, 7]))  # 10
# [51,99,100,99,81,87,83,47,46,91,15]
print(Solution().getSum([51, 99, 100, 99, 81, 87, 83, 47, 46, 91, 15]))  # 1290
