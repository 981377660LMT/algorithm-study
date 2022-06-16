from typing import List
from collections import Counter

# 大餐 是指 恰好包含两道不同餐品 的一餐，其美味程度之和等于 2 的幂。
# 1 <= deliciousness.length <= 10^5
# 0 <= deliciousness[i] <= 220
# 返回你可以用数组中的餐品做出的不同 大餐 的数量


class Solution:
    def countPairs(self, deliciousness: List[int]) -> int:
        C = Counter(deliciousness)

        res = 0
        for num in C:
            for p in range(32):
                need = pow(2, p) - num
                if num == need:
                    res += C[num] * (C[num] - 1) // 2
                elif num < need and need in C:
                    res += C[num] * C[need]

        return res % int(1e9 + 7)


print(Solution().countPairs(deliciousness=[1, 3, 5, 7, 9]))
# 输出：4
# 解释：大餐的美味程度组合为 (1,3) 、(1,7) 、(3,5) 和 (7,9) 。
# 它们各自的美味程度之和分别为 4 、8 、8 和 16 ，都是 2 的幂。
