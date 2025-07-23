# 3624. 位计数深度为 K 的整数数目 II
#
# 给你一个整数数组 nums。

# 对于任意正整数 x，定义以下序列：

# p0 = x
# pi+1 = popcount(pi)，对于所有 i >= 0，其中 popcount(y) 表示整数 y 的二进制表示中 1 的个数。
# 这个序列最终会收敛到值 1。

# popcount-depth（位计数深度）定义为满足 pd = 1 的最小整数 d >= 0。

# 例如，当 x = 7（二进制表示为 "111"）时，该序列为：7 → 3 → 2 → 1，因此 7 的 popcount-depth 为 3。

# 此外，给定一个二维整数数组 queries，其中每个 queries[i] 可以是以下两种类型之一：

# [1, l, r, k] - 计算在区间 [l, r] 中，满足 nums[j] 的 popcount-depth 等于 k 的索引 j 的数量。
# [2, idx, val] - 将 nums[idx] 更新为 val。
# 返回一个整数数组 answer，其中 answer[i] 表示第 i 个类型为 [1, l, r, k] 的查询的结果。
#
# !0<=k<=5

from typing import List


class Solution:
    def popcountDepth(self, nums: List[int], queries: List[List[int]]) -> List[int]: ...
