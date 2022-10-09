"""找出前缀异或的原始数组

相当于对前缀和/前缀异或数组求差分
"""
from typing import List


# 给你一个长度为 n 的 整数 数组 pref 。找出并返回满足下述条件且长度为 n 的数组 `arr` ：
# pref[i] = arr[0] ^ arr[1] ^ ... ^ arr[i].


class Solution:
    def findArray(self, pref: List[int]) -> List[int]:
        preXor = [0] + pref
        return [b ^ a for a, b in zip(preXor, preXor[1:])]


print(Solution().findArray(pref=[5, 2, 0, 3, 1]))
