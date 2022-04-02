from typing import List


class Solution:
    def solve(self, pieces: List[List[int]], target: List[int]) -> int:
        # 组内不能换位=>当成一个整体，用首部元素作为key
        chunk = {nums[0]: nums for nums in pieces if nums}
        res = []
        for num in target:
            if num in chunk:
                res.extend(chunk[num])
        return res == target


# 是否能将部分块前后连接，使其成为target
# pieces中的所有数都不同


print(Solution().solve(pieces=[[1], [3, 4], [5, 6], [2]], target=[1, 2, 3, 4]))
