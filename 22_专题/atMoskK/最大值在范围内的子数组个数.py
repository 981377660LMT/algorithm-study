# 0 ≤ n ≤ 100,000

# 子数组的全部元素都不大于 k的子数组个数
# atMostK
class Solution:
    def solve(self, nums, lo, hi):
        def atMostK(k: int) -> int:
            """子数组的最大值不大于k的子数组个数"""
            res = 0
            pre = 0
            for num in nums:
                if num <= k:
                    pre += 1
                else:
                    pre = 0
                res += pre
            return res

        return (atMostK(hi) - atMostK(lo - 1)) % int(1e9 + 7)


print(Solution().solve(nums=[1, 5, 3, 2], lo=1, hi=4))
# We have the following sublists where 1 ≤ max(A) ≤ 4

# [1]
# [3]
# [3, 2]
# [2]
