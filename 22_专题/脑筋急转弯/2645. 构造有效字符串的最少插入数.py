# 给你一个字符串 word ，
# 你可以向其中任何位置插入 "a"、"b" 或 "c" 任意次，
# 返回使 word 有效 需要插入的最少字母数。
# !如果字符串可以由 "abc" 串联多次得到，则认为该字符串 有效 。


# 脑筋急转弯
# !如果si+1<=si的下标i有k个,那么完成插入后,至少有(k+1)个abc


class Solution:
    def addMinimum(self, word: str) -> int:
        rev = sum(pre >= cur for pre, cur in zip(word, word[1:]))
        return 3 * (1 + rev) - len(word)


assert Solution().addMinimum("aaa") == 6
