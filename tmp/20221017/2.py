from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个由 正 整数组成的数组 nums 。

# 你必须取出数组中的每个整数，反转其中每个数位，并将反转后得到的数字添加到数组的末尾。这一操作只针对 nums 中原有的整数执行。

# 返回结果数组中 不同 整数的数目。
class Solution:
    def countDistinctIntegers(self, nums: List[int]) -> int:
        return len(set(nums) | set(int(str(x)[::-1]) for x in nums))
