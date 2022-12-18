# 通过连接另一个数组的子数组得到一个数组
# !连接子数组

from typing import List


class Solution:
    def canChoose(self, groups: List[List[int]], nums: List[int]) -> bool:
        s = "".join([f"#{num}#" for num in nums])
        start = 0
        for group in groups:
            target = "".join([f"#{num}#" for num in group])
            pos = s.find(target, start)
            if pos == -1:
                return False
            start = pos + len(target)
        return True


assert Solution().canChoose(groups=[[1, -1, -1], [3, -2, 0]], nums=[1, -1, 0, 1, -1, -1, 3, -2, 0])
