from collections import Counter
from itertools import permutations
from typing import List

# 请你返回 nums[i] + nums[j] （两个字符串连接）结果等于 target 的下标 (i, j) （需满足 i != j）的数目。


class Solution:
    def numOfPairs(self, nums: List[str], target: str) -> int:
        """哈希表做法：枚举前缀+后缀"""
        counter = Counter(nums)
        res = 0
        for i in range(1, len(target)):  # 枚举非空前后缀
            prefix = target[:i]
            suffix = target[i:]
            if prefix == suffix:
                res += (counter[prefix] - 1) * counter[prefix]
            else:
                res += counter[prefix] * counter[suffix]
        return res

        # return sum("".join(pair) == target for pair in permutations(nums, 2))


print(Solution().numOfPairs(nums=["777", "7", "77", "77"], target="7777"))
# 输出：4
# 解释：符合要求的下标对包括：
# - (0, 1)："777" + "7"
# - (1, 0)："7" + "777"
# - (2, 3)："77" + "77"
# - (3, 2)："77" + "77"
