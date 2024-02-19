# 100229. 最长公共前缀的长度
# https://leetcode.cn/problems/find-the-length-of-the-longest-common-prefix/description/
# 你需要找出属于 arr1 的整数 x 和属于 arr2 的整数 y 组成的所有数对 (x, y) 之中最长的公共前缀的长度。
# 返回所有数对之中最长公共前缀的长度。如果它们之间不存在公共前缀，则返回 0 。


from typing import List


class Solution:
    def longestCommonPrefix(self, arr1: List[int], arr2: List[int]) -> int:
        A = list(map(str, arr1))
        B = list(map(str, arr2))
        pre = set()
        for s in A:
            for i in range(1, len(s) + 1):
                pre.add(s[:i])

        res = 0
        for s in B:
            for i in range(1, len(s) + 1):
                if s[:i] in pre:
                    res = max(res, i)
                else:
                    break
        return res
