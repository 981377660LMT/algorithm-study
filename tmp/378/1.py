from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个 正整数 数组 nums 。

# 你需要检查是否可以从数组中选出 两个或更多 元素，满足这些元素的按位或运算（ OR）结果的二进制表示中 至少 存在一个尾随零。

# 例如，数字 5 的二进制表示是 "101"，不存在尾随零，而数字 4 的二进制表示是 "100"，存在两个尾随零。

# 如果可以选择两个或更多元素，其按位或运算结果存在尾随零，返回 true；否则，返回 false 。


class Solution:
    def hasTrailingZeros(self, nums: List[int]) -> bool:
        return len([v for v in nums if v & 1 == 0]) >= 2


print(Solution().hasTrailingZeros([1, 3, 5, 7, 9]))
