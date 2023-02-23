# 给你两个下标从 0 开始的数组 nums1 和 nums2 ，和一个二维数组 queries 表示一些操作。总共有 3 种类型的操作：

# 操作类型 1 为 queries[i] = [1, l, r] 。你需要将 nums1 从下标 l 到下标 r 的所有 0 反转成 1 或将 1 反转成 0 。l 和 r 下标都从 0 开始。
# 操作类型 2 为 queries[i] = [2, p, 0] 。对于 0 <= i < n 中的所有下标，令 nums2[i] = nums2[i] + nums1[i] * p 。
# 操作类型 3 为 queries[i] = [3, 0, 0] 。求 nums2 中所有元素的和。
# 请你返回一个数组，包含所有第三种操作类型的答案。
from typing import List
from BitSetInt import BitSet


# bitset api: onesCount/flipRange
class Solution:
    def handleQuery(
        self, nums1: List[int], nums2: List[int], queries: List[List[int]]
    ) -> List[int]:
        bs = BitSet.fromlist(nums1)
        sum_ = sum(nums2)
        res = []
        for op, a, b in queries:
            if op == 1:
                bs.flip_range(a, b)
            elif op == 2:
                sum_ += a * bs.bit_count()
            else:
                res.append(sum_)
        return res
