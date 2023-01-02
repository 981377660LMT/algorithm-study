# https://leetcode.cn/problems/rank-transform-of-an-array/

"""给你一个整数数组 arr ，请你将数组中的每个元素替换为它们排序后的序号。"""
from bisect import bisect_right
from typing import List


class Solution:
    def arrayRankTransform1(self, arr: List[int]) -> List[int]:
        """离散+二分查询"""
        allNums = sorted(set(arr))
        res = []
        for num in arr:
            pos = bisect_right(allNums, num)
            res.append(pos)
        return res

    def arrayRankTransform2(self, arr: List[int]) -> List[int]:
        """离散+哈希表查询"""
        allNums = sorted(set(arr))
        mp = {allNums[i]: i + 1 for i in range(len(allNums))}
        return [mp[num] for num in arr]


# cpp 的离散化操作 (upper_bound是bisect_right)
# class Solution {
# public:
#     vector<int> arrayRankTransform(vector<int>& arr) {
#         vector<int> b = arr;
#         sort(b.begin(), b.end());
#         b.erase(unique(b.begin(), b.end()), b.end());
#         for(auto &x: arr)
#             x = upper_bound(b.begin(), b.end(), x) - b.begin();
#         return arr;
#     }
# };
