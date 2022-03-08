# 此题是2035. 将数组分成两个数组并最小化数组和的差的进阶版

# 1 <= nums.length <= 30
from typing import List


# 我们要将 nums 数组中的每个元素移动到 A 数组 或者 B 数组中，使得 A 数组和 B 数组不为空，并且 average(A) == average(B) 。
# 如果可以完成则返回true ， 否则返回 false  。

# 1.将原数组每个元素减去所有元素的平均值
# 问题变成：找出若干个元素组成集合 B，这些元素的和为 0。

# 2. 直接枚举不行，需要折半枚举，求出两半所有子序列和，看1.两半能否和为0或者2.每一半能够成立

EPS = 1e-5


class Solution:
    def splitArraySameAverage(self, nums: List[int]) -> bool:
        def getSum(arr: List[int]) -> List[int]:
            """返回子序列的和"""

            def dfs(index: int, curSum: int, state: int) -> None:
                if index == len(arr):
                    res[state] = curSum
                    return
                dfs(index + 1, curSum, state)
                dfs(index + 1, curSum + arr[index], state | (1 << index))

            n = len(arr)
            res = [0] * (1 << n)
            dfs(0, 0, 0)
            return res

        n = len(nums)
        if n == 1:
            return False
        avg = sum(nums) / n
        nums = [num - avg for num in nums]

        half = n // 2
        s1, s2 = getSum(nums[:half]), getSum(nums[half:])

        # 先判断是否能独自构成答案，注意去掉全不取的情况
        for num in s1[1:]:
            if abs(num) < EPS:
                return True
        for num in s2[1:]:
            if abs(num) < EPS:
                return True

        # 左右两方都不贡献元素或者都全部贡献元素的情况需要剔除，因为这样的分类会导致某一个集合没有元素，违背题目意思
        s1[0], s1[-1] = int(1e20), int(1e20)
        s2[0], s2[-1] = int(1e20), int(1e20)
        s1, s2 = sorted(s1), sorted(s2)

        # 然后双指针看两个排序数组和是否和为0
        i, j = 0, len(s2) - 1
        while i < len(s1) and j >= 0:
            curSum = s1[i] + s2[j]
            if abs(curSum) < EPS:
                return True
            if curSum > 0:
                j -= 1
            else:
                i += 1

        return False


print(Solution().splitArraySameAverage([3, 1]))
print(Solution().splitArraySameAverage(nums=[1, 2, 3, 4, 5, 6, 7, 8]))

