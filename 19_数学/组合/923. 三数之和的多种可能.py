from typing import List
from itertools import combinations_with_replacement
from collections import Counter

MOD = int(1e9 + 7)

# 返回满足 i < j < k 且 A[i] + A[j] + A[k] == target 的元组 i, j, k 的数量


class Solution:
    def threeSumMulti(self, arr: List[int], target: int) -> int:
        c = Counter(arr)
        # combinations_with_replacement 作用 ：返回指定长度的组合，组合内元素可重复
        # print(*combinations_with_replacement(c, 2))
        res = 0
        for i, j in combinations_with_replacement(c, 2):
            k = target - i - j
            # 三个相等/有两个相等(只取一种情况)/都不想等(只取一种情况)
            if i == j == k:
                res += c[i] * (c[i] - 1) * (c[i] - 2) // 6
            elif i == j != k:
                res += (c[i] * (c[i] - 1) // 2) * c[k]
            elif k > i and k > j:
                res += c[i] * c[j] * c[k]

        return res % MOD


print(Solution().threeSumMulti(arr=[1, 1, 2, 2, 3, 3, 4, 4, 5, 5], target=8))
