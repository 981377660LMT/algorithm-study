"""奇偶分组"""

from collections import Counter, defaultdict
from typing import List


# 给你两个正整数数组 nums 和 target ，两个数组长度相等。
# 在一次操作中，你可以选择两个 不同 的下标 i 和 j ，其中 0 <= i, j < nums.length ，并且：
# 令 nums[i] = nums[i] + 2 且
# 令 nums[j] = nums[j] - 2 。
# !如果两个数组中每个元素出现的频率相等，我们称两个数组是 相似 的。
# 请你返回将 nums 变得与 target 相似的最少操作次数。
# 测试数据保证 nums 一定能变得与 target 相似。

# !奇数只能变成奇数，偶数只能变成偶数。
# !分别考虑奇数数组和目标的奇数数组以及偶数数组和目标的偶数数组，一定是维持当前的顺序找目标更优。（否则总距离会更长）


class Solution:
    def makeSimilar(self, nums: List[int], target: List[int], k=2) -> int:
        if sum(nums) != sum(target):
            return -1
        if k == 0:
            return 0 if Counter(nums) == Counter(target) else -1

        g1, g2 = defaultdict(list), defaultdict(list)
        for i in range(len(nums)):
            g1[nums[i] % k].append(nums[i])
            g2[target[i] % k].append(target[i])
        for g in g1.values():
            g.sort()
        for g in g2.values():
            g.sort()

        res = 0
        for mod in range(k):
            if len(g1[mod]) != len(g2[mod]):
                return -1
            for a, b in zip(g1[mod], g2[mod]):
                res += abs(a - b) // k
        return res // 2


assert Solution().makeSimilar(nums=[8, 12, 6], target=[2, 14, 10]) == 2
assert Solution().makeSimilar(nums=[1, 1, 1, 1, 1], target=[1, 1, 1, 1, 1]) == 0
