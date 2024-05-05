from typing import List
from collections import Counter

# 子序列：选哪几个
class Solution:
    def maxSubsequence(self, nums: List[int], k: int) -> List[int]:
        cand = sorted(nums, reverse=True)[:k]
        freq = Counter(cand)
        res = []
        for num in nums:
            if freq[num] > 0:
                freq[num] -= 1
                res.append(num)
        return res


print(Solution().maxSubsequence(nums=[-1, -2, 3, 4], k=3))
# 输出：[-1,3,4]
# 解释：
# 子序列有最大和：-1 + 3 + 4 = 6 。
