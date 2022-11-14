# 2035. 将数组分成两个数组并最小化数组和的差的进阶版

# 我们要将 nums 数组中的每个元素移动到 A 数组 或者 B 数组中，
# !使得 A 数组和 B 数组不为空，并且 average(A) == average(B) 。
# 如果可以完成则返回true ，否则返回 false  。
# !1 <= nums.length <= 30
# !0 <= nums[i] <= 1e4

# 1.将原数组每个元素减去所有元素的平均值
# !问题变成：找出若干个元素组成非空真子集 A，这些元素的和为 0 (从而剩下的元素组成集合 B，它们的和也同样为 0)。
# 2. 直接枚举不行，需要折半枚举，求出两半所有子序列和，看1.两半能否和为0或者
# 2.每一半能够成立


from typing import List, Set


def subsetSum(nums: List[int]) -> Set[int]:
    """O(2^n)求所有`非空子集`的可能和"""
    dp = set()
    for cur in nums:
        dp |= {(cur + pre) for pre in (dp | {0})}
    return dp


class Solution:
    def splitArraySameAverage(self, nums: List[int]) -> bool:
        """
        折半枚举 O(n*2^(n/2))

        # !每个数减去平均数后,变成nums中是否存在和为0的非空真子集
        1. 左边可以凑出0/右边可以凑出0
        2. 左边+右边可以凑出0
        """
        n = len(nums)
        if n <= 1:
            return False

        # 注意平均数可能为浮点数，为了精确，
        # 可以将 nums 的数都先乘以一个合适的倍数 ( n//gcd(sum(nums),n) 或者 n) 使得平均数为整数
        nums = [num * n for num in nums]
        avg = sum(nums) // n
        nums = [num - avg for num in nums]

        leftSums, rightSums = subsetSum(nums[: n // 2]), subsetSum(nums[n // 2 :])
        if 0 in leftSums or 0 in rightSums:
            return True

        # !左右两边选,注意左右两方全部贡献元素的情况需要剔除，
        # 因为这样的分类会导致某一个集合没有元素，违背题目意思。
        leftSums.discard(sum(nums[: n // 2]))
        rightSums.discard(sum(nums[n // 2 :]))

        # 最后将两个数组排序，利用双指针判断是否存在答案。
        nums1, nums2 = sorted(leftSums), sorted(rightSums)
        i, j = 0, len(nums2) - 1
        while i < len(nums1) and j >= 0:
            if nums1[i] + nums2[j] == 0:
                return True
            elif nums1[i] + nums2[j] < 0:
                i += 1
            else:
                j -= 1

        return False

    def splitArraySameAverage2(self, nums: List[int]) -> bool:
        """背包dp O(n^2 * sum(nums))

        #!集合A的平均值等于B的平均值 <=>
        #!`集合A的平均值`等于`nums的平均值` <=>
        #!在数组中选k个数,使得和为k*avg => 背包dp

        !剪枝1:划分成两个子数组,对长度短的那个dp
        !剪枝2:sum(A)*k = sum(nums)*n 要求 `sum(A)*k 模 n 为 0`
        """
        n = len(nums)
        if n <= 1:
            return False
        mid, sum_ = n // 2, sum(nums)
        if all(sum_ * k % n != 0 for k in range(1, mid + 1)):
            return False

        dp = [set() for _ in range(mid + 1)]
        dp[0].add(0)
        for num in nums:
            ndp = [s.copy() for s in dp]
            for i in range(mid):
                for pre in dp[i]:
                    curSum = pre + num
                    if curSum * n == sum_ * (i + 1):  # A和nums的平均值相等
                        return True
                    ndp[i + 1].add(curSum)
            dp = ndp
        return False


print(Solution().splitArraySameAverage([3, 1]))  # False
print(Solution().splitArraySameAverage(nums=[1, 2, 3, 4, 5, 6, 7, 8]))
