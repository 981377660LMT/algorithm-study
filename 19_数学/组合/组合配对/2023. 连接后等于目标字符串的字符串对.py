from collections import Counter
from itertools import permutations
from typing import List

# 请你返回 nums[i] + nums[j] （两个字符串连接）结果等于 target 的下标 (i, j) （需满足 i != j）的数目。


class Solution:
    def numOfPairs(self, nums: List[str], target: str) -> int:
        return sum("".join(pair) == target for pair in permutations(nums, 2))
        freq = Counter(nums)
        res = 0
        for pre, _ in freq.items():
            if target.startswith(pre):
                suf = target[len(pre) :]
                if pre != suf:
                    res += freq[suf] * freq[pre]
                else:
                    res += freq[pre] * (freq[pre] - 1)

        return res


print(Solution().numOfPairs(nums=["777", "7", "77", "77"], target="7777"))
# 输出：4
# 解释：符合要求的下标对包括：
# - (0, 1)："777" + "7"
# - (1, 0)："7" + "777"
# - (2, 3)："77" + "77"
# - (3, 2)："77" + "77"

