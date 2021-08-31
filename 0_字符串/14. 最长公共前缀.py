# description 编写一个函数来查找字符串数组中的最长公共前缀
# 仅需最大、最小字符串的最长公共前缀


from typing import List


class Solution:
    def longestCommonPrefix(self, strs: List[str]) -> str:
        res = ''
        z = list(zip(*strs))
        for item in z:
            if len(set(item)) == 1:
                res += item[0]
            else:
                break
        return res


print(Solution().longestCommonPrefix(["flower", "flow", "flight"]))

