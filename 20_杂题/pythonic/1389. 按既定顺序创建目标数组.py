from typing import List


class Solution:
    def createTargetArray(self, nums: List[int], index: List[int]) -> List[int]:
        res = []
        for v, i in zip(nums, index):
            res[i:i] = [v]
            # res.insert(i, v)
        return res

