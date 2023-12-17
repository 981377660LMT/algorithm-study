# 2964. Number of Divisible Triplet Sums.py
# 和模d为0的三元组个数
# https://leetcode.cn/problems/number-of-divisible-triplet-sums/description/
# 1 <= nums.length <= 1000
# 1 <= nums[i] <= 109
# 1 <= d <= 109
# !固定左端点，剩下的问题就是两数之和，O(n^2)


from collections import defaultdict
from typing import List


class Solution:
    def divisibleTripletCount(self, nums: List[int], d: int) -> int:
        counter = defaultdict(int)
        n, res = len(nums), 0
        for i in range(n):
            for j in range(i + 1, n):
                need = (-nums[i] - nums[j]) % d
                res += counter[need]
            counter[nums[i] % d] += 1
        return res
