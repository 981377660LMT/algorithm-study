from typing import List
from collections import Counter

# 返回 arr 中第 k 个 独一无二的字符串 。如果 少于 k 个独一无二的字符串，那么返回 空字符串 "" 。
class Solution:
    def kthDistinct(self, arr: List[str], k: int) -> str:
        num_freq = Counter(arr)
        nums = []
        for x in arr:
            if num_freq[x] == 1:
                nums.append(x)
        return '' if len(nums) < k else nums[k - 1]


# 其实用生成器最好
