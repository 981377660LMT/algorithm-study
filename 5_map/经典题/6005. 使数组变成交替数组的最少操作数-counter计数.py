from collections import Counter
from typing import List, Tuple

MOD = int(1e9 + 7)

# 返回使数组变成交替数组的 最少操作数 。
class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        n = len(nums)
        g1 = nums[::2]
        g2 = nums[1::2]
        c1 = Counter(g1).most_common()
        c2 = Counter(g2).most_common()

        cand1, cand1Count = c1[0] if c1 else (0, 0)
        cand2, cand2Count = c1[1] if len(c1) > 1 else (0, 0)
        cand3, cand3Count = c2[0] if c2 else (0, 0)
        cand4, cand4Count = c2[1] if len(c2) > 1 else (0, 0)
        if cand1 != cand3:
            return n - cand1Count - cand3Count
        else:
            return min(n - cand2Count - cand3Count, n - cand1Count - cand4Count)


print(Solution().minimumOperations(nums=[3, 1, 3, 2, 4, 3]))
print(Solution().minimumOperations(nums=[1, 2, 2, 2, 2]))
print(Solution().minimumOperations(nums=[1, 2,]))
print(Solution().minimumOperations(nums=[1, 2, 3]))
print(Solution().minimumOperations(nums=[1,]))

