from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个 二维 整数数组 nums1 和 nums2.

# nums1[i] = [idi, vali] 表示编号为 idi 的数字对应的值等于 vali 。
# nums2[i] = [idi, vali] 表示编号为 idi 的数字对应的值等于 vali 。
# 每个数组都包含 互不相同 的 id ，并按 id 以 递增 顺序排列。

# 请你将两个数组合并为一个按 id 以递增顺序排列的数组，并符合下述条件：

# 只有在两个数组中至少出现过一次的 id 才能包含在结果数组内。
# 每个 id 在结果数组中 只能出现一次 ，并且其对应的值等于两个数组中该 id 所对应的值求和。如果某个数组中不存在该 id ，则认为其对应的值等于 0 。
# 返回结果数组。返回的数组需要按 id 以递增顺序排列。
class Solution:
    def mergeArrays(self, nums1: List[List[int]], nums2: List[List[int]]) -> List[List[int]]:
        mp1 = {x[0]: x[1] for x in nums1}
        mp2 = {x[0]: x[1] for x in nums2}
        res = []
        for k in sorted(set(mp1.keys()) | set(mp2.keys())):
            res.append([k, mp1.get(k, 0) + mp2.get(k, 0)])
        return res
