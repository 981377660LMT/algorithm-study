from typing import List
from itertools import product
from collections import Counter

MOD = int(1e9 + 7)

# 923. 三数之和的多种可能
# 返回满足 i < j < k 且 A[i] + A[j] + A[k] == target 的元组 i, j, k 的数量

# 3 <= arr.length <= 3000
# 0 <= arr[i] <= 100
# 0 <= target <= 300

# 也可以dpIndexRemain


class Solution:
    def threeSumMulti(self, arr: List[int], target: int) -> int:
        C = Counter(arr)
        res = 0
        for n1, n2 in product(C, repeat=2):
            n3 = target - n1 - n2
            if not n1 <= n2 <= n3:
                continue
            # 三个相等/有两个相等/都不等
            if n1 == n2 == n3:
                res += C[n1] * (C[n1] - 1) * (C[n1] - 2) // 6
            elif n1 == n2 != n3:
                res += (C[n1] * (C[n1] - 1) // 2) * C[n3]
            elif n1 != n2 == n3:
                res += (C[n2] * (C[n2] - 1) // 2) * C[n1]
            else:
                res += C[n1] * C[n2] * C[n3]

        return res % MOD


print(Solution().threeSumMulti(arr=[1, 1, 2, 2, 3, 3, 4, 4, 5, 5], target=8))
