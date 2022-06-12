# description 编写一个函数来查找字符串数组中的最长公共前缀
# 仅需最大、最小字符串的最长公共前缀

from typing import List


class Solution:
    def longestCommonPrefix(self, strs: List[str]) -> str:
        """纵向扫描"""
        res = []
        cols = zip(*strs)
        for col in cols:
            if len(set(col)) == 1:
                res.append(col[0])
            else:
                break
        return ''.join(res)

    # 只需要比较最大最小的公共前缀就是整个数组的公共前缀
    def longestCommonPrefix2(self, strs):
        if not strs:
            return ""
        s1 = min(strs)
        s2 = max(strs)
        for i, x in enumerate(s1):
            if x != s2[i]:
                return s2[:i]
        return s1

    def longestCommonPrefix3(self, strs: str):
        """二分(通用思想)"""
        if not strs:
            return ""
        s1 = min(strs)
        s2 = max(strs)
        for i, x in enumerate(s1):
            if x != s2[i]:
                return s2[:i]
        return s1


#
# zip法

print(Solution().longestCommonPrefix(["flower", "flow", "flight"]))

