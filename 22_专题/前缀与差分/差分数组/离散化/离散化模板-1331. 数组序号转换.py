# https://leetcode.cn/problems/rank-transform-of-an-array/

"""给你一个整数数组 arr ，请你将数组中的每个元素替换为它们排序后的序号。"""
from bisect import bisect_right
from typing import List


class Solution:
    def arrayRankTransform1(self, arr: List[int]) -> List[int]:
        """离散+二分查询"""
        sl = sorted(set(arr))
        res = []
        for num in arr:
            pos = bisect_right(sl, num)
            res.append(pos)
        return res

    def arrayRankTransform2(self, arr: List[int]) -> List[int]:
        """离散+哈希表查询"""
        sl = sorted(set(arr))
        mapping = {sl[i]: i + 1 for i in range(len(sl))}
        return [mapping[num] for num in arr]
