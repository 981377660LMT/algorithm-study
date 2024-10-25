# 如果数组中存在至少重复 k 次且长度为 m 的模式，则返回 true ，否则返回  false 。
from typing import List


class Solution:
    def containsPattern(self, arr: List[int], m: int, k: int) -> bool:
        return any(arr[i : i + m] * k == arr[i : i + m * k] for i in range(len(arr) - m * k + 1))

