# 3299. 连续子序列的和
# https://leetcode.cn/problems/sum-of-consecutive-subsequences/description/
# 如果一个长度为 n 的数组 arr 符合下面其中一个条件，可以称它 连续：
# 对于所有的 1 <= i < n，arr[i] - arr[i - 1] == 1。
# 对于所有的 1 <= i < n，arr[i] - arr[i - 1] == -1。
# 数组的 值 是其元素的和。
# 例如，[3, 4, 5] 是一个值为 12 的连续数组，并且 [9, 8] 是另一个值为 17 的连续数组。而 [3, 4, 3] 和 [8, 6] 都不连续。
# 给定一个整数数组 nums，返回所有 连续 非空子序列的 值 之和。
# 由于答案可能很大，返回它对 1e9 + 7 取模 的值。
# 注意 长度为 1 的数组也被认为是连续的。
#
# !维护四个字典, 分别表示: 以 v 结尾的递(增/减)子序列(数量/和)


from collections import defaultdict
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def getSum(self, nums: List[int]) -> int:
        c1, s1 = defaultdict(int), defaultdict(int)  # increasing
        c2, s2 = defaultdict(int), defaultdict(int)  # decreasing
        for v in nums:
            s1[v] = (s1[v] + s1[v - 1] + v * (c1[v - 1] + 1)) % MOD
            c1[v] = (c1[v] + c1[v - 1] + 1) % MOD
            s2[v] = (s2[v] + s2[v + 1] + v * (c2[v + 1] + 1)) % MOD
            c2[v] = (c2[v] + c2[v + 1] + 1) % MOD
        return (sum(s1.values()) + sum(s2.values()) - sum(nums)) % MOD
