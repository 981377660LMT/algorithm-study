# 前缀和的前缀和 是求区间里所有子数组的和
from functools import lru_cache
from itertools import accumulate
from typing import List


class SubArraySumManager:
    def __init__(self, nums: List[int]) -> None:
        self.nums = nums
        self.p1 = [0] + list(accumulate(nums))
        self.p2 = [0] + list(accumulate(self.p1))
        self.rp1 = ([0] + list(accumulate(nums[::-1])))[::-1]
        self.rp2 = ([0] + list(accumulate(self.rp1[::-1])))[::-1]
        print(self.p1, self.p2, self.rp1, self.rp2)

    def querySubArraySum(self, left: int, right: int) -> int:
        """O(logn)查询[left,right]闭区间内所有子数组的和

        O(logn)的解法
        为记mid为[L,R]中点，可以把[L,R]内的子数组划分为包含mid的和不包含mid的
        包含mid的可以由`querySubArraySumInclude`O(1)求出，
        不包含mid的可以分治到[L,mid-1]和[mid+1,R]两个子区间上递归求解

        #todo 有没有O(1)的方法呢
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

