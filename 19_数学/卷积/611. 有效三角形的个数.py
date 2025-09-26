# 611. 有效三角形的个数
# https://leetcode.cn/problems/valid-triangle-number/solutions/2432875/zhuan-huan-cheng-abcyong-xiang-xiang-shu-1ex3/?envType=daily-question&envId=2025-09-26
# 给定一个包含非负整数的数组 nums ，返回其中可以组成三角形三条边的三元组个数。
#
# 1 <= nums.length <= 1000
# 0 <= nums[i] <= 1000
#
# 值域不大时可以fft

from math import comb
from typing import List
from scipy.signal import convolve


class Solution:
    def triangleNumber(self, nums: List[int]) -> int:
        max_ = max(nums)
        counter = [0] * (max_ + 1)
        for x in nums:
            counter[x] += 1
        c0 = counter[0]
        counter[0] = 0

        counter2 = convolve(counter, counter).round().tolist()

        pos = len(nums) - c0
        res = comb(pos, 3)
        s = 0  # (a,b) s.t. a + b <= c
        for c in range(1, max_ + 1):
            c2 = counter2[c]
            if c % 2 == 0:
                c2 -= counter[c // 2]
            s += c2 // 2
            res -= s * counter[c]
        return res
