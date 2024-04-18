# 2935. 找出强数对的最大异或值 II
# https://leetcode.cn/problems/maximum-strong-pair-xor-ii/solutions/2523213/0-1-trie-hua-dong-chuang-kou-pythonjavac-gvv2/
# 给你一个下标从 0 开始的整数数组 nums 。如果一对整数 x 和 y 满足以下条件，则称其为 强数对 ：
# |x - y| <= min(x, y)
# 你需要从 nums 中选出两个整数，且满足：这两个整数可以形成一个强数对，并且它们的按位异或（XOR）值是在该数组所有强数对中的 最大值 。
# 返回数组 nums 所有可能的强数对中的 最大 异或值。
# 注意，你可以选择同一个整数两次来形成一个强数对。


# !一边遍历数组，一边记录每个 key 对应的最大的 nums[i]。
# 由于数组已经排好序，所以每个 key 对应的 x=nums[i] 一定是当前最大的，只要 2x≥y，就说明这个比特位可以是 1。

from typing import List


class Solution:
    def maximumStrongPairXor(self, nums: List[int]) -> int:
        nums = sorted(nums)
        bitLen = max(nums, default=0).bit_length()
        res, mask = 0, 0
        for b in range(bitLen - 1, -1, -1):
            mask |= 1 << b
            cand = res | (1 << b)  # 试填，这个比特位可以是1吗？
            maxNum = dict()
            for v in nums:
                maskV = v & mask
                if cand ^ maskV in maxNum and maxNum[cand ^ maskV] * 2 >= v:
                    res = cand
                    break
                maxNum[maskV] = v
        return res
