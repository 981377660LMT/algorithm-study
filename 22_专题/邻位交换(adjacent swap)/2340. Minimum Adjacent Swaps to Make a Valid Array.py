# 求邻位交换的最小次数
# 使得数组最大值在数组最右边  最小值在数组最左边
from typing import List
from minAdjacentSwap import minAdjacentSwap1


class Solution:
    def minimumSwaps(self, nums: List[int]) -> int:
        """直接把目标数组求出来 再求逆序对的个数"""
        n = len(nums)
        min_, max_ = min(nums), max(nums)
        i1, i2 = nums.index(min_), next((i for i in range(n - 1, -1, -1) if nums[i] == max_))
        i1, i2 = sorted((i1, i2))
        target = [min_] + nums[:i1] + nums[i1 + 1 : i2] + nums[i2 + 1 :] + [max_]  # 数组的删除用切片
        return minAdjacentSwap1(nums, target)

    def minimumSwaps2(self, nums: List[int]) -> int:
        """只要计算最近的位置交换"""
        n = len(nums)
        min_, max_ = min(nums), max(nums)
        mini, maxi = nums.index(min_), next((i for i in range(n - 1, -1, -1) if nums[i] == max_))
        res = mini + (n - 1 - maxi)
        # !注意mini > maxi的情况 有一次交换可以两全其美 即最大值右移一次 最小值左移一次
        return res - 1 if mini > maxi else res


print(Solution().minimumSwaps(nums=[3, 4, 5, 5, 3, 1]))
print(Solution().minimumSwaps(nums=[2, 1]))
print(Solution().minimumSwaps(nums=[35, 25, 30, 25, 31, 39, 35]))
print(Solution().minimumSwaps(nums=[9]))
print(Solution().minimumSwaps2(nums=[3, 4, 5, 5, 3, 1]))
print(Solution().minimumSwaps2(nums=[2, 1]))
print(Solution().minimumSwaps2(nums=[35, 25, 30, 25, 31, 39, 35]))
print(Solution().minimumSwaps2(nums=[9]))
