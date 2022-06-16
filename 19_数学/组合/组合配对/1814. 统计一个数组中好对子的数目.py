from collections import Counter
from itertools import product
from typing import List

# A[i] + rev(A[j]) == A[j] + rev(A[i])
# A[i] - rev(A[i]) == A[j] - rev(A[j])
# B[i] = A[i] - rev(A[i])

# Then it becomes an easy question that,
# how many pairs in B with B[i] == B[j]

MOD = int(1e9 + 7)


class Solution:
    def countNicePairs(self, A: List[int]) -> int:
        """一遍遍历 前不看后"""
        res = 0
        counter = Counter()

        # 保证i<j
        for num1 in A:
            num2 = int(str(num1)[::-1])
            res += counter[num1 - num2]
            counter[num1 - num2] += 1

        return res % MOD

    def countNicePairs2(self, A: List[int]) -> int:
        """先全部存起来再统计"""
        res = 0
        C = Counter(num - int(str(num)[::-1]) for num in A)
        for count in C.values():
            res += count * (count - 1) // 2
        return res % MOD


print(Solution().countNicePairs(A=[42, 11, 1, 97]))
# 输出：2
# 解释：两个坐标对为：
#  - (0,3)：42 + rev(97) = 42 + 79 = 121, 97 + rev(42) = 97 + 24 = 121 。
#  - (1,2)：11 + rev(1) = 11 + 1 = 12, 1 + rev(11) = 1 + 11 = 12 。
