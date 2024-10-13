from itertools import accumulate
from typing import List


class SubArraySumManager:
    def __init__(self, nums: List[int]) -> None:
        self.nums = nums
        self.p1 = list(accumulate(nums, initial=0))
        self.p2 = list(accumulate(self.p1, initial=0))
        self.rp1 = (list(accumulate(nums[::-1], initial=0)))[::-1]
        self.rp2 = (list(accumulate(self.rp1[::-1], initial=0)))[::-1]

        self.pi1 = None  # i*nums[i] 前缀和
        self.pi2 = None  # i*i*nums[i] 前缀和

    def querySubArraySum(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有子数组的和

        a[i] 被计算 (i-L+1)*(R-i+1) 次
        [L,R]里所有子数组和为 ∑(L,R) nums[i]*(-i^2+(R+L)*i-LR+R-L+1)
        预处理出 nums[i]、i*nums[i]、i*i*nums[i] 的前缀和即可
        """
        if self.pi1 is None:
            self.pi1 = [0]
            for i, num in enumerate(self.nums):
                self.pi1.append(self.pi1[-1] + num * i)
        if self.pi2 is None:
            self.pi2 = [0]
            for i, num in enumerate(self.nums):
                self.pi2.append(self.pi2[-1] + num * i * i)

        sum1 = self.pi2[left] - self.pi2[right + 1]
        sum2 = (left + right) * (self.pi1[right + 1] - self.pi1[left])
        sum3 = (right + 1 - left - left * right) * (self.p1[right + 1] - self.p1[left])
        return sum1 + sum2 + sum3

    def querySubArraySumStartsAt(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有以nums[left]开头的子数组的和"""
        assert 0 <= left <= right <= len(self.nums) - 1
        sum1 = self.rp1[left] * (right - left + 1)
        sum2 = self.rp2[left + 1] - self.rp2[right + 2]
        return sum1 - sum2

    def querySubArraySumEndsAt(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有以nums[right]结尾的子数组的和

        (p1[right+1]-p1[left])+(p1[right+1]-p1[left+1])+...+(p1[right+1]-p1[right]]) 即
        p1[right+1]*(right-left+1)+p2[right+1]-p2[left]
        """
        assert 0 <= left <= right <= len(self.nums) - 1
        sum1 = self.p1[right + 1] * (right - left + 1)
        sum2 = self.p2[right + 1] - self.p2[left]
        return sum1 - sum2

    def querySubArraySumInclude(self, left: int, right: int, include: int) -> int:
        """O(1)查询[left,right]闭区间内所有包含include下标的子数组的和"""
        assert 0 <= left <= include <= right <= len(self.nums) - 1
        sum1 = (self.p2[right + 2] - self.p2[include + 1]) * (include - left + 1)
        sum2 = (self.p2[include + 1] - self.p2[left]) * (right - include + 1)
        return sum1 - sum2

    def querySubArrayOccurrence(self, left: int, right: int, i: int) -> int:
        """O(1)查询[left,right]闭区间内第i个元素在区间内所有子数组中出现的次数"""
        assert 0 <= left <= right <= len(self.nums) - 1
        return (right - i + 1) * (i - left + 1)
