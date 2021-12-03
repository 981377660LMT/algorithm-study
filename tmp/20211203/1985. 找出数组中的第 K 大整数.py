from typing import List
from heapq import heappop, heappush

# 返回 nums 中表示第 k 大整数的字符串。
class Solution:
    def kthLargestNumber(self, nums: List[str], k: int) -> str:
        pq = []
        for char in nums:
            heappush(pq, int(char))
            if len(pq) > k:
                heappop(pq)
        return str(pq[0])


print(Solution().kthLargestNumber(nums=["3", "6", "7", "10"], k=4))
# 输出："3"
# 解释：
# nums 中的数字按非递减顺序排列为 ["3","6","7","10"]
# 其中第 4 大整数是 "3"
