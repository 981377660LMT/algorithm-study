from typing import List

# 请你返回 任意 一个由 n 个 各不相同 的整数组成的数组，并且这 n 个数相加和为 0 。
class Solution:
    def sumZero(self, n: int) -> List[int]:
        return list(range(1 - n, n, 2))

