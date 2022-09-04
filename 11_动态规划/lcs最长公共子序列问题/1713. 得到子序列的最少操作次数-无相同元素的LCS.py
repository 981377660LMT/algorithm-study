# 1713. 得到子序列的最少操作次数
# 无相同元素的LCS可以转化为LIS


from typing import List
from bisect import bisect_left, bisect_right


def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    res = [nums[0]]
    for i in range(1, n):
        pos = bisect_left(res, nums[i]) if isStrict else bisect_right(res, nums[i])
        if pos >= len(res):
            res.append(nums[i])
        else:
            res[pos] = nums[i]

    return len(res)


class Solution:
    def minOperations(self, target: List[int], arr: List[int]) -> int:
        """
        给你一个数组 target,包含若干 互不相同 的整数,
        以及另一个整数数组 arr,arr 可能 包含重复元素。

        请你返回 最少 操作次数,使得 target 成为 arr 的一个子序列。
        """
        n = len(target)
        indexMap = {num: i for i, num in enumerate(target)}
        nums = [indexMap[num] for num in arr if num in indexMap]
        return n - LIS(nums, isStrict=True)


if __name__ == "__main__":
    print(Solution().minOperations(target=[5, 1, 3], arr=[9, 4, 2, 3, 4]))
