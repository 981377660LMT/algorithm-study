# https://leetcode.cn/problems/greatest-common-divisor-traversal/
# 于两个下标 i 和 j（i != j），当且仅当 gcd(nums[i], nums[j]) > 1 时，
# 我们可以在两个下标之间通行，其中 gcd 是两个数的 最大公约数 。
# 你需要判断 nums 数组中 任意 两个满足 i < j 的下标 i 和 j ，
# 是否存在若干次通行可以从 i 遍历到 j 。

# 1 <= nums.length <= 105
# 1 <= nums[i] <= 105


# !注意特判全1的情形


from typing import List
from 埃氏筛和并查集 import EratosthenesSieve, UnionFindArray

E = EratosthenesSieve(int(1e5) + 10)


class Solution:
    def canTraverseAllPairs(self, nums: List[int]) -> bool:
        if len(nums) <= 1:
            return True
        if 1 in nums:
            return False

        max_ = max(nums, default=0)
        uf = UnionFindArray(max_ + 1)
        for num in nums:
            for p in E.getPrimeFactors(num):
                uf.union(num, p)

        roots = set(uf.find(num) for num in nums)
        return len(roots) == 1
