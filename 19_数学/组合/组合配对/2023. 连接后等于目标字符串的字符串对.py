from collections import Counter
from typing import List

# 请你返回 nums[i] + nums[j] （两个字符串连接）结果等于 target 的下标 (i, j) （需满足 i != j）的数目。
class Solution:
    def numOfPairs(self, nums: List[str], target: str) -> int:
        freq = Counter(nums)
        res = 0
        for prefix, count in freq.items():
            if target.startswith(prefix):
                suffix = target[len(prefix) :]
                if prefix != suffix:
                    res += freq[suffix] * freq[prefix]
                else:
                    res += freq[prefix] * (freq[prefix] - 1)

        return res


print(Solution().numOfPairs(nums=["777", "7", "77", "77"], target="7777"))
# 输出：4
# 解释：符合要求的下标对包括：
# - (0, 1)："777" + "7"
# - (1, 0)："7" + "777"
# - (2, 3)："77" + "77"
# - (3, 2)："77" + "77"

