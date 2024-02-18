# 3035. 回文字符串的最大数量
# https://leetcode.cn/problems/maximum-palindromes-after-operations/description/
# 给定一些words，可以任意交换两个单词的字母.
# !返回在执行一些操作后，words 中可以包含的回文字符串的 最大 数量.
#
# 先把所有字母都取出来，然后考虑如何填入各个字符串。
# 如果一个奇数长度字符串最终是回文串，那么它正中间的那个字母填什么都可以。
# 既然如此，不妨先把左右的字母填了，最后在往正中间填入字母。


from collections import Counter
from typing import List


class Solution:
    def maxPalindromesAfterOperations(self, words: List[str]) -> int:
        counter = Counter("".join(words))  # 统计每个字符的数量
        counts = sum(v // 2 for v in counter.values())
        lens = sorted(len(w) // 2 for w in words)
        res = 0
        for len_ in lens:
            if counts < len_:
                break
            counts -= len_
            res += 1
        return res
