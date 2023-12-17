# 区间内头尾相同子串的数量
# https://leetcode.cn/problems/number-of-same-end-substrings/description/
# 给定一个字符串和一个查询数组，返回每个查询[left,right]对应范围内，所有头尾相同的子串的数量。
# 2 <= s.length <= 3 * 104
# s consists only of lowercase English letters.
# 1 <= queries.length <= 3 * 104
# queries[i] = [li, ri]
# 0 <= li <= ri < s.length

from typing import List
from alphaPresum import alphaPresum


class Solution:
    def sameEndSubstringCount(self, s: str, queries: List[List[int]]) -> List[int]:
        S = alphaPresum(s)
        res = []
        for left, right in queries:
            cur = 0
            for i in range(26):
                freq = S(left, right + 1, i + 97)
                cur += freq * (freq + 1) // 2
            res.append(cur)
        return res
