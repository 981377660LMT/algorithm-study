#  给你一个 下标从 0 开始 的整数数组 nums ，返回满足下述条件的 不同 四元组 (a, b, c, d) 的 数目 ：
#  nums[a] + nums[b] + nums[c] == nums[d] ，且
#  a < b < c < d

from typing import List
from itertools import combinations
from collections import defaultdict

# 用3个哈希表分别存储1个数及其个数，
# 2个数之和及其个数，3个数之和及其个数，
# 通过循环迭代，
# 把同时在第3个哈希表和原数组中出现的值对应的个数相加即为所有满足题意的个数

# O(3)
# 等式可转换成nums[a] + nums[b] == nums[d] - nums[c]，先用哈希表存储一边的结果，再遍历计算另一边的结果是否存在于哈希表。
class Solution:
    def countQuadruplets1(self, nums: List[int]) -> int:
        return sum(a + b + c == d for a, b, c, d in combinations(nums, 4))

    def countQuadruplets(self, nums: List[int]) -> int:
        n = len(nums)
        counter = defaultdict(list)
        for i in range(n - 3):
            for j in range(i + 1, n - 2):
                # 表示最大取j
                counter[nums[i] + nums[j]].append(j)

        res = 0
        for i in range(2, n):
            for j in range(i + 1, n):
                lis = counter[nums[j] - nums[i]]
                if not lis:
                    continue
                # i需要大于之前的最大值
                res += sum(i > limit for limit in lis)

        return res


print(Solution().countQuadruplets([1, 2, 3, 6]))
