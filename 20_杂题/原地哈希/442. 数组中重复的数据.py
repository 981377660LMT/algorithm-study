from typing import List


# 遍历时将 nums[|val|-1] 取反，若已为负说明 val 出现过两次
class Solution:
    def findDuplicates(self, nums: List[int]) -> List[int]:
        res = []
        for v in nums:
            pos = abs(v) - 1
            if nums[pos] < 0:
                res.append(pos + 1)
            else:
                nums[pos] = -nums[pos]
        return res
