from typing import List


class Solution:
    def minInterval(self, intervals: List[List[int]], queries: List[int]) -> List[int]:
        ...


print(Solution().minInterval(intervals=[[1, 4], [2, 4], [3, 6], [4, 4]], queries=[2, 3, 4, 5]))
