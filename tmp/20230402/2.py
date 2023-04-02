from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 。请你创建一个满足以下条件的二维数组：

# 二维数组应该 只 包含数组 nums 中的元素。
# 二维数组中的每一行都包含 不同 的整数。
# 二维数组的行数应尽可能 少 。
# 返回结果数组。如果存在多种答案，则返回其中任何一种。

# 请注意，二维数组的每一行上可以存在不同数量的元素。
class Solution:
    def findMatrix(self, nums: List[int]) -> List[List[int]]:
        # S = set(nums)
        C = Counter(nums)
        res = []
        while True:
            cur = []
            for k in C:
                if C[k] > 0:
                    cur.append(k)
                    C[k] -= 1
            if not cur:
                break
            res.append(cur)
        return res


print(Solution().findMatrix(nums=[1, 2, 3, 4]))
