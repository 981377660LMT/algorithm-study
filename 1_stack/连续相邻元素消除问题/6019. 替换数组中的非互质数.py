from math import gcd
from typing import List, Tuple
from collections import defaultdict

# 1 <= nums.length <= 105
# 1 <= nums[i] <= 105

# 从 nums 中找出 任意 两个 相邻 的 非互质 数。
# 如果不存在这样的数，终止 这一过程。
# 否则，删除这两个数，并 替换 为它们的 最小公倍数
# 只要还能找出两个相邻的非互质数就继续 重复 这一过程。


# 时间复杂度：O(n⋅logM)，其中 n 为 nums 的长度，M 为 nums 的最大值（求 gcd 需要O(logM)）。


class Solution:
    def replaceNonCoprimes(self, nums: List[int]) -> List[int]:
        if len(nums) == 1:
            return nums

        res = []
        for num in nums:
            res.append(num)
            while len(res) >= 2 and gcd(res[-2], res[-1]) != 1:
                n1, n2 = res.pop(), res.pop()
                res.append(n1 * n2 // gcd(n1, n2))

        return res


print(Solution().replaceNonCoprimes(nums=[6, 4, 3, 2, 7, 6, 2]))
print(Solution().replaceNonCoprimes([287, 41, 49, 287, 899, 23, 23, 20677, 5, 825]))
# 预期：
# [2009,20677,825]
