# 此题是2035. 将数组分成两个数组并最小化数组和的差的进阶版

# 1 <= nums.length <= 30
from math import gcd
from typing import List, Set


# 我们要将 nums 数组中的每个元素移动到 A 数组 或者 B 数组中，使得 A 数组和 B 数组不为空，并且 average(A) == average(B) 。
# 如果可以完成则返回true ， 否则返回 false  。

# 1.将原数组每个元素减去所有元素的平均值
# 问题变成：找出若干个元素组成集合 B，这些元素的和为 0。

# 2. 直接枚举不行，需要折半枚举，求出两半所有子序列和，看1.两半能否和为0或者
# 2.每一半能够成立


def subsetSum(nums: List[int]) -> Set[int]:
    """求所有非空子集的可能和"""
    dp = set()
    for cur in nums:
        dp |= {(cur + pre) for pre in (dp | {0})}
    return dp


class Solution:
    def splitArraySameAverage(self, nums: List[int]) -> bool:
        """每个数减去平均数后，变成nums中是否存在和为0的非空真子集
        
        暴力枚举超时
        折半枚举
        在左边选出x 在右边选出-x 这样就找到了一个子集
        注意不能全选和全不选
        """
        if len(nums) <= 1:
            return False

        # 注意平均数可能为浮点数，为了精确，可以将 nums 的数都先乘以一个合适的 mul 使得平均数为整数
        n, sum_ = len(nums), sum(nums)
        mul = n // gcd(sum_, n)
        nums = [x * mul - mul * sum_ // n for x in nums]

        sums1, sums2 = subsetSum(nums[: n // 2]), subsetSum(nums[n // 2 :])
        if 0 in sums1 or 0 in sums2:
            return True

        # !排除0后,左右两方都不贡献元素或者都全部贡献元素的情况需要剔除，
        # 因为这样的分类会导致某一个集合没有元素，违背题目意思。
        # 最后将两个数组排序，利用双指针判断是否存在答案。
        sums1.discard(sum(nums[: n // 2]))
        sums2.discard(sum(nums[n // 2 :]))
        nums1, nums2 = sorted(sums1), sorted(sums2)
        i, j = 0, len(nums2) - 1
        while i < len(nums1) and j >= 0:
            if nums1[i] + nums2[j] == 0:
                return True
            elif nums1[i] + nums2[j] < 0:
                i += 1
            else:
                j -= 1
        return False


print(Solution().splitArraySameAverage([3, 1]))
print(Solution().splitArraySameAverage(nums=[1, 2, 3, 4, 5, 6, 7, 8]))

