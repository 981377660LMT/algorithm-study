# 前缀和的前缀和 是求区间里所有子数组的和
from functools import lru_cache
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

    def __querySubArraySum(self, left: int, right: int) -> int:
        """O(logn)查询[left,right]闭区间内所有子数组的和

        O(logn)的解法
        为记mid为[L,R]中点，可以把[L,R]内的子数组划分为包含mid的和不包含mid的
        包含mid的可以由`querySubArraySumInclude`O(1)求出，
        不包含mid的可以分治到[L,mid-1]和[mid+1,R]两个子区间上递归求解

        #todo 有没有O(1)的方法呢
        a[i] 被计算 (i-L+1)*(R-i+1) 次
        [L,R]里所有子数组和为 ∑(L,R) nums[i]*(-i^2+(R+L)*i-LR+R-L+1)
        预处理出 nums[i]、i*nums[i]、i*i*nums[i] 的前缀和即可
        """

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left > right:
                return 0
            if left == right:
                return self.nums[left]
            mid = (left + right) >> 1
            return (
                dfs(left, mid - 1)
                + dfs(mid + 1, right)
                + self.querySubArraySumInclude(left, right, mid)
            )

        assert 0 <= left <= right <= len(self.nums) - 1
        return dfs(left, right)

    def querySubArraySumStartsAt(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有以nums[left]开头的子数组的和
        """
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
        """O(1)查询[left,right]闭区间内所有包含include下标的子数组的和
        """
        assert 0 <= left <= include <= right <= len(self.nums) - 1
        sum1 = (self.p2[right + 2] - self.p2[include + 1]) * (include - left + 1)
        sum2 = (self.p2[include + 1] - self.p2[left]) * (right - include + 1)
        return sum1 - sum2


if __name__ == '__main__':
    manager = SubArraySumManager([1, 3, 1, 2])
    assert manager.querySubArraySum(0, 2) == 18
    assert manager.querySubArraySumStartsAt(1, 3) == 13
    assert manager.querySubArraySumEndsAt(1, 3) == 11
    assert manager.querySubArraySumInclude(1, 3, 2) == 14

