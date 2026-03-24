# 3845. 最大子数组异或值_区间最大异或
# https://leetcode.cn/problems/maximum-subarray-xor-with-bounded-range/description/
# 给你一个非负整数数组 nums 和一个整数 k。
# 你需要选择 nums 的一个 子数组，使得该子数组中元素的 最大值 与 最小值 之间的差值不超过 k。这个子数组的 值 定义为子数组中所有元素按位异或（XOR）的结果。
# 返回一个整数，表示所选子数组可能获得的 最大值 。
# 子数组 是数组中任意连续、非空 的元素序列。


from itertools import accumulate
from operator import xor

from sortedcontainers import SortedList

from XorTrie import BinaryTrie


class Solution:
    def maxXor(self, nums: list[int], k: int) -> int:
        prexor = list(accumulate(nums, xor, initial=0))
        trie = BinaryTrie(max(prexor), addLimit=len(nums) + 1, allowMultipleElements=True)
        sl = SortedList()
        res, left = 0, 0
        for right, cur in enumerate(nums):
            trie.add(prexor[right])
            sl.add(cur)
            while sl[-1] - sl[0] > k:
                trie.discard(prexor[left])
                sl.remove(nums[left])
                left += 1
            trie.xorAll(prexor[right + 1])
            res = max(res, trie.maximum())
            trie.xorAll(prexor[right + 1])
        return res
