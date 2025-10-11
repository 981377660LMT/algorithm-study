# https://leetcode.cn/problems/longest-common-prefix/description/
# 作用：对字符串进行精准匹配，预计算最长公共前缀lcp后，只需要看(lcp+1)个字符是否相同即可。
# !例如单词集合["flower","flow","flight"]，lcp=2，"fl"，则只需要比较word[:3]是否与text[:3]相同即可。

from typing import List


class Solution:
    def longestCommonPrefix(self, strs: List[str]) -> str:
        if not strs:
            return ""
        min_, max_ = min(strs), max(strs)
        for i, c in enumerate(min_):
            if c != max_[i]:
                return min_[:i]
        return min_
