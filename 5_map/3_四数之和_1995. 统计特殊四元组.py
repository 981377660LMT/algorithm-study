#  给你一个 下标从 0 开始 的整数数组 nums ，返回满足下述条件的 不同 四元组 (a, b, c, d) 的 数目 ：

#  nums[a] + nums[b] + nums[c] == nums[d] ，且
#  a < b < c < d
from collections import defaultdict
from typing import List

# 用3个哈希表分别存储1个数及其个数，
# 2个数之和及其个数，3个数之和及其个数，
# 通过循环迭代，
# 把同时在第3个哈希表和原数组中出现的值对应的个数相加即为所有满足题意的个数


class Solution:
    def countQuadruplets(self, nums: List[int]) -> int:
        d1, d2, d3 = defaultdict(int), defaultdict(int), defaultdict(int)
        res = 0
        for i in nums:
            res += d3[i]
            for j in d2:
                d3[i + j] += d2[j]
            for j in d1:
                d2[i + j] += d1[j]
            d1[i] += 1
        print(d1, d2, d3, sep='\n')
        return res


print(Solution().countQuadruplets([1, 2, 3, 6]))
